package models

import (
	"database/sql"
	"fmt"
	"gotrading/config"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	tableNameSignalEvents = "signal_events"
)

var DbConnection *sql.DB

// product_codeと時刻を連結させたテーブルを返す処理を定義
func GetCandleTableName(productCode string, duration time.Duration) string {
	return fmt.Sprintf("%s_%s", productCode, duration) // %s_%s とすることで文字列を連結させている
}

// Databaseに対して実行するQueryと実行処理を定義
func init() {
	var err error
	DbConnection, err = sql.Open(config.Config.SQLDriver, config.Config.DbName)
	if err != nil {
		log.Fatalln(err)
	}

	// DB接続時にtableが存在しない場合は生成するQueryを定義
	cmd := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            time DATETIME PRIMARY KEY NOT NULL,
            product_code STRING,
            side STRING,
            price FLOAT,
            size FLOAT)`, tableNameSignalEvents)
	DbConnection.Exec(cmd)

	for _, duration := range config.Config.Durations {
		// tableName => ex) BTC_JPY_1m
		tableName := GetCandleTableName(config.Config.ProductCode, duration)
		c := fmt.Sprintf(`
            CREATE TABLE IF NOT EXISTS %s (
            time DATETIME PRIMARY KEY NOT NULL,
            open FLOAT,
            close FLOAT,
            high FLOAT,
            low FLOAT,
			volume FLOAT)`, tableName)
		DbConnection.Exec(c)
	}
}
