//
// FilePath    : video-trim\utils.go
// Author      : jiaopengzi
// Blog        : https://jiaopengzi.com
// Copyright   : Copyright (c) 2026 by jiaopengzi, All Rights Reserved.
// Description : 工具
//

package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 初始化目录
func initDir() {
	// 初始化, 确保上传和输出目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("failed to create upload directory: %v", err)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}
}

// getLocalIP 获取本机局域网 IP, 用于提示用户
func getLocalIP() string {
	// 获取本机网络接口地址列表
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	// 选择第一个非回环的 IPv4 地址作为局域网提示
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return "127.0.0.1"
}

// parseHeadSec 从请求中解析裁剪开头秒数
func parseHeadSec(r *http.Request) int {
	// 从表单中解析裁剪开头秒数, 若无或无效则使用默认值
	head := headTrimSeconds

	if headStr := r.FormValue("head"); headStr != "" {
		if v, err := strconv.Atoi(headStr); err == nil && v >= 0 {
			head = v
		}
	}

	return head
}

// parseTailSec 从请求中解析去尾秒数
func parseTailSec(r *http.Request) int {
	tail := tailSeconds

	if tailStr := r.FormValue("tail"); tailStr != "" {
		if v, err := strconv.Atoi(tailStr); err == nil && v >= 0 {
			tail = v
		}
	}

	return tail
}

// ensurePostMethod 确保请求方法为 POST, 否则直接响应错误
func ensurePostMethod(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != "POST" {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return false
	}

	return true
}

// parseMultipartFormOrRespond 解析 multipart 表单并在出错时直接响应
func parseMultipartFormOrRespond(w http.ResponseWriter, r *http.Request) bool {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		// 选择语言并加载翻译
		lang := detectLangFromRequest(r)
		i18n := getLocale(lang)

		if strings.Contains(err.Error(), "request body too large") || strings.Contains(err.Error(), "http: request body too large") {
			msg := fmt.Sprintf(i18n[KeyRequestBodyTooLarge], humanReadableBytes(maxUploadSize))

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			fmt.Fprintln(w, msg)

			return false
		}

		// 通用解析错误信息
		http.Error(w, i18n[KeyRequestParseError], http.StatusBadRequest)

		return false
	}

	return true
}

// getUploadedFilesOrRespond 获取上传的文件列表并在无文件时直接响应
func getUploadedFilesOrRespond(w http.ResponseWriter, r *http.Request) ([]*multipart.FileHeader, bool) {
	files := r.MultipartForm.File["videos"]
	if len(files) == 0 {
		// 选择语言并加载翻译
		lang := detectLangFromRequest(r)
		i18n := getLocale(lang)

		msg := i18n[KeySelectAtLeastOne]

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<script>alert(%s);location.href='/';</script>`, strconv.Quote(msg))

		return nil, false
	}

	return files, true
}

// checkFilesSizeOrRespond 检查上传文件大小并在超出限制时直接响应
func checkFilesSizeOrRespond(w http.ResponseWriter, files []*multipart.FileHeader, lang string) bool {
	for _, hdr := range files {
		if hdr.Size > maxUploadSize {
			i18n := getLocale(lang)

			prefix := i18n[KeyFileTooLargePrefix]
			suffix := i18n[KeyFileTooLargeSuffix]
			end := i18n[KeyFileTooLargeEnd]

			msg := prefix + hdr.Filename + suffix + humanReadableBytes(maxUploadSize) + end

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			fmt.Fprintln(w, msg)

			return false
		}
	}

	return true
}

// checkFilesMagicOrRespond 使用魔法数字(文件签名)校验上传文件是否为视频格式
func checkFilesMagicOrRespond(w http.ResponseWriter, files []*multipart.FileHeader, lang string) bool {
	const sniffLen = 64

	i18n := getLocale(lang)

	for _, hdr := range files {
		f, err := hdr.Open()
		if err != nil {
			log.Printf("open uploaded file error: %v", err)

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, `<script>alert(%s);location.href='/';</script>`, strconv.Quote(fmt.Sprintf(i18n[KeyCannotReadFile], hdr.Filename)))

			return false
		}

		buf := make([]byte, sniffLen)
		n, err := io.ReadFull(f, buf)
		f.Close()

		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Printf("read uploaded file error: %v", err)

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, `<script>alert(%s);location.href='/';</script>`, strconv.Quote(fmt.Sprintf(i18n[KeyCannotReadFile], hdr.Filename)))

			return false
		}

		if n == 0 {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, `<script>alert(%s);location.href='/';</script>`, strconv.Quote(fmt.Sprintf(i18n[KeyFileEmptyOrUnreadable], hdr.Filename)))

			return false
		}

		if !isVideoMagic(buf[:n]) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, `<script>alert(%s);location.href='/';</script>`, strconv.Quote(fmt.Sprintf(i18n[KeyNotSupportedVideo], hdr.Filename)))

			return false
		}
	}

	return true
}

// isVideoMagic 根据常见视频文件的魔法数字(文件签名)判断是否可能为视频
func isVideoMagic(b []byte) bool {
	// 检查 MP4 / MOV: 通常在偏移 4 处包含 "ftyp"
	if len(b) >= 12 {
		if string(b[4:8]) == "ftyp" {
			return true
		}
	}

	if len(b) >= 4 {
		// Matroska / WebM (EBML) 头: 0x1A45DFA3
		if b[0] == 0x1A && b[1] == 0x45 && b[2] == 0xDF && b[3] == 0xA3 {
			return true
		}

		// AVI: "RIFF" ... "AVI "
		if string(b[0:4]) == "RIFF" && len(b) >= 12 && string(b[8:12]) == "AVI " {
			return true
		}

		// FLV
		if len(b) >= 3 && string(b[0:3]) == "FLV" {
			return true
		}

		// MPEG-TS: 以 0x47 (sync byte) 开头的流(简单判断, 可能有误判但覆盖常见 TS)
		if b[0] == 0x47 {
			return true
		}
	}

	return false
}

// processSingleFile 保存上传文件、调用 ffmpeg, 并返回输出文件名
func processSingleFile(hdr *multipart.FileHeader, idx int, headSec int, tailSec int) (string, error) {
	// 打开上传的文件头, 获取读取流
	f, err := hdr.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	// 推断文件扩展名, 默认使用 .mp4
	ext := filepath.Ext(hdr.Filename)
	if ext == "" {
		ext = ".mp4"
	}

	base := filepath.Base(hdr.Filename)
	nameOnly := strings.TrimSuffix(base, ext)

	// 在 uploads 目录生成临时输入文件路径, 防止文件名冲突
	inputPath := filepath.Join(uploadDir, fmt.Sprintf("input_%d_%d%s", time.Now().Unix(), idx, ext))

	// 将上传文件写入磁盘
	outFile, err := os.Create(inputPath)
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(outFile, f); err != nil {
		outFile.Close()
		os.Remove(inputPath)

		return "", err
	}

	outFile.Close()

	// 生成输出文件名并确保不会覆盖已有文件
	outName := fmt.Sprintf("%s-head%s", nameOnly, ext)
	outputPath := filepath.Join(outputDir, outName)

	if _, err := os.Stat(outputPath); err == nil {
		outName = fmt.Sprintf("%s-head-%d%s", nameOnly, time.Now().Unix(), ext)
		outputPath = filepath.Join(outputDir, outName)
	}

	// 调用 ffmpeg 进行剪切处理
	if err := runFFmpeg(inputPath, outputPath, headSec, tailSec); err != nil {
		// 如果处理失败, 清理临时输入文件并返回错误
		os.Remove(inputPath)
		return "", err
	}

	// 处理完成后删除临时输入文件并返回输出文件名
	os.Remove(inputPath)

	return outName, nil
}

// runFFmpeg 简单包装 ffmpeg 调用, 校验并规范化参数以避免可控的命令注入
func runFFmpeg(inputPath, outputPath string, headSec int, tailSec int) error {
	// 执行流程：校验参数 -> 解析并校验路径 -> 构建参数 -> 执行 ffmpeg
	if err := validateHeadTail(&headSec, tailSec); err != nil {
		return err
	}

	// 寻找 ffmpeg 可执行文件路径
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("ffmpeg not found in PATH: %w", err)
	}

	// 解析并校验输入/输出路径
	absInput, absOutput, err := resolveAndValidatePaths(inputPath, outputPath)
	if err != nil {
		return err
	}

	// 构建 ffmpeg 参数
	args, err := buildFFmpegArgs(absInput, absOutput, headSec, tailSec)
	if err != nil {
		return err
	}

	// 执行 ffmpeg 命令
	cmd := exec.Command(ffmpegPath, args...)

	// 捕获输出以便调试
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w: %s", err, string(out))
	}

	return nil
}

// validateHeadTail 校验并规范化 head/tail 参数
func validateHeadTail(headSec *int, tailSec int) error {
	if *headSec < 0 {
		return fmt.Errorf("invalid head seconds")
	}

	if tailSec < 0 {
		return fmt.Errorf("invalid tail seconds")
	}

	const maxHead = 24 * 3600 // 不允许超过一天
	if *headSec > maxHead {
		*headSec = maxHead
	}

	return nil
}

// resolveAndValidatePaths 返回输入/输出的绝对路径并校验它们位于受控目录下
func resolveAndValidatePaths(inputPath, outputPath string) (string, string, error) {
	absInput, err := filepath.Abs(inputPath)
	if err != nil {
		return "", "", fmt.Errorf("invalid input path: %w", err)
	}

	absOutput, err := filepath.Abs(outputPath)
	if err != nil {
		return "", "", fmt.Errorf("invalid output path: %w", err)
	}

	absUploadDir, err := filepath.Abs(uploadDir)
	if err != nil {
		return "", "", fmt.Errorf("invalid upload dir: %w", err)
	}

	absOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		return "", "", fmt.Errorf("invalid output dir: %w", err)
	}

	if !(absInput == absUploadDir || strings.HasPrefix(absInput, absUploadDir+string(os.PathSeparator))) {
		return "", "", fmt.Errorf("input path not allowed")
	}

	if !(absOutput == absOutputDir || strings.HasPrefix(absOutput, absOutputDir+string(os.PathSeparator))) {
		return "", "", fmt.Errorf("output path not allowed")
	}

	return absInput, absOutput, nil
}

// buildFFmpegArgs 根据是否需要去尾构建 ffmpeg 参数, 必要时会调用 ffprobe 获取时长
func buildFFmpegArgs(absInput, absOutput string, headSec int, tailSec int) ([]string, error) {
	if tailSec <= 0 {
		return []string{
			"-ss", strconv.Itoa(headSec),
			"-i", absInput,
			"-c", "copy",
			"-avoid_negative_ts", "make_zero",
			absOutput,
		}, nil
	}

	// 需要去尾：先查找 ffprobe 并获取时长
	ffprobePath, err := exec.LookPath("ffprobe")
	if err != nil {
		return nil, fmt.Errorf("ffprobe not found in PATH: %w", err)
	}

	duration, err := getMediaDuration(ffprobePath, absInput)
	if err != nil {
		return nil, fmt.Errorf("failed to get media duration: %w", err)
	}

	if duration <= 0 {
		return nil, fmt.Errorf("invalid media duration: %v", duration)
	}

	// 计算结束时间并验证
	end := duration - float64(tailSec)
	if end <= float64(headSec) {
		return nil, fmt.Errorf("head + tail exceeds media duration")
	}

	dur := end - float64(headSec)
	durStr := strconv.FormatFloat(dur, 'f', 3, 64)

	args := []string{
		"-ss", strconv.Itoa(headSec),
		"-i", absInput,
		"-t", durStr,
		"-c", "copy",
		"-avoid_negative_ts", "make_zero",
		absOutput,
	}

	return args, nil
}

// getMediaDuration 使用 ffprobe 获取媒体文件时长(秒)
func getMediaDuration(ffprobePath, input string) (float64, error) {
	cmd := exec.Command(ffprobePath, "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", input)

	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	s := strings.TrimSpace(string(out))
	if s == "" {
		return 0, fmt.Errorf("empty duration from ffprobe")
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	return f, nil
}

// humanReadableBytes 将字节数格式化为人类可读的字符串(例如 8.0 MB)
func humanReadableBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}

	div := int64(unit)
	exp := 0

	for n/div >= unit {
		div *= unit
		exp++
	}

	value := float64(n) / float64(div)
	units := []string{"KB", "MB", "GB", "TB"}
	u := units[0]

	if exp >= 0 && exp < len(units) {
		u = units[exp]
	}

	return fmt.Sprintf("%.1f %s", value, u)
}

// generateResponse 输出处理完成后的下载页面
func generateResponse(w http.ResponseWriter, processed []string, r *http.Request) {
	// 生成处理完成后的 HTML 页面, 列出可下载文件
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	jsItems := []string{}
	for _, fn := range processed {
		jsItems = append(jsItems, strconv.Quote(fn))
	}

	jsArray := "[" + strings.Join(jsItems, ",") + "]"

	// 选择语言并加载翻译
	lang := detectLangFromRequest(r)
	i18n := getLocale(lang)

	title := i18n[KeyProcessedTitle]
	downloadAll := i18n[KeyDownloadAll]
	downloadText := i18n[KeyDownload]
	returnUpload := i18n[KeyReturnUpload]

	// 输出 HTML 内容
	fmt.Fprintf(w, "<!DOCTYPE html><html><head><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width,initial-scale=1\"><title>%s</title><style>body{font-family:system-ui,Arial;background:#f7f8fa;margin:0;padding:12px} .wrap{max-width:720px;margin:0 auto} h2{margin:0 0 12px;font-size:18px} .downloadAllBtn{display:block;width:100%%;padding:12px;border-radius:12px;background:#10b981;color:#fff;border:0;font-size:16px;margin-bottom:12px} .list{display:flex;flex-direction:column;gap:10px} .item{display:flex;align-items:center;justify-content:space-between;background:#fff;padding:12px;border-radius:12px;box-shadow:0 6px 18px rgba(2,6,23,0.06)} .name{flex:1;font-size:14px;color:#0f172a;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;margin-right:10px} .actions{display:flex;gap:8px} .btn{background:#0366d6;color:#fff;padding:8px 12px;border-radius:10px;text-decoration:none;font-size:14px} .muted{color:#6b7280;font-size:13px;margin-top:8px} .returnBtn{display:inline-block;margin-top:10px;padding:10px 14px;border-radius:10px;background:#0366d6;color:#fff;text-decoration:none;font-weight:600}</style></head><body>", title)
	fmt.Fprintln(w, "<div class=\"wrap\">")
	fmt.Fprintf(w, "<h2>%s</h2>", title)
	fmt.Fprintf(w, "<button id=\"downloadAll\" class=\"downloadAllBtn\">%s</button>", downloadAll)
	fmt.Fprintln(w, "<div class=\"list\">")

	// 列出所有处理成功的文件
	for _, fn := range processed {
		link := fmt.Sprintf("/download/%s", fn)
		fmt.Fprintf(w, "<div class=\"item\"><div class=\"name\">%s</div><div class=\"actions\"><a class=\"btn\" href=\"%s\" download>%s</a></div></div>", fn, link, downloadText)
	}

	fmt.Fprintln(w, "</div>")

	// 若无文件被处理成功, 提示用户
	if len(processed) == 0 {
		muted := i18n[KeyNoProcessedFilesHint]
		fmt.Fprintf(w, "<p class=\"muted\">%s</p>", muted)
	}

	// 返回上传页面链接
	fmt.Fprintf(w, "<p><a class=\"returnBtn\" href=\"/\">%s</a></p>", returnUpload)

	// 页面端提供「下载全部」功能: 逐个请求 /download/ 并触发保存
	fmt.Fprintf(w, `<script>var files=%s;(function(){document.getElementById('downloadAll').addEventListener('click',async function(){if(!files||!files.length)return;this.disabled=true;for(let i=0;i<files.length;i++){let f=files[i];try{let resp=await fetch('/download/'+encodeURIComponent(f));if(!resp.ok){console.error('fetch failed',f);continue}let blob=await resp.blob();let url=URL.createObjectURL(blob);let a=document.createElement('a');a.href=url;a.download=f;document.body.appendChild(a);a.click();a.remove();URL.revokeObjectURL(url);}catch(e){console.error(e)}await new Promise(r=>setTimeout(r,500));}this.disabled=false});})();</script>`, jsArray)
	fmt.Fprintln(w, "</div>")
	fmt.Fprintln(w, "</body></html>")
}
