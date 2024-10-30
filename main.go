package main

import (
	"DownloadHisKLine/config"
	"DownloadHisKLine/save"
	"DownloadHisKLine/stock"
	"log"
)

var (
	// 配置信息
	configFile = "config.json"
)

func main() {
	// 读取配置信息
	cfg := &config.Config{}
	cfg.Init(configFile)

	//初始化数据库连接
	db := save.NewDBMysql(cfg)
	err := db.Init()
	if err != nil {
		log.Printf("db init failed")
		return
	}

	download := &stock.DownLoadHisKline{}
	download.Init(cfg, db)

	// 1.进行历史数据采集服务
	if cfg.DownloadAll {
		go download.ProcDownLoadAllHisKLine()
	}
	// 2.提供历史数据查询服务

}
