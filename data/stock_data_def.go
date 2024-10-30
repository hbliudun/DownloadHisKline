package data

type AStockKey struct {
	exchange int    //市场
	code     string //code
}

// tushare http 请求包头信息
type HttpReqHead struct {
	ApiName string `json:"api_name"`
	Token   string `json:"token"`
	Params  any    `json:"params"`
	Fields  any    `json:"fields"`
}

// TushareRespPackHead tushare http 应答包头
type TushareRespPackHead struct {
	Code int                 `json:"code"`
	Msg  string              `json:"msg"`
	Data TushareRespPackData `json:"data"`
}

// TushareRespPackData tushare http 应答包体
type TushareRespPackData struct {
	Fields []string `json:"fields"`
	Items  [][]any  `json:"items"`
}

// stock_basic 信息
type StockBasicInfo struct {
	Ts_code    string `json:"ts_code"`
	Symbol     string `json:"symbol"`
	Name       string `json:"name"`
	Area       string `json:"area"`
	Industry   string `json:"industry"`
	Cnspell    string `json:"cnspell"`
	Market     string `json:"market"`
	Listdate   string `json:"list_date"`
	Actname    string `json:"act_name"`
	Actenttype string `json:"act_ent_type"`
}

type StockBasicData struct {
	Fields []string   `json:"fields"`
	Items  [][]string `json:"items"`
}

type StockInfoResp struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data StockBasicData `json:"data"`
}

type BaseInfoParam struct {
	ListStatus string `json:"list_status"`
	TsCode     string `json:"ts_code"`
}

// api daily params
type DailyParam struct {
	TsCode    string `json:"ts_code"`
	TradeDate string `json:"trade_date "`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// api daily kline
type DailyKLineData struct {
	TsCode    string  `json:"ts_code"`
	TradeDate string  `json:"trade_date"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	PreClose  float64 `json:"pre_close"`
	Change    float64 `json:"change"`
	PctChg    float64 `json:"pct_chg"`
	Vol       float64 `json:"vol"`
	Amount    float64 `json:"amount"`
}
