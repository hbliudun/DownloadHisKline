package main

import (
	"DownloadHisKLine/config"
	"DownloadHisKLine/httpserver"
)

var (
	// 配置信息
	configFile = "E:\\data\\code\\go\\DownloadHisKLine\\bin\\conf.json"
)

func main() {
	// 读取配置信息
	cfg := &config.Config{}
	cfg.Init(configFile)

	server := httpserver.HttpDataServer{}
	server.Init(cfg)
	server.Start()
}
