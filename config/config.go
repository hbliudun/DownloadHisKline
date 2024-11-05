package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	//tushare匹配
	Token   string `json:"token"`
	Address string `json:"address"`

	// 数据库配置
	DbUser string `json:"db_user"`
	DbPass string `json:"db_passwd"`
	Ip     string `json:"ip"`
	Port   int    `json:"port"`
	DbName string `json:"db_name"`

	// 启动时是否下载所有品种 所有历史K线
	DownloadAll bool `json:"download_all"`
	// 每天定时下载历史数据时间
	DownloadTime string `json:"download_time"`

	GoHttpPort   string `json:"go_http_port"`
	MaxConnPerIp int    `json:"max_conn_per_ip"`
	MaxConn      int    `json:"max_conn"`

	// 同时存储数据最大并发数
	MaxSaveDataChans int `json:"max_save_data_chans"`
	// 单个通道最大消息数量
	MaxChanSize int `json:"max_chan_size"`
}

// Init 初始化配置信息
func (c *Config) Init(cfgFile string) {
	content, err := os.ReadFile(cfgFile)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	err = json.Unmarshal(content, c)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	// Let's print the unmarshalled data!
	c.PrintConf()
}

// UpdateConf 更新配置文件
func (c *Config) UpdateConf(cfgFile string) {
	byteValue, err := json.Marshal(c)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	os.WriteFile(cfgFile, byteValue, 0644)
}

// PrintConf 打印配置信息
func (c *Config) PrintConf() {
	log.Printf("token: %s\n", c.Token)
	log.Printf("address: %s\n", c.Address)
}
