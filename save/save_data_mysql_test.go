package save

import (
	"DownloadHisKLine/config"
	"DownloadHisKLine/data"
	"testing"
)

func TestSaveDailyKLine(t *testing.T) {
	cfg := &config.Config{}
	cfg.Init("E:\\data\\code\\go\\DownloadHisKLine\\bin\\conf.json")

	db := NewDBMysql(cfg)
	err := db.Init()
	if err != nil {
		t.Errorf("Init failed: %v", err)
		return
	}
	// 保存日线数据
	kline := &data.DailyKLineData{
		TsCode:    "000001.SZSE",
		TradeDate: "20230101",
		Open:      10.0,
		High:      11.0,
		Low:       9.0,
		Close:     10.5,
		PreClose:  10.0,
		Change:    0.5,
		PctChg:    0.5,
		Vol:       1000.0,
		Amount:    10000.0,
	}
	kLines := []*data.DailyKLineData{kline}
	err = db.SaveDailyKLine(kLines)
	if err != nil {
		t.Errorf("SaveDailyKLine failed: %v", err)
		return
	}
	err = db.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
		return
	}
}

func TestGetAllAStocks(t *testing.T) {

	err := DbMysqlTest()
	if err != nil {
		t.Errorf("DbMysqlTest failed: %v", err)
		return
	}
}

func TestDBMysql_SelectDbBarOverview(t *testing.T) {
	cfg := &config.Config{}
	cfg.Init("E:\\data\\code\\go\\DownloadHisKLine\\bin\\conf.json")

	db := NewDBMysql(cfg)
	err := db.Init()
	if err != nil {
		t.Errorf("Init failed: %v", err)
		return
	}
	view, err := db.SelectDbBarOverview("000001", "SZSE", "d")
	if err != nil {
		t.Errorf("SelectDbBarOverview failed: %v", err)
		return
	}
	t.Logf("view: %v", view)

	err = db.SaveDbBarOverView(view)
	if err != nil {
		t.Errorf("SaveDbBarOverView failed: %v", err)
		return
	}
}

func TestDBMysql_QueryDbBarOverView(t *testing.T) {
	cfg := &config.Config{}
	cfg.Init("E:\\data\\code\\go\\DownloadHisKLine\\bin\\conf.json")

	db := NewDBMysql(cfg)
	err := db.Init()
	if err != nil {
		t.Errorf("Init failed: %v", err)
		return
	}
	view, err := db.QueryDbBarOverView("000001", "SZSE", "d")
	if err != nil {
		t.Errorf("SelectDbBarOverview failed: %v", err)
		return
	}
	t.Logf("view: %v", view)

}
