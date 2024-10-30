package stock

import (
	"DownloadHisKLine/config"
	"DownloadHisKLine/data"
	"DownloadHisKLine/save"
	"testing"
)

func TestDownloadAllHisKLine(t *testing.T) {
	cfg := &config.Config{}

	cfg.Init("E:\\data\\code\\go\\DownloadHisKLine\\bin\\conf.json")
	client := data.NewTuShareHttpCliet(cfg)
	client.Init()

	db := save.NewDBMysql(cfg)

	download := &DownLoadHisKline{client, db}

	_, err := download.DownloadAllHisKLine()
	if err != nil {
		t.Error(err)
	}
}

func TestDownloadSingleHisKLine(t *testing.T) {
	cfg := &config.Config{}
	cfg.Init("E:\\data\\code\\go\\DownloadHisKLine\\bin\\conf.json")
	db := save.NewDBMysql(cfg)
	err := db.Init()
	if err != nil {
		t.Error(err)
	}

	download := &DownLoadHisKline{}
	download.Init(cfg, db)

	_, err = download.DownloadSingleHisKLine("000001.SZ")
	if err != nil {
		t.Error(err)
	}
}
