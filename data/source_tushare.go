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

// GetSingleAStockInfo 获取单个股票信息 tsCode: 000001.SZ
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

func (t *TuShareHttpCliet) GetAllAStockInfo() ([]*StockBasicInfo, error) {

	sendParams := &BaseInfoParam{
		ListStatus: "L",
	}

	log.Printf("sendParams:%v", sendParams)
	body, err := t.tushareHttpPost("stock_basic", sendParams, "")
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	// 解析品种列表
	return t.parseStockBasicInfoResp(body)
}

// DownloadAllHisKLine 下载所有历史K线
// paramsJson ts_code 股票代码 trade_date 交易日期 start_date 开始日期 end_date 结束日期
func (t *TuShareHttpCliet) DownloadHisKLine(tsCode string, tradeDate string, startDate string, endDate string) ([]*DailyKLineData, error) {
	jsonParams := &DailyParam{
		TsCode:    tsCode,
		TradeDate: tradeDate,
		StartDate: startDate,
		EndDate:   endDate,
	}
	body, err := t.tushareHttpPost("daily", jsonParams, "")
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}
	// 解析品种列表
	return t.parseDailyKLineResp(body)
}

func (t *TuShareHttpCliet) tushareHttpPost(api string, params any, fields string) ([]byte, error) {

	jsonBody := &HttpReqHead{
		ApiName: api,
		Token:   t.token,
		Params:  params,
		Fields:  fields,
	}

	jsonData, err := json.Marshal(jsonBody)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return nil, err
	}
	log.Printf("post json data :%s\n", string(jsonData))
	resp, err := http.Post(t.ip, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error sending request:", err)
		return nil, err
	}

	defer resp.Body.Close()
	// 处理响应
	if resp.StatusCode != http.StatusOK {
		log.Println("Error StatusCode :", resp.StatusCode)
		return nil, err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error ReadAll :", err)
		return nil, err
	}
	log.Printf("recv response body[%s]", string(bodyBytes))
	return bodyBytes, nil
}

func (t *TuShareHttpCliet) parseStockBasicInfoResp(resp []byte) ([]*StockBasicInfo, error) {
	// 解析品种列表
	data := &StockInfoResp{}
	err := json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	if data.Code != 0 {
		return nil, err
	}

	var ArrStocks []*StockBasicInfo
	for _, item := range data.Data.Items {
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

func (t *TuShareHttpCliet) parseDailyKLineResp(resp []byte) ([]*DailyKLineData, error) {
	// 解析品种列表
	data := &TushareRespPackHead{}
	err := json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	if data.Code != 0 {
		return nil, err
	}

	var ArrKlines []*DailyKLineData
	for _, item := range data.Data.Items {
		index := 0
		kline := &DailyKLineData{
			TsCode:    item[index].(string),
			TradeDate: item[index+1].(string),
			Open:      item[index+2].(float64),
			High:      item[index+3].(float64),
			Low:       item[index+4].(float64),
			Close:     item[index+5].(float64),
			PreClose:  item[index+6].(float64),
			Change:    item[index+7].(float64),
			PctChg:    item[index+8].(float64),
			Vol:       item[index+9].(float64),
			Amount:    item[index+10].(float64),
		}
		ArrKlines = append(ArrKlines, kline)
	}
	return ArrKlines, nil
}
