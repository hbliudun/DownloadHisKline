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

	// 上次更新日期
	lastUpdateDate string
}

func (dl *DownLoadHisKline) Init(config *config.Config, db *save.DBMysql) {
	dl.conf = config
	dl.stocks = make(chan *data.StockBasicInfo, config.MaxChanSize)
	dl.stocksDaily = make(chan *data.StockBasicInfo, config.MaxChanSize)
	dl.Db = db
	// 创建 go 协程采集历史数据
	dl.client = data.NewTuShareHttpCliet(config)
	dl.client.Init()
	dl.lastUpdateDate = "2006-01-02"
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
	lastErr := err
	for _, stock := range stocks {
		dl.stocks <- stock
		downLoadCounts++
	}
	return downLoadCounts, lastErr
}

func (dl *DownLoadHisKline) handleSaveKLineToDb() {

	defer func() { log.Println("handleSaveKLineToDb exit") }()
	for {
		select {
		case stock := <-dl.stocks:
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
			//todo 保存品种历史数据统计信息

		//当前通道无数据时，等待30秒无数据则退出
		case <-time.After(30 * time.Second):
			if len(dl.stocks) == 0 {
				log.Println("handleSaveKLineToDb 等待10s 无数据 exit")
				return
			}
		}
	}
}

func (dl *DownLoadHisKline) handleSaveKLineToDbDaily() {
	endDate := time.Now().Format("2006-01-02")
	defer func() { log.Println("handleSaveKLineToDbDaily exit") }()

	for {
		select {
		case <-dl.stocksDaily:
			// 下载历史K线
			stock := <-dl.stocksDaily
			//todo 查询数据库该品种已有数据最新日期
			startDate := endDate
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
		//当前通道无数据时，等待30秒无数据则退出
		case <-time.After(30 * time.Second):
			if len(dl.stocks) == 0 {
				log.Println("handleSaveKLineToDbDaily 等待30s 无数据 exit")
				return
			}
		}
	}

}

// 获取所有历史K线信息
func (dl *DownLoadHisKline) ProcDownLoadAllHisKLine() {
	dl.lastUpdateDate = time.Now().Format("2006-01-02")
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
	lastDate := dl.lastUpdateDate
	downloadTime, _ := time.Parse("15:04:05", "17:00:00")

	for {
		curDate := time.Now().Format("2006-01-02")
		curTime := time.Now()
		if lastDate != curDate && downloadTime.Before(curTime) {
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
				dl.stocksDaily <- stock
				downLoadCounts++
			}
			lastDate = curDate

			log.Println("update daily date ok, downLoadCounts:", downLoadCounts)
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
