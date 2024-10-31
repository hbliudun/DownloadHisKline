package httpserver

import (
	"DownloadHisKLine/config"
	"DownloadHisKLine/save"
	"DownloadHisKLine/stock"
	"github.com/gin-gonic/gin"
	"log"
)

type HttpDataServer struct {
	config       *config.Config
	dataDownload *stock.DownLoadHisKline
	dataSave     *save.DBMysql
	gs           *gin.Engine
}

func (server *HttpDataServer) Init(config *config.Config) {
	server.config = config

	//初始化数据库连接
	server.dataSave = save.NewDBMysql(config)
	err := server.dataSave.Init()
	if err != nil {
		log.Printf("db init failed")
		return
	}

	//初始化数据下载服务
	server.dataDownload = &stock.DownLoadHisKline{}
	server.dataDownload.Init(config, server.dataSave)

	// gin http服务 用于查询数据
	server.gs = gin.Default()
	server.gs.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong!",
		})
	})
	//
	server.gs.GET("/query_stock", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "download ok!",
		})
	})

}

func (server *HttpDataServer) Start() {

	// 判断是否需要进行历史数据采集服务
	if server.config.DownloadAll {
		go server.dataDownload.ProcDownLoadAllHisKLine()
		server.config.DownloadAll = false
	}

	// todo 每天定时进行当天数据采集服务

	// 启动http数据查询服务
	server.gs.Run(server.config.GoHttpPort)
}
