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
	stocksDaily chan *data.StockBasicInfo
	barsSaveDb  chan []*data.DailyKLineData

	// 上次更新日期
	lastUpdateDate string
}

func (dl *DownLoadHisKline) Init(config *config.Config, db *save.DBMysql) {
	dl.conf = config
	dl.stocksDaily = make(chan *data.StockBasicInfo, config.MaxChanSize)
	dl.barsSaveDb = make(chan []*data.DailyKLineData, config.MaxChanSize*config.MaxSaveDataChans)
	dl.Db = db
	// 创建 go 协程采集历史数据
	dl.client = data.NewTuShareHttpCliet(config)
	dl.client.Init()
	dl.lastUpdateDate = "2006-01-02"

	// 启动databar 入库线程
	go dl.handleDbBarData()

	// 每天定时进行当天数据采集服务
	go dl.ProcDownloadDaily()
}

// 每日更新K线数据业务处理线程
func (dl *DownLoadHisKline) handleSaveKLineToDbDaily() {

	defer func() { log.Println("handleSaveKLineToDbDaily exit") }()

	for {
		select {
		case stock := <-dl.stocksDaily:
			// 下载历史K线
			endDate := time.Now().Format("20060102")

			// 查询数据库该品种已有数据最新日期
			symbol := stock.Ts_code[0:6]
			exchange := data.GetExchangeTushare2Vn(stock.Ts_code[7:])
			view, err := dl.Db.SelectDbBarOverview(symbol, exchange, "d")
			if err != nil {
				log.Println("SelectDbBarOverview failed, err:", err)
				break
			}

			var startDate string
			if view.Count == 0 {
				startDate = "19000101"
			} else {
				sTime, err := time.Parse("2006-01-02 15:04:05", view.End)
				if err == nil {
					startDate = sTime.Format("20060102")
				} else {
					startDate = "19000101"
				}
			}

			if startDate == endDate {
				log.Printf("data is newest  ")
				break
			}

			// 下载历史K线
			//time.Sleep(1500 * time.Millisecond)
			kLines, err := dl.client.DownloadHisKLine(stock.Ts_code, "", startDate, endDate)
			if err != nil {
				log.Println("DownloadHisKLine failed, err:", err)
			}
			log.Printf("DownloadHisKLine %s,counts:%d ", stock.Ts_code, len(kLines))
			dl.barsSaveDb <- kLines

		//当前通道无数据时，等待30秒无数据则退出
		case <-time.After(30 * time.Second):
			if len(dl.stocksDaily) == 0 {
				log.Println("handleSaveKLineToDbDaily 等待30s 无数据 exit")
				return
			}
		}
	}

}

// handleDbBarData databar数据入库处理线程
func (dl *DownLoadHisKline) handleDbBarData() {
	for {
		select {
		case bars := <-dl.barsSaveDb:
			// 保存K线信息至数据库
			err := dl.Db.SaveDailyKLine(bars)
			if err != nil {
				log.Println("SaveDailyKLine failed, err:", err)
			}
			if len(bars) > 0 {
				TsCode := bars[0].TsCode
				symbol := TsCode[0:6]
				exchange := data.GetExchangeTushare2Vn(TsCode[7:])

				// 在databar 统计 databarview
				view, err := dl.Db.SelectDbBarOverview(symbol, exchange, "d")
				if err != nil {
					log.Println("SelectDbBarOverview failed, err:", err)
					break
				}
				// 保存品种历史数据统计信息
				err = dl.Db.SaveDbBarOverView(view)
				if err != nil {
					log.Println("SaveDbBarOverView failed, err:", err)
				}
			}
		}
	}
}

// 每天定时获取最新K线信息
func (dl *DownLoadHisKline) ProcDownloadDaily() {
	//todo 每天盘后定时下载最新数据
	lastDate := dl.lastUpdateDate
	downloadTime, _ := time.Parse("15:04:05", "17:00:00")

	for {
		curDate := time.Now().Format("2006-01-02")
		curTime := time.Now()
		// 到定时采集数据时间
		if lastDate != curDate && downloadTime.Before(curTime) {
			// 获取全市场品种
			stocks, err := dl.client.GetAllAStockInfo()
			if err != nil {
				log.Println(err)
				break
			}

			// 启动 多个routine 读取k线数据
			for i := 0; i < dl.conf.MaxSaveDataChans; i++ {
				go dl.handleSaveKLineToDbDaily()
			}
			// 所有品种放入chan 用于K线数据采集
			for _, stock := range stocks {
				dl.stocksDaily <- stock
			}
			lastDate = curDate
			dl.lastUpdateDate = lastDate

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
