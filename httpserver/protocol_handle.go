package httpserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (server *HttpDataServer) GetStock(c *gin.Context) {
	dataGet := &ReqStockGet{}
	resp := &ReqStockGetResp{
		Code: SERVER_OK,
		Msg:  "ok",
	}

	err := c.BindJSON(dataGet)
	if err != nil {
		resp.Code = SERVER_PARSE_PACK_ERR
		resp.Msg = err.Error()
		respError(c, http.StatusBadRequest, resp)
		return
	}

	// 根据参数读取数据
	bars, err := server.dataSave.QueryDailyKLine(dataGet.Code, dataGet.Exchange, dataGet.InterVal, dataGet.Start, dataGet.End)
	if err != nil {
		resp.Code = SERVER_QUERY_DB_ERR
		resp.Msg = err.Error()
		respError(c, http.StatusServiceUnavailable, resp)
		return
	}
	resp.BarData = bars
	RespSuccess(c, resp)
}

func RespSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func respError(c *gin.Context, code int, errMsg interface{}) {
	c.JSON(code, gin.H{
		"data": errMsg,
	})
}
