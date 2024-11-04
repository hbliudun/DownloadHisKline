package data

import (
	"DownloadHisKLine/config"
	"log"
	"testing"
)

func TestGetAllAStocks(t *testing.T) {
	cfg := &config.Config{}
	cfg.Init("E:\\data\\code\\go\\DownloadHisKLine\\bin\\conf.json")

	client := NewTuShareHttpCliet(cfg)
	client.Init()
	stocks, err := client.GetAllAStockInfo()
	if err != nil {
		log.Println("GetAllAStockInfo failed")
		return
	}

	for _, stock := range stocks {
		//t.Log(stock)
		log.Println(stock)
	}
}

func TestGetSingleAStocks(t *testing.T) {
	cfg := &config.Config{}

	cfg.Init("E:\\data\\code\\go\\DownloadHisKLine\\bin\\conf.json")
	client := NewTuShareHttpCliet(cfg)
	client.Init()
	stocks, err := client.GetSingleAStockInfo("000001.SZ")

	if err != nil {
		log.Println("GetSingleAStockInfo failed")
		return
	}
	t.Log(len(stocks))
}

func TestDownloadHisKLine(t *testing.T) {
	cfg := &config.Config{}

	cfg.Init("E:\\data\\code\\go\\DownloadHisKLine\\bin\\conf.json")
	client := NewTuShareHttpCliet(cfg)
	client.Init()

	stocks, err := client.GetSingleAStockInfo("000001.SZ")

	if err != nil {
		log.Println("GetSingleAStockInfo failed")
		return
	}

	for _, stock := range stocks {
		t.Log(stock)

		klines, _ := client.DownloadHisKLine(stock.Ts_code, "", stock.Listdate, "")
		if klines != nil {
			for _, kline := range klines {
				log.Println(kline)
			}
		}
		// 终止测试 只测试一个品种
		break
	}

}
