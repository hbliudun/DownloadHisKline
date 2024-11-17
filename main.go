package main

import (
	"DownloadHisKLine/httpserver"
	"flag"

	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
)

var (
	// 配置信息
	configFile = "bin\\conf.json"
)

// 读取 参数 初始化
func init() {
	// 初始化默认配置
	flag.StringVar(&configFile, "c", "conf.json", "配置文件路径")

	// 初始化日志
	InitLog()
}
func main() {
	flag.Parse()
	server := httpserver.HttpDataServer{}
	server.Init(configFile)
	defer server.Close()

	server.Start()

}

func InitLog() {
	logger := &lumberjack.Logger{
		Filename:   "./log.log",
		MaxSize:    200,  // 日志文件大小，单位是 MB
		MaxBackups: 10,   // 最大过期日志保留个数
		MaxAge:     28,   // 保留过期文件最大时间，单位 天
		Compress:   true, // 是否压缩日志，默认是不压缩。这里设置为true，压缩日志
	}

	log.SetOutput(logger) // logrus 设置日志的输出方式
}
