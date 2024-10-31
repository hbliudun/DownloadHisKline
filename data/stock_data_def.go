package data

// Exchange.
var (
	// Chinese
	CFFEX = "CFFEX" // China Financial Futures Exchange
	SHFE  = "SHFE"  // Shanghai Futures Exchange
	CZCE  = "CZCE"  // Zhengzhou Commodity Exchange
	DCE   = "DCE"   // Dalian Commodity Exchange
	INE   = "INE"   // Shanghai International Energy Exchange
	GFEX  = "GFEX"  // Guangzhou Futures Exchange
	SSE   = "SSE"   // Shanghai Stock Exchange
	SZSE  = "SZSE"  // Shenzhen Stock Exchange
	BSE   = "BSE"   // Beijing Stock Exchange
	SHHK  = "SHHK"  // Shanghai-HK Stock Connect
	SZHK  = "SZHK"  // Shenzhen-HK Stock Connect
	SGE   = "SGE"   // Shanghai Gold Exchange
	WXE   = "WXE"   // Wuxi Steel Exchange
	CFETS = "CFETS" // CFETS Bond Market Maker Trading System
	XBOND = "XBOND" // CFETS X-Bond Anonymous Trading System

	// Global
	SMART    = "SMART"    // Smart Router for US stocks
	NYSE     = "NYSE"     // New York Stock Exchnage
	NASDAQ   = "NASDAQ"   // Nasdaq Exchange
	ARCA     = "ARCA"     // ARCA Exchange
	EDGEA    = "EDGEA"    // Direct Edge Exchange
	ISLAND   = "ISLAND"   // Nasdaq Island ECN
	BATS     = "BATS"     // Bats Global Markets
	IEX      = "IEX"      // The Investors Exchange
	AMEX     = "AMEX"     // American Stock Exchange
	TSE      = "TSE"      // Toronto Stock Exchange
	NYMEX    = "NYMEX"    // New York Mercantile Exchange
	COMEX    = "COMEX"    // COMEX of CME
	GLOBEX   = "GLOBEX"   // Globex of CME
	IDEALPRO = "IDEALPRO" // Forex ECN of Interactive Brokers
	CME      = "CME"      // Chicago Mercantile Exchange
	ICE      = "ICE"      // Intercontinental Exchange
	SEHK     = "SEHK"     // Stock Exchange of Hong Kong
	HKFE     = "HKFE"     // Hong Kong Futures Exchange
	SGX      = "SGX"      // Singapore Global Exchange
	CBOT     = "CBT"      // Chicago Board of Trade
	CBOE     = "CBOE"     // Chicago Board Options Exchange
	CFE      = "CFE"      // CBOE Futures Exchange
	DME      = "DME"      // Dubai Mercantile Exchange
	EUREX    = "EUX"      // Eurex Exchange
	APEX     = "APEX"     // Asia Pacific Exchange
	LME      = "LME"      // London Metal Exchange
	BMD      = "BMD"      // Bursa Malaysia Derivatives
	TOCOM    = "TOCOM"    // Tokyo Commodity Exchange
	EUNX     = "EUNX"     // Euronext Exchange
	KRX      = "KRX"      // Korean Exchange
	OTC      = "OTC"      // OTC Product (Forex/CFD/Pink Sheet Equity)
	IBKRATS  = "IBKRATS"  // Paper Trading Exchange of IB

	// Special Function
	LOCAL = "LOCAL" // For local generated data
)

// vnpy 交易所映射
var mapExchangeVn2Tushare = map[string]string{
	CFFEX: "CFX",
	SHFE:  "SHF",
	CZCE:  "ZCE",
	DCE:   "DCE",
	INE:   "INE",
	SSE:   "SH",
	SZSE:  "SZ",
	BSE:   "BJ",
	GFEX:  "GFE",
}

var mapExchangeTushare2Vn = map[string]string{
	"CFX": CFFEX,
	"SHF": SHFE,
	"ZCE": CZCE,
	"DCE": DCE,
	"INE": INE,
	"SH":  SSE,
	"SZ":  SZSE,
	"BJ":  BSE,
	"GFE": GFEX,
}

func GetExchangeVn2Tushare(exchange string) string {
	return mapExchangeVn2Tushare[exchange]
}

func GetExchangeTushare2Vn(exchange string) string {
	return mapExchangeTushare2Vn[exchange]
}

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
