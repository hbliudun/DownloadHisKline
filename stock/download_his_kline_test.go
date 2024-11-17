package stock

import (
	"DownloadHisKLine/config"
	"DownloadHisKLine/save"
	"testing"
)

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
