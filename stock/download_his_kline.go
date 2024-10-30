package stock

import (
	"DownloadHisKLine/config"
	"DownloadHisKLine/data"
	"DownloadHisKLine/save"
	"log"
)

type DownLoadHisKline struct {
	client *data.TuShareHttpCliet
	Db     *save.DBMysql
}

func (dl *DownLoadHisKline) Init(config *config.Config, db *save.DBMysql) {

	dl.Db = db
	// 创建 go 协程采集历史数据
	dl.client = data.NewTuShareHttpCliet(config)
	dl.client.Init()
}

// DownloadAllHisKLine 下载所有品种历史K线
func (dl *DownLoadHisKline) DownloadAllHisKLine() (int, error) {

	// 获取全市场品种
	stocks, err := dl.client.GetAllAStockInfo()
	if err != nil {
		return 0, err
	}

	downLoadCounts := 0
	// 遍历所有品种 开始下载K线信息
	last_err := err
	for _, stock := range stocks {
		// 下载历史K线
		kLines, err := dl.client.DownloadHisKLine(stock.Ts_code, "", stock.Listdate, "")
		if err != nil {
			last_err = err
			continue
		}

		// 保存K线信息至数据库
		err = dl.Db.SaveDailyKLine(kLines)
		if err != nil {
			last_err = err
			continue
		}
		downLoadCounts++
	}
	return downLoadCounts, last_err
}

func (dl *DownLoadHisKline) ProcDownLoadAllHisKLine() {
	nums, err := dl.DownloadAllHisKLine()
	if err != nil {
		log.Println("DownloadAllHisKLine failed, err:", err)
		return
	}
	log.Println("DownloadAllHisKLine nums:", nums)
}

func (dl *DownLoadHisKline) DownloadSingleHisKLine(ts_code string) (int, error) {

	// 获取全市场品种
	// 获取全市场品种
	stocks, err := dl.client.GetSingleAStockInfo(ts_code)
	if err != nil {
		return 0, err
	}

	downLoadCounts := 0
	last_err := err
	// 遍历所有品种 开始下载K线信息
	for _, stock := range stocks {
		// 下载历史K线
		klines, err := dl.client.DownloadHisKLine(stock.Ts_code, "", stock.Listdate, "")
		if err != nil {
			last_err = err
			continue
		}

		// 保存K线信息至数据库
		log.Printf("klines:%v", klines)
		err = dl.Db.SaveDailyKLine(klines)
		if err != nil {
			last_err = err
			continue
		}
		downLoadCounts++
	}
	return downLoadCounts, last_err
}
