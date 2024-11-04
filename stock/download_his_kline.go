package stock

import (
	"DownloadHisKLine/config"
	"DownloadHisKLine/data"
	"DownloadHisKLine/save"
	"log"
	"time"
)

type DownLoadHisKline struct {
	conf        *config.Config
	client      *data.TuShareHttpCliet
	Db          *save.DBMysql
	stocks      chan *data.StockBasicInfo
	stocksDaily chan *data.StockBasicInfo
}

func (dl *DownLoadHisKline) Init(config *config.Config, db *save.DBMysql) {
	dl.conf = config
	// 初始化通道数8个
	dl.stocks = make(chan *data.StockBasicInfo, config.MaxSaveDataChans)
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

	// 启动go routine
	for i := 0; i < dl.conf.MaxSaveDataChans; i++ {
		go dl.handleSaveKLineToDb()
	}

	downLoadCounts := 0
	// 遍历所有品种 开始下载K线信息

	last_err := err
	for _, stock := range stocks {
		dl.stocks <- stock
		downLoadCounts++
	}
	return downLoadCounts, last_err
}

func (dl *DownLoadHisKline) handleSaveKLineToDb() {
	select {
	case <-dl.stocks:
		// 下载历史K线
		stock := <-dl.stocks
		// 下载历史K线
		kLines, err := dl.client.DownloadHisKLine(stock.Ts_code, "", stock.Listdate, "")
		if err != nil {
			log.Println("DownloadHisKLine failed, err:", err)
		}
		// 保存K线信息至数据库
		err = dl.Db.SaveDailyKLine(kLines)
		if err != nil {
			log.Println("SaveDailyKLine failed, err:", err)
		}
	}
}

func (dl *DownLoadHisKline) handleSaveKLineToDbDaily() {
	select {
	case <-dl.stocksDaily:
		// 下载历史K线
		stock := <-dl.stocksDaily
		// 查询已有数据最新日期
		startDate := stock.Listdate
		endDate := ""
		// 下载历史K线
		kLines, err := dl.client.DownloadHisKLine(stock.Ts_code, "", startDate, endDate)
		if err != nil {
			log.Println("DownloadHisKLine failed, err:", err)
		}
		// 保存K线信息至数据库
		err = dl.Db.SaveDailyKLine(kLines)
		if err != nil {
			log.Println("SaveDailyKLine failed, err:", err)
		}
	}
}

// 获取所有历史K线信息
func (dl *DownLoadHisKline) ProcDownLoadAllHisKLine() {
	nums, err := dl.DownloadAllHisKLine()
	if err != nil {
		log.Println("DownloadAllHisKLine failed, err:", err)
		return
	}
	log.Println("DownloadAllHisKLine nums:", nums)
}

// 每天定时获取最新K线信息
func (dl *DownLoadHisKline) ProcDownloadDaily() {
	//todo 每天盘后定时下载最新数据
	lastDate := time.Now().Format("2006-01-02")

	for {
		curDate := time.Now().Format("2006-01-02")
		if lastDate != curDate {
			// 获取全市场品种
			stocks, err := dl.client.GetAllAStockInfo()
			if err != nil {
				log.Println()
				break
			}

			// 启动go routine
			for i := 0; i < dl.conf.MaxSaveDataChans; i++ {
				go dl.handleSaveKLineToDbDaily()
			}

			downLoadCounts := 0
			// 遍历所有品种 开始下载K线信息
			for _, stock := range stocks {
				dl.stocks <- stock
				downLoadCounts++
			}

			lastDate = curDate
			log.Println("update daily date ok")
		} else {
			time.Sleep(10 * time.Second)
		}
	}

}

func (dl *DownLoadHisKline) DownloadSingleHisKLine(ts_code string) (int, error) {
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
