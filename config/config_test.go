package config

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	cfg := &Config{}
	cfg.Init("E:\\data\\code\\go\\DownloadHisKLine\\bin\\conf.json")

}
