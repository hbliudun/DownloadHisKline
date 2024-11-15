package httpserver

import (
	"DownloadHisKLine/config"
	"DownloadHisKLine/save"
	"DownloadHisKLine/stock"
	"github.com/gin-gonic/gin"
	"log"
)

type HttpDataServer struct {
	cfgFile        string
	config         *config.Config
	dataDownload   *stock.DownLoadHisKline
	dataSave       *save.DBMysql
	gs             *gin.Engine
	lastUpdateDate string
}

func (server *HttpDataServer) Init(cfgFile string) {
	server.cfgFile = cfgFile
	// 读取配置信息
	server.config = &config.Config{}
	server.config.Init(server.cfgFile)

	//初始化数据库连接
	server.dataSave = save.NewDBMysql(server.config)
	err := server.dataSave.Init()
	if err != nil {
		log.Printf("db init failed")
		return
	}

	//初始化数据下载服务
	server.dataDownload = &stock.DownLoadHisKline{}
	server.dataDownload.Init(server.config, server.dataSave)

	// gin http服务 用于查询数据
	server.gs = gin.Default()
	server.gs.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong!",
		})
	})

	// 查询股票信息
	server.gs.GET("/stock", server.GetStock)
	// 更新股票信息
	server.gs.POST("/stock", func(c *gin.Context) {

	})

}

func (server *HttpDataServer) Start() {
	// 启动http数据查询服务
	server.gs.Run(server.config.GoHttpPort)
}

func (server *HttpDataServer) Close() {
	if server.dataSave != nil {
		server.dataSave.Close()
	}

	if server.config != nil {
		server.config.UpdateConf(server.cfgFile)
	}
}
