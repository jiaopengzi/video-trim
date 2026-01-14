//
// FilePath    : video-trim\i18n_lang_zh.go
// Author      : jiaopengzi
// Blog        : https://jiaopengzi.com
// Copyright   : Copyright (c) 2026 by jiaopengzi, All Rights Reserved.
// Description : 中文翻译
//

package main

var langZH = map[string]string{
	KeyLanguageName:          "中文",
	KeyTitle:                 "视频裁剪工具",
	KeyHeaderUpload:          "上传视频(裁剪前/后 N 秒)",
	KeyChooseVideo:           "选择视频",
	KeyUploadButton:          "上传并处理",
	KeyUploadingText:         "正在上传并处理...",
	KeyClearButton:           "清理已上传与输出文件",
	KeyConfirmClear:          "确认清理所有已上传和输出文件吗？此操作不可恢复。",
	KeySelectAtLeastOne:      "请选择至少一个视频文件后再上传",
	KeyFileTooLargePrefix:    "文件 \"",
	KeyFileTooLargeSuffix:    "\" 超过单文件允许大小 ",
	KeyFileTooLargeEnd:       ", 请减少文件大小后重试。",
	KeyHeadLabel:             "掐头 N 秒(可修改)",
	KeyTailLabel:             "掐尾 N 秒(可修改, 默认 0)",
	KeyHint:                  "处理完成后会自动跳转到下载页面；确保浏览器和当前服务端在同一局域网。",
	KeyProcessedTitle:        "处理完成, 点击下载: ",
	KeyDownloadAll:           "下载全部",
	KeyDownload:              "下载",
	KeyReturnUpload:          "返回上传页面",
	KeyRemove:                "移除",
	KeyAlertNoTrim:           "裁剪开头和结尾均为 0, 无需处理",
	KeyNoProcessedFilesHint:  "没有文件被成功处理, 请检查源文件或 FFmpeg 日志。",
	KeyRequestBodyTooLarge:   "文件太大, 最大允许上传大小为 %s。请减少文件大小后重试。",
	KeyRequestParseError:     "请求体太大或无法解析表单",
	KeyCannotReadFile:        "无法读取文件 %s",
	KeyFileEmptyOrUnreadable: "文件 %s 为空或无法读取",
	KeyNotSupportedVideo:     "文件 %s 不是受支持的视频格式(魔法数字校验失败)",
}
