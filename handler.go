//
// FilePath    : video-trim\handler.go
// Author      : jiaopengzi
// Blog        : https://jiaopengzi.com
// Copyright   : Copyright (c) 2026 by jiaopengzi, All Rights Reserved.
// Description : HTTP 处理函数
//

package main

import (
	"fmt"
	"html/template"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// handleHome 处理首页请求, 渲染上传页面
func handleHome(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// 解析并执行模板
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		log.Printf("parse template error: %v", err)
		http.Error(w, "parse template error", http.StatusInternalServerError)

		return
	}

	// 选择语言并加载翻译
	lang := detectLangFromRequest(r)
	i18n := getLocale(lang)
	data := struct {
		Head              int
		Tail              int
		MaxUpload         int64
		MaxUploadReadable string
		I18n              map[string]string
		AvailableLocales  []LocaleMeta
		Lang              string
	}{
		Head:              headTrimSeconds,
		Tail:              tailSeconds,
		MaxUpload:         maxUploadSize,
		MaxUploadReadable: humanReadableBytes(maxUploadSize),
		I18n:              i18n,
		AvailableLocales:  GetAvailableLocales(lang),
		Lang:              lang,
	}

	// 执行模板并写入响应
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("template execute error: %v", err)
		http.Error(w, "template execute error", http.StatusInternalServerError)

		return
	}
}

// handleUpload 处理文件上传和剪切请求
func handleUpload(w http.ResponseWriter, r *http.Request) {
	// 仅允许 POST 方法
	if !ensurePostMethod(w, r) {
		return
	}

	// 解析 multipart 表单
	if !parseMultipartFormOrRespond(w, r) {
		return
	}

	// 获取上传的文件列表
	files, ok := getUploadedFilesOrRespond(w, r)
	if !ok {
		return
	}

	// 选择语言并加载翻译
	lang := detectLangFromRequest(r)

	// 检查每个文件大小是否超出限制
	if !checkFilesSizeOrRespond(w, files, lang) {
		return
	}

	// 魔法数字校验, 确保上传文件看起来像视频文件
	if !checkFilesMagicOrRespond(w, files, lang) {
		return
	}

	// 解析要剪掉的秒数和去尾秒数
	headSec := parseHeadSec(r)
	tailSec := parseTailSec(r)

	// 如果 head 和 tail 都为 0, 则无需处理
	if headSec == 0 && tailSec == 0 {
		respondNoTrim(w, r)
		return
	}

	// 逐个处理文件并生成响应
	processed := processUploadedFiles(files, headSec, tailSec)

	// 生成并返回下载页面
	generateResponse(w, processed, r)
}

// detectLangFromRequest 返回请求中优先级为: ?lang -> cookie(lang) -> defaultLang 的语言代码
func detectLangFromRequest(r *http.Request) string {
	qLang := r.URL.Query().Get("lang")

	lang := qLang
	if lang == "" {
		if c, err := r.Cookie("lang"); err == nil {
			lang = c.Value
		} else {
			lang = defaultLang
		}
	}

	return lang
}

// respondNoTrim 使用 i18n 输出 JS alert 并重定向回首页
func respondNoTrim(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	lang := detectLangFromRequest(r)
	i18n := getLocale(lang)

	msg := i18n[KeyAlertNoTrim]
	esc := template.JSEscapeString(msg)
	fmt.Fprintf(w, "<script>alert('%s');location.href='/';</script>", esc)
}

// processUploadedFiles 逐个调用 processSingleFile 并返回成功的输出文件名列表
func processUploadedFiles(files []*multipart.FileHeader, headSec, tailSec int) []string {
	processed := []string{}

	for idx, hdr := range files {
		outName, err := processSingleFile(hdr, idx, headSec, tailSec)
		if err != nil {
			log.Printf("process file %s error: %v", hdr.Filename, err)
			continue
		}

		processed = append(processed, outName)
	}

	return processed
}

// saveFile 处理文件下载请求
func handleDownload(w http.ResponseWriter, r *http.Request) {
	// 提取文件名并在输出目录中寻找对应文件, 若不存在则返回 404
	// 使用 filepath.Base 防止路径遍历
	filename := filepath.Base(strings.TrimPrefix(r.URL.Path, "/download/"))

	filePath := filepath.Join(outputDir, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// 使用标准库直接提供文件下载
	http.ServeFile(w, r, filePath)
}

// 处理清理 uploads 和 outputs 目录下的所有内容(保留目录)
func handleClear(w http.ResponseWriter, r *http.Request) {
	// 仅允许 POST 请求触发清理操作
	if r.Method != "POST" {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	// 遍历 uploads 和 outputs 目录, 删除所有子项但保留目录本身
	dirs := []string{uploadDir, outputDir}
	for _, d := range dirs {
		entries, err := os.ReadDir(d)
		if err != nil {
			log.Printf("read dir error: %v", err)
			continue
		}

		for _, e := range entries {
			p := filepath.Join(d, e.Name())
			if err := os.RemoveAll(p); err != nil {
				log.Printf("remove %s error: %v", p, err)
			}
		}
	}

	// 清理完成后重定向回首页
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
