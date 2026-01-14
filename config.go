//
// FilePath    : video-trim\config.go
// Author      : jiaopengzi
// Blog        : https://jiaopengzi.com
// Copyright   : Copyright (c) 2026 by jiaopengzi, All Rights Reserved.
// Description : 配置文件读取
//

package main

import (
	"log"

	"github.com/spf13/viper"
)

// 配置键
const (
	keyUploadDir           = "upload_dir"            // 上传目录
	keyOutputDir           = "output_dir"            // 输出目录
	keyHeadTrimSeconds     = "head_trim_seconds"     // 掐头秒数
	keyTailSeconds         = "tail_seconds"          // 去尾秒数
	keyServerPort          = "server_port"           // 服务器端口
	keyMaxUploadSize       = "max_upload_size"       // 单个最大上传大小
	keyReadTimeoutSeconds  = "read_timeout_seconds"  // 读取超时秒数
	keyWriteTimeoutSeconds = "write_timeout_seconds" // 写入超时秒数
	keyIdleTimeoutSeconds  = "idle_timeout_seconds"  // 空闲连接超时秒数
)

// 可配置变量(会被 config.yaml 覆盖)
var (
	uploadDir       = "./uploads" // 上传文件存放目录
	outputDir       = "./outputs" // 输出文件存放目录
	headTrimSeconds = 6           // 掐头:多少秒
	tailSeconds     = 0           // 去尾:多少秒
	serverPort      = ":7778"     // 服务器监听端口
	// 超时配置(秒)
	readTimeoutSeconds        = 15         // 读取超时(秒)
	writeTimeoutSeconds       = 60         // 写入超时(秒)
	idleTimeoutSeconds        = 120        // 空闲连接超时(秒)
	maxUploadSize       int64 = 2048 << 20 // 2048 MB, 单个文件最大允许上传大小(字节)
)

// 读取配置文件(如果存在)
// readConfig 读取 workspace 根目录下的 config.yaml(如果存在),
// 仅使用标准库解析简单的 key: value 配置行, 支持注释和引号。
func readConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// defaults
	viper.SetDefault(keyUploadDir, uploadDir)
	viper.SetDefault(keyOutputDir, outputDir)
	viper.SetDefault(keyHeadTrimSeconds, headTrimSeconds)
	viper.SetDefault(keyTailSeconds, tailSeconds)
	viper.SetDefault(keyServerPort, serverPort)
	viper.SetDefault(keyMaxUploadSize, maxUploadSize)
	viper.SetDefault(keyReadTimeoutSeconds, readTimeoutSeconds)
	viper.SetDefault(keyWriteTimeoutSeconds, writeTimeoutSeconds)
	viper.SetDefault(keyIdleTimeoutSeconds, idleTimeoutSeconds)

	if err := viper.ReadInConfig(); err != nil {
		// 如果配置文件不存在则使用默认值
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return
		}

		log.Printf("读取 config.yaml 失败: %v", err)

		return
	}

	// 覆盖变量
	if v := viper.GetString(keyUploadDir); v != "" {
		uploadDir = v
	}

	if v := viper.GetString(keyOutputDir); v != "" {
		outputDir = v
	}

	if v := viper.GetInt(keyHeadTrimSeconds); v >= 0 {
		headTrimSeconds = v
	}

	if v := viper.GetInt(keyTailSeconds); v >= 0 {
		tailSeconds = v
	}

	if v := viper.GetString(keyServerPort); v != "" {
		serverPort = v
	}

	if v := viper.GetInt64(keyMaxUploadSize); v > 0 {
		maxUploadSize = v
	}

	if v := viper.GetInt(keyReadTimeoutSeconds); v >= 0 {
		readTimeoutSeconds = v
	}

	if v := viper.GetInt(keyWriteTimeoutSeconds); v >= 0 {
		writeTimeoutSeconds = v
	}

	if v := viper.GetInt(keyIdleTimeoutSeconds); v >= 0 {
		idleTimeoutSeconds = v
	}
}
