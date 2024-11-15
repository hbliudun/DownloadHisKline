package save

import (
	"DownloadHisKLine/config"
	"DownloadHisKLine/data"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
	"time"
)

type DBMysql struct {
	//DBBase
	config *config.Config
	db     *sql.DB
	rwLock sync.RWMutex

	DbUser string
	DbPass string
	Ip     string
	Port   int
	DbName string
}

// NewDBMysql 创建数据库连接对象
func NewDBMysql(conf *config.Config) *DBMysql {
	return &DBMysql{
		config: conf,
	}
}

// Init 初始化数据库连接
func (db *DBMysql) Init() error {

	// 读取配置信息
	address := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", db.config.DbUser, db.config.DbPass, db.config.Ip, db.config.Port, db.config.DbName)
	// 初始化数据库连接
	Db, err := sql.Open("mysql", address)
	Db.SetMaxOpenConns(64)
	if err != nil {
		return err
	}

	db.db = Db
	return nil
}

// SaveDailyKLine 保存日线数据
func (db *DBMysql) SaveDailyKLine(klines []*data.DailyKLineData) error {

	sqlStr := "replace into dbbardata(`symbol`, `exchange`, `datetime`, `interval`, `volume`, `turnover`, `open_interest`, `open_price`, `high_price`, `low_price`, `close_price`) values"

	var values string
	for i, dayline := range klines {
		symbol := dayline.TsCode[0:6]
		exchange := data.GetExchangeTushare2Vn(dayline.TsCode[7:])

		val := fmt.Sprintf("('%s','%s','%s','%s',%f,%f,%d,%f,%f,%f,%f)", symbol, exchange, dayline.TradeDate, "d", dayline.Vol, dayline.Amount, 0, dayline.Open, dayline.High, dayline.Low, dayline.Close)
		values = values + val
		if i == len(klines)-1 {
			values = values + ";"
		} else {
			values = values + ","
		}
	}
	sqlStr = sqlStr + values
	db.rwLock.Lock()
	defer db.rwLock.Unlock()
	result, err := db.db.Exec(sqlStr)

	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

// 查询 dabaBar
func (db *DBMysql) QueryDailyKLine(symbol string, exchange string, interval string, startDate string, endDate string) ([]*DBBarData, error) {
	startTime, err := time.Parse("20060102", startDate)
	if err != nil {
		return nil, err
	}
	endTime, err := time.Parse("20060102", endDate)
	if err != nil {
		return nil, err
	}
	start := startTime.Format("2006-01-02 00:00:00")
	end := endTime.Format("2006-01-02 00:00:00")

	sqlStr := "select `datetime`,`open_price`,`high_price`,`low_price`,`close_price`,`volume`,`turnover`,`open_interest` from `dbbardata` where `interval`=? and symbol=? and exchange=? and datetime>=? and datetime<=? limit 0,10;"

	db.rwLock.RLock()
	defer db.rwLock.RUnlock()
	rows, err := db.db.Query(sqlStr, interval, symbol, exchange, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bars []*DBBarData
	for rows.Next() {
		bar := &DBBarData{Interval: interval, Symbol: symbol, Exchange: exchange}
		err = rows.Scan(&bar.Datetime, &bar.Open, &bar.High, &bar.Low, &bar.Close, &bar.Volume, &bar.Turnover, &bar.OpenInterest)
		bars = append(bars, bar)
	}
	return bars, nil
}

func (db *DBMysql) SelectDbBarOverview(symbol string, exchange string, interval string) (*DBBarOverview, error) {
	view := &DBBarOverview{Symbol: symbol, Exchange: exchange, Interval: interval}
	//sqlStr := "select * from dbbaroverview where symbol =? and exchange =? and interval =?"
	sqlStr := "select (select count(*) from dbbardata where `interval`=? and symbol=? and exchange=? ) as count, ifnull(min(datetime),'0001-01-01 00:00:00') as start,ifnull(max(datetime),'0001-01-01 00:00:00')as end from `dbbardata` where `interval`=? and symbol=? and exchange=?;"

	db.rwLock.RLock()
	defer db.rwLock.RUnlock()
	rows, err := db.db.Query(sqlStr, interval, symbol, exchange, interval, symbol, exchange)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&view.Count, &view.Start, &view.End)
		if err != nil {
			return nil, err
		}
		break
	}
	return view, nil
}

// 更新品种统计信息
func (db *DBMysql) SaveDbBarOverView(view *DBBarOverview) error {

	sqlStr := "replace into dbbaroverview(`symbol`, `exchange`, `interval`, `count`, `start`, `end`) values (?,?,?,?,?,?)"

	db.rwLock.Lock()
	defer db.rwLock.Unlock()
	result, err := db.db.Exec(sqlStr, view.Symbol, view.Exchange, view.Interval, view.Count, view.Start, view.End)

	if err != nil {
		return err
		// insert into failed , update
		//sqlStr = "update dbbaroverview set `count`=?, `start`=?, `end`=? where `symbol`=? and `exchange`=? and `interval`=?"
		//result, err = db.db.Exec(sqlStr, view.Count, view.Start, view.End, view.Symbol, view.Exchange, view.Interval)
		//if err != nil {
		//	return err
		//}
	} else {
		_, err = result.LastInsertId()
		if err != nil {
			return err
		}
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

// QueryDbBarOverView 查询品种统计信息
func (db *DBMysql) QueryDbBarOverView(symbol string, exchange string, interval string) (*DBBarOverview, error) {
	view := &DBBarOverview{Symbol: symbol, Exchange: exchange, Interval: interval}
	// 不要select * 速度慢
	sqlStr := "select `count`,`start`,`end` from `dbbaroverview` where `interval`=? and symbol=? and exchange=?;"

	db.rwLock.RLock()
	defer db.rwLock.RUnlock()
	rows, err := db.db.Query(sqlStr, interval, symbol, exchange, interval, symbol, exchange)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&view.Count, &view.Start, &view.End)
		if err != nil {
			return nil, err
		}
		break
	}
	return view, nil
}

// Close 关闭数据库连接
func (db *DBMysql) Close() error {
	return db.db.Close()
}

func DbMysqlTest() error {
	// 初始化数据库连接
	Db, err := sql.Open("mysql", "root:zth123456.@tcp(192.168.3.6:13306)/vnpy-test")
	if err != nil {
		log.Println("Error opening database:", err)
		return err
	}
	// 关闭数据库连接
	defer Db.Close()

	// 执行insert
	sqlStr := "insert into dbbardata(`symbol`, `exchange`, `datetime`, `interval`, `volume`, `turnover`, `open_interest`, `open_price`, `high_price`, `low_price`, `close_price`)" +
		"values (?,?,?,?,?,?,?,?,?,?,?)"

	result, err := Db.Exec(sqlStr, "000002", "SZ", "20230101", "d", 1000, 10000, 0, 10.0, 11.0, 9.0, 10.5)
	if err != nil {
		log.Println("Error Exec:", err)
		return err
	}

	// 获取插入数据的ID
	_, err = result.LastInsertId()
	if err != nil {
		log.Println("Error LastInsertId:", err)
		return err
	}

	// 获取受影响的行数
	_, err = result.RowsAffected()
	if err != nil {
		log.Println("Error RowsAffected:", err)
		return err
	}
	return nil
}
