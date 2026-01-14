//
// FilePath    : video-trim\main.go
// Author      : jiaopengzi
// Blog        : https://jiaopengzi.com
// Copyright   : Copyright (c) 2026 by jiaopengzi, All Rights Reserved.
// Description : 简单的视频裁剪工具
//

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// 读取配置文件
	readConfig()

	// 初始化目录
	initDir()

	// 初始化并加载本地化文件
	EnsureLocaleExists()
	loadLocales()

	// 路由注册
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/download/", handleDownload)
	http.HandleFunc("/clear", handleClear)

	// 打印本机局域网 IP, 方便访问
	localIP := getLocalIP()
	fmt.Printf("\n✅ Open in browser: http://%s%s\n\n", localIP, serverPort)
	fmt.Println("Ensure your browser and server are on the same LAN!")

	// 启动 HTTP 服务并设置超时以避免资源耗尽 (读取超时, 写入超时, 空闲连接超时 从配置读取，单位: 秒)
	srv := &http.Server{
		Addr:         serverPort,
		Handler:      http.DefaultServeMux,
		ReadTimeout:  time.Duration(readTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(writeTimeoutSeconds) * time.Second,
		IdleTimeout:  time.Duration(idleTimeoutSeconds) * time.Second,
	}

	log.Printf("starting server on %s", serverPort)
	log.Fatal(srv.ListenAndServe())
}
