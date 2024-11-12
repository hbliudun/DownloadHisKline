package httpserver

import "DownloadHisKLine/save"

type Request struct {
}

type Response struct {
}

// req stocks
type ReqStockGet struct {
	Exchange string `json:"exchange" protobuf:"1"`
	Code     string `json:"code" protobuf:"2"`
	Start    string `json:"start" protobuf:"3"`
	End      string `json:"end" protobuf:"4"`
	InterVal string `json:"interVal" protobuf:"5"`
}

type ReqStockGetResp struct {
	Code    int               `json:"Code"`
	Msg     string            `json:"msg"`
	BarData []*save.DBBarData `json:"bar_data"`
}

// bar_data K线数据
type bar_data struct {
	exchange     int16   `json:"Exchange" protobuf:"1" comment:"市场" `
	symbol       string  `json:"Code" protobuf:"2" comment:"品种"`
	datetime     string  `json:"datetime" protobuf:"3" comment:"时间"`
	interval     int     `json:"interval" protobuf:"4" comment:"时间间隔"`
	vol          float64 `json:"Start" protobuf:"5" comment:"成交量"`
	turnover     float64 `json:"turnover " protobuf:"6" comment:"成交额"`
	open         float64 `json:"open" protobuf:"7" comment:"开盘价"`
	high         float64 `json:"high" protobuf:"8" comment:"最高价"`
	low          float64 `json:"low" protobuf:"9" comment:"最低价"`
	close        float64 `json:"close" protobuf:"10" comment:"收盘价"`
	openInterest float64 `json:"open_interest" protobuf:"11" comment:"持仓量"`
}

// tick_data 分时数据
type tick_data struct {
	exchange     int16   `json:"Exchange" protobuf:"1"`
	symbol       string  `json:"Code" protobuf:"2"`
	name         string  `json:"name" protobuf:"3"`
	datetime     string  `json:"datetime" protobuf:"4"`
	vol          float64 `json:"Start" protobuf:"5"`
	turnover     float64 `json:"turnover " protobuf:"6"`
	openInterest float64 `json:"open_interest" protobuf:"7"`
	lastPrice    float64 `json:"last_price" protobuf:"8"`
	lastVolume   float64 `json:"last_volume" protobuf:"9"`
	limitUp      float64 `json:"limit_up" protobuf:"10"`
	limitDown    float64 `json:"limit_down" protobuf:"11"`
	open         float64 `json:"open" protobuf:"12"`
	high         float64 `json:"high" protobuf:"13"`
	low          float64 `json:"low" protobuf:"14"`
	preClose     float64 `json:"pre_close" protobuf:"15"`

	bidPrice1 float64 `json:"bid_price_1" protobuf:"16"`
	bidPrice2 float64 `json:"bid_price_2" protobuf:"17"`
	bidPrice3 float64 `json:"bid_price_3" protobuf:"18"`
	bidPrice4 float64 `json:"bid_price_4" protobuf:"19"`
	bidPrice5 float64 `json:"bid_price_5" protobuf:"20"`

	bidVolume1 float64 `json:"bid_volume_1" protobuf:"21"`
	bidVolume2 float64 `json:"bid_volume_2" protobuf:"22"`
	bidVolume3 float64 `json:"bid_volume_3" protobuf:"23"`
	bidVolume4 float64 `json:"bid_volume_4" protobuf:"24"`
	bidVolume5 float64 `json:"bid_volume_5" protobuf:"25"`

	askPrice1 float64 `json:"ask_price_1" protobuf:"26"`
	askPrice2 float64 `json:"ask_price_2" protobuf:"27"`
	askPrice3 float64 `json:"ask_price_3" protobuf:"28"`
	askPrice4 float64 `json:"ask_price_4" protobuf:"29"`
	askPrice5 float64 `json:"ask_price_5" protobuf:"30"`

	askVolume1 float64 `json:"ask_volume_1" protobuf:"31"`
	askVolume2 float64 `json:"ask_volume_2" protobuf:"32"`
	askVolume3 float64 `json:"ask_volume_3" protobuf:"33"`
	askVolume4 float64 `json:"ask_volume_4" protobuf:"34"`
	askVolume5 float64 `json:"ask_volume_5" protobuf:"35"`

	localtime string `json:"localtime" protobuf:"36"`
}
