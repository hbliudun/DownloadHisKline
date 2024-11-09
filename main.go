package main

import (
	"DownloadHisKLine/httpserver"
	"flag"
)

var (
	// 配置信息
	configFile = "bin\\conf.json"
)

// 读取 参数 初始化
func init() {
	flag.StringVar(&configFile, "c", "conf.json", "配置文件路径")
}
func main() {
	flag.Parse()
	server := httpserver.HttpDataServer{}
	server.Init(configFile)
	defer server.Close()

	server.Start()

}
