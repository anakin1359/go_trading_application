package models

import (
	"fmt"
	"gotrading/bitflyer"
	"time"
)

// Candle Stickの構造体を定義
type Candle struct {
	ProductCode string
	Duration    time.Duration
	Time        time.Time
	Open        float64
	Close       float64
	High        float64
	Low         float64
	Volume      float64
}

// Candle Stickを生成するコンストラクタを定義
func NewCandle(productCode string, duration time.Duration, timeDate time.Time, open, close, high, low, volume float64) *Candle {
	return &Candle{
		productCode,
		duration,
		timeDate,
		open,
		close,
		high,
		low,
		volume,
	}
}

// GetCandleTableName関数を利用してCandleが保存されるべきテーブルを取得
func (c *Candle) TableName() string {
	return GetCandleTableName(c.ProductCode, c.Duration)
}

// Candle Stick生成までの大まかな流れ
// 1. h, m, s単位のキャンドルスティックを生成する           => Create
// 2. 生成したキャンドルスティックをデータベースに反映させる => Save

// 1. SQLクエリを発行してテーブルにレコードを追加する処理を定義
func (c *Candle) Create() error {
	// 上記のTableName関数で取得した各columnの値を「%s」に渡してSQL Queryを生成
	cmd := fmt.Sprintf("INSERT INTO %s (time, open, close, high, low, volume) VALUES (?, ?, ?, ?, ?, ?)", c.TableName())

	// sql実行後の処理結果は不要なため「_」で破棄し、エラーが発生した時のみエラーハンドリングの処理に渡せるように定義
	_, err := DbConnection.Exec(cmd, c.Time.Format(time.RFC3339), c.Open, c.Close, c.High, c.Low, c.Volume)
	if err != nil {
		return err
	}
	return err
}

// 2. 上記SQLクエリでINSERTしたレコードを保存する処理を定義
func (c *Candle) Save() error {
	// TableName関数で取得した内容を「%s」に渡して、timeで取得した単位のレコード情報に書き換える
	cmd := fmt.Sprintf("UPDATE %s SET open = ?, close = ?, high = ?, low = ?, volume = ? WHERE time = ?", c.TableName())
	_, err := DbConnection.Exec(cmd, c.Open, c.Close, c.High, c.Low, c.Volume, c.Time.Format(time.RFC3339))
	if err != nil {
		return err
	}
	return err
}

// 3. product-code, 時刻形式(秒,分,時間), 時刻を引数に渡して、SELECT文で保存したレコード情報を取得する処理を定義
func GetCandle(productCode string, duration time.Duration, dateTime time.Time) *Candle {
	tableName := GetCandleTableName(productCode, duration)
	cmd := fmt.Sprintf("SELECT time, open, close, high, low, volume FROM  %s WHERE time = ?", tableName)
	row := DbConnection.QueryRow(cmd, dateTime.Format(time.RFC3339)) // QueryRowはマッチしたレコードを1行出力するメソッド
	var candle Candle
	err := row.Scan(&candle.Time, &candle.Open, &candle.Close, &candle.High, &candle.Low, &candle.Volume)
	if err != nil {
		return nil
	}
	return NewCandle(productCode, duration, candle.Time, candle.Open, candle.Close, candle.High, candle.Low, candle.Volume)
}

func CreateCandleWithDuration(ticker bitflyer.Ticker, productCode string, duration time.Duration) bool {
	currentCandle := GetCandle(productCode, duration, ticker.TruncateDateTime(duration))
	price := ticker.GetMidPrice()
	if currentCandle == nil {
		candle := NewCandle(productCode, duration, ticker.TruncateDateTime(duration),
			price, price, price, price, ticker.Volume)
		candle.Create()
		return true
	}

	if currentCandle.High <= price {
		currentCandle.High = price
	} else if currentCandle.Low >= price {
		currentCandle.Low = price
	}
	currentCandle.Volume += ticker.Volume
	currentCandle.Close = price
	currentCandle.Save()
	return false
}

// dfcandle.goで定義した情報を全て取得してデータを成形する処理を定義
func GetAllCandle(productCode string, duration time.Duration, limit int) (dfCandle *DataFrameCandle, err error) {
	// base.goで定義したproductCodeと時刻情報を連結したテーブル名を取得する関数を使用してテーブル名を定義
	tableName := GetCandleTableName(productCode, duration)

	// SQLクエリ(tableNameのテーブルから指定した情報を降順に並び替え、且つ取得上限ありで取得し、取得した情報を昇順で並び変えて表示)
	cmd := fmt.Sprintf(`SELECT * FROM (
		SELECT time, open, close, high, low, volume FROM %s ORDER BY time DESC LIMIT ?
		) ORDER BY time ASC;`, tableName)
	rows, err := DbConnection.Query(cmd, limit)
	if err != nil {
		return
	}

	// 処理の最後に必ずデータベースへの接続を切断
	defer rows.Close()

	// DataFrameCandleポインタを呼び出してdfCandleに格納、productCodeとdurationを格納
	dfCandle = &DataFrameCandle{}
	dfCandle.ProductCode = productCode
	dfCandle.Duration = duration

	// Nextメソッドを使用してデータベース接続処理を行う(処理が成功したらTrue、失敗したらFalseを返す)
	for rows.Next() {
		var candle Candle
		candle.ProductCode = productCode
		candle.Duration = duration

		// Scanメソッドを使用して各項目の情報をデータベースから1行ずつ取得
		rows.Scan(&candle.Time, &candle.Open, &candle.Close, &candle.High, &candle.Low, &candle.Volume)
		dfCandle.Candles = append(dfCandle.Candles, candle)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	return dfCandle, nil
}
