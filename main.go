package main

import (
	"DownloadHisKLine/httpserver"
)

var (
	// 配置信息
	configFile = "E:\\data\\code\\go\\DownloadHisKLine\\bin\\conf.json"
)

func main() {

	server := httpserver.HttpDataServer{}
	server.Init(configFile)
	defer server.Close()

	server.Start()

}
