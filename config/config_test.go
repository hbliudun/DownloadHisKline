package config

import (
	"testing"
)

var cfgPath = "E:\\data\\code\\go\\DownloadHisKLine\\bin\\conf.json"

// TestReadConfig 测试读取配置文件
func TestReadConfig(t *testing.T) {
	cfg := &Config{}
	cfg.Init(cfgPath)

}

// TestUpdateFile 测试更新配置文件
func TestUpdateFile(t *testing.T) {
	cfg := &Config{}
	cfg.Init(cfgPath)

	cfg.DownloadAll = true
	cfg.GoHttpPort = "18081"

	cfg.UpdateConf(cfgPath)

}
