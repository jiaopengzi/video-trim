package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// LocaleMeta 描述可用语言
type LocaleMeta struct {
	Code     string
	Name     string
	Selected bool
}

var (
	// locales 存储所有加载的语言文件内容
	locales = map[string]map[string]string{}

	// availableLocales 存储可用语言列表
	availableLocales = []LocaleMeta{}

	// defaultLang 默认语言代码
	defaultLang = "zh"
)

// loadLocales 从 locales 目录加载所有 json 翻译文件
func loadLocales() {
	dir := "locales"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 目录不存在：创建并写入默认文件
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("创建 locales 目录失败: %v", err)
			return
		}

		EnsureLocaleExists()
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		log.Printf("扫描 locale 目录失败: %v", err)
		return
	}

	for _, f := range files {
		base := filepath.Base(f)
		code := strings.TrimSuffix(base, filepath.Ext(base))

		b, err := os.ReadFile(f)
		if err != nil {
			log.Printf("读取 locale 文件 %s 失败: %v", f, err)
			continue
		}

		var m map[string]string
		if err := json.Unmarshal(b, &m); err != nil {
			log.Printf("解析 locale 文件 %s 失败: %v", f, err)
			continue
		}

		locales[code] = m
	}

	// 构建可用语言列表，尝试从每个 locale 中读取 LanguageName
	availableLocales = []LocaleMeta{}

	for code, m := range locales {
		name := code
		if v, ok := m[KeyLanguageName]; ok && v != "" {
			name = v
		}

		availableLocales = append(availableLocales, LocaleMeta{Code: code, Name: name})
	}

	if len(locales) == 0 {
		log.Printf("警告: 未检测到 locales 中的翻译文件，默认使用内置中文。")
	}
}

// getLocale 返回指定语言的翻译 map，如不存在则回退到默认
func getLocale(code string) map[string]string {
	if code == "" {
		code = defaultLang
	}

	if m, ok := locales[code]; ok {
		return m
	}

	if m, ok := locales[defaultLang]; ok {
		return m
	}

	// 最后回退为内置中文提示
	return langZH
}

// GetAvailableLocales 返回可用语言及当前选择标识
func GetAvailableLocales(current string) []LocaleMeta {
	res := []LocaleMeta{}

	for _, l := range availableLocales {
		sel := false
		if l.Code == current {
			sel = true
		}

		res = append(res, LocaleMeta{Code: l.Code, Name: l.Name, Selected: sel})
	}

	// 如果没有从文件加载到任何语言, 则提供中文和英文两个选项
	if len(res) == 0 {
		res = []LocaleMeta{
			{Code: "zh", Name: "中文", Selected: current == "zh"},
			{Code: "en", Name: "English", Selected: current == "en"},
		}
	}

	return res
}

// EnsureLocaleExists 可用于在运行时检查并创建基础 locales
func EnsureLocaleExists() {
	// 如果 locales 目录下没有文件, 写入基础 zh.json 和 en.json
	dir := "locales"

	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		log.Printf("扫描 locale 目录失败: %v", err)
		return
	}

	if len(files) > 0 {
		return
	}

	if err = os.MkdirAll(dir, 0755); err != nil {
		log.Printf("创建 locales 目录失败: %v", err)
		return
	}

	// 写入默认中文和英文文件
	writeLocaleFile(filepath.Join(dir, "zh.json"), langZH)
	writeLocaleFile(filepath.Join(dir, "en.json"), langEN)
}

func writeLocaleFile(path string, m map[string]string) {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Printf("marshal locale %s failed: %v\n", path, err)
		return
	}

	if err = os.WriteFile(path, b, 0600); err != nil {
		fmt.Printf("write locale file %s failed: %v\n", path, err)
	}
}
