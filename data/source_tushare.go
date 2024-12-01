package data

import (
	"DownloadHisKLine/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

/*
20240522230128-a4912f53-70a9-4539-9ea6-dbc60b641bca


import tushare as ts
# 初始化pro接口
pro = ts.pro_api('20241114201856-4f34155e-3cd2-4a42-b798-df1713e24f07')
pro._DataApi__http_url = 'http://tsapi.majors.ltd:7000'

一定要把上面的代码原封不动的替换token！！！

#常规接口
df1 = pro.daily(ts_code='000001.SZ', start_date='20180701', end_date='20180718')
#通用行情接口
df2 = ts.pro_bar(api=pro,ts_code='000001.SZ', adj='qfq', start_date='20180101', end_date='20181011)


*/

type TuShareHttpCliet struct {
	config  *config.Config
	client  *http.Client
	ip      string
	port    int
	address string

	token string
}

func NewTuShareHttpCliet(cfg *config.Config) *TuShareHttpCliet {
	return &TuShareHttpCliet{
		client: &http.Client{},
		config: cfg,
	}
}

func (t *TuShareHttpCliet) Init() {
	t.ip = t.config.Address
	t.token = t.config.Token
}

// GetSingleAStockInfo 获取单个股票基本信息 tsCode: 000001.SZ
func (t *TuShareHttpCliet) GetSingleAStockInfo(tsCode string) ([]*StockBasicInfo, error) {
	//TODO 添加单个股票信息处理
	sendParams := &BaseInfoParam{
		ListStatus: "L",
		TsCode:     tsCode,
	}

	body, err := t.tushareHttpPost("stock_basic", sendParams, "")
	if err != nil {
		return nil, err
	}
	// 解析品种列表
	return t.parseStockBasicInfoResp(body)
}

// GetAllAStockInfo 获取所有股票基本信息
func (t *TuShareHttpCliet) GetAllAStockInfo() ([]*StockBasicInfo, error) {

	sendParams := &BaseInfoParam{
		ListStatus: "L",
	}

	log.Printf("GetAllAStockInfo sendParams:%v", sendParams)
	body, err := t.tushareHttpPost("stock_basic", sendParams, "")
	if err != nil {
		fmt.Println("GetAllAStockInfo Error reading response body:", err)
		return nil, err
	}

	// 解析品种列表
	return t.parseStockBasicInfoResp(body)
}

// DownloadHisKLine 下载历史K线 tsCode:股票代码(000001.SZ) tsCode:交易日期 startDate:开始日期 endDate:结束日期
func (t *TuShareHttpCliet) DownloadHisKLine(tsCode string, tradeDate string, startDate string, endDate string) ([]*DailyKLineData, error) {
	jsonParams := &DailyParam{
		TsCode:    tsCode,
		TradeDate: tradeDate,
		StartDate: startDate,
		EndDate:   endDate,
	}
	body, err := t.tushareHttpPost("daily", jsonParams, "")
	if err != nil {
		fmt.Println("DownloadHisKLine Error reading response body:", err)
		return nil, err
	}
	// 解析品种列表
	return t.parseDailyKLineResp(body)
}

// tushareHttpPost 发送http请求
func (t *TuShareHttpCliet) tushareHttpPost(api string, params any, fields string) ([]byte, error) {

	jsonBody := &HttpReqHead{
		ApiName: api,
		Token:   t.token,
		Params:  params,
		Fields:  fields,
	}

	jsonData, err := json.Marshal(jsonBody)
	if err != nil {
		log.Println("tushareHttpPost Error marshaling JSON:", err)
		return nil, err
	}
	log.Printf("post json data :%s\n", string(jsonData))
	resp, err := http.Post(t.ip, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("tushareHttpPost Error sending request:", err)
		return nil, err
	}

	defer resp.Body.Close()
	// 处理响应
	if resp.StatusCode != http.StatusOK {
		log.Println("tushareHttpPost Error StatusCode :", resp.StatusCode)
		return nil, err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("tushareHttpPost Error ReadAll :", err)
		return nil, err
	}
	return bodyBytes, nil
}

// parseStockBasicInfoResp 解析股票基本信息
func (t *TuShareHttpCliet) parseStockBasicInfoResp(resp []byte) ([]*StockBasicInfo, error) {
	// 解析品种列表
	tdata := &StockInfoResp{}
	err := json.Unmarshal(resp, &tdata)
	if err != nil {
		return nil, err
	}

	if tdata.Code != 0 {
		return nil, ErrTushare{err: tdata.Msg}
	}

	var ArrStocks []*StockBasicInfo
	for _, item := range tdata.Data.Items {
		index := 0
		stock := &StockBasicInfo{Ts_code: item[index], Symbol: item[index+1], Name: item[index+2],
			Area: item[index+3], Industry: item[index+4], Cnspell: item[index+5],
			Market:   item[index+6],
			Listdate: item[index+7], Actname: item[index+8],
			Actenttype: item[index+9]}
		ArrStocks = append(ArrStocks, stock)
	}

	return ArrStocks, nil
}

// parseDailyKLineResp 解析股票K线数据
func (t *TuShareHttpCliet) parseDailyKLineResp(resp []byte) ([]*DailyKLineData, error) {
	// 解析品种列表
	tdata := &TushareRespPackHead{}
	err := json.Unmarshal(resp, &tdata)
	if err != nil {
		return nil, err
	}

	if tdata.Code != 0 {
		log.Println("parseDailyKLineResp Error Code :", tdata.Code, " Msg:", tdata.Msg)
		return nil, ErrTushare{err: tdata.Msg}
	}

	var ArrKlines []*DailyKLineData
	for _, item := range tdata.Data.Items {
		index := 0
		kline := &DailyKLineData{
			TsCode:    item[index].(string),
			TradeDate: item[index+1].(string),
			Open:      convertToFloat64(item[index+2]),
			High:      convertToFloat64(item[index+3]),
			Low:       convertToFloat64(item[index+4]),
			Close:     convertToFloat64(item[index+5]),
			PreClose:  convertToFloat64(item[index+6]),
			Change:    convertToFloat64(item[index+7]),
			PctChg:    convertToFloat64(item[index+8]),
			Vol:       convertToFloat64(item[index+9]),
			Amount:    convertToFloat64(item[index+10]),
		}
		ArrKlines = append(ArrKlines, kline)
	}
	return ArrKlines, nil
}

func convertToString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch value.(type) {
	case float64:
		return value.(string)
	default:
		return ""
	}
	return ""
}

func convertToFloat64(value interface{}) float64 {
	if value == nil {
		return 0.0
	}

	switch value.(type) {
	case float64:
		return value.(float64)
	}
	return 0.0
}

func convertToFloat32(value interface{}) float32 {
	if value == nil {
		return 0.0
	}

	switch value.(type) {
	case float32:
		return value.(float32)
	}
	return 0.0
}

func convertToInt(value interface{}) int {
	if value == nil {
		return 0
	}

	switch value.(type) {
	case int:
		return value.(int)
	}
	return 0
}

func convertToInt64(value interface{}) int64 {
	if value == nil {
		return 0
	}

	switch value.(type) {
	case int64:
		return value.(int64)
	}
	return 0
}
