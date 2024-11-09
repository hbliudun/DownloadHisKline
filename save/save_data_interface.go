package save

import "DownloadHisKLine/data"

type DBBase interface {
	Init() error
	//SaveDailyKLine(data any) error
	Close() error
	SaveDailyKLine(klines []*data.DailyKLineData) error
}
