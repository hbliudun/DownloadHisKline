package save

import (
	"DownloadHisKLine/config"
	"DownloadHisKLine/data"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type DBBase interface {
	Init() error
	//SaveDailyKLine(data any) error
	Close() error
	SaveDailyKLine(klines []*data.DailyKLineData) error
}

type DBMysql struct {
	//DBBase
	config *config.Config
	db     *sql.DB

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
	for _, dayline := range klines {
		sqlStr := "insert into dbbardata(`symbol`, `exchange`, `datetime`, `interval`, `volume`, `turnover`, `open_interest`, `open_price`, `high_price`, `low_price`, `close_price`) values (?,?,?,?,?,?,?,?,?,?,?)"
		symbol := dayline.TsCode[0:6]
		exchange := data.GetExchangeTushare2Vn(dayline.TsCode[7:])
		result, err := db.db.Exec(sqlStr, symbol, exchange, dayline.TradeDate, "d", dayline.Vol, dayline.Amount, 0, dayline.Open, dayline.High, dayline.Low, dayline.Close)

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
	}

	return nil
}

func (db *DBMysql) SelectDbBarOverview(symbol string, exchange string, interval string) (*DBBarOverview, error) {
	view := &DBBarOverview{Symbol: symbol, Exchange: exchange, Interval: interval}
	//sqlStr := "select * from dbbaroverview where symbol =? and exchange =? and interval =?"
	sqlStr := "select (select count(*) from dbbardata where `interval`=? and symbol=? and exchange=? ) as count, min(datetime)as start,max(datetime)as end from `dbbardata` where `interval`=? and symbol=? and exchange=?;"
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

func (db *DBMysql) SaveDbBarOverView(view *DBBarOverview) error {

	sqlStr := "replace into dbbaroverview(`symbol`, `exchange`, `interval`, `count`, `start`, `end`) values (?,?,?,?,?,?)"
	result, err := db.db.Exec(sqlStr, view.Symbol, view.Exchange, view.Interval, view.Count, view.Start, view.End)

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
