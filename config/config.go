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

	DownloadAll  bool   `json:"download_all"`
	DownloadTime string `json:"download_time"`
}

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

func (c *Config) PrintConf() {
	log.Printf("token: %s\n", c.Token)
	log.Printf("address: %s\n", c.Address)
}
