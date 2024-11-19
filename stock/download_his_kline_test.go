package stock

import (
	"DownloadHisKLine/config"
	"DownloadHisKLine/save"
	"testing"
	"time"
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

func TestTimeBefore(t *testing.T) {
	downloadTime, _ := time.Parse("15:04:05", "17:00:00")
	dlTime := GetTimeInt(downloadTime)
	curTime := GetTimeInt(time.Now())
	ret := dlTime >= curTime
	t.Logf("ret: %v", ret)
}
