package models

import (
	"time"

	"github.com/markcheno/go-talib"
)

// データフレームの構造体を定義
type DataFrameCandle struct {
	ProductCode string        `json:"product_code"`
	Duration    time.Duration `json:"duration"`
	Candles     []Candle      `json:"candles"`
	Smas        []Sma         `json:"smas,omitempty"`
}

// 単純移動平均線(simple-moving-average)
type Sma struct {
	Period int       `json:"period,omitempty"` // Period: 終値, omitempty: jsonに変換した際に数値が0（または値が何も入っていない状態）の場合は省略する
	Values []float64 `json:"values,omitempty"`
}

// []Candle配列にデータを格納し、Candle Chartで表示するための設定

// 時刻情報の取得処理を定義
func (df *DataFrameCandle) Times() []time.Time {
	// candle chartで出力するための情報をスライスで定義
	s := make([]time.Time, len(df.Candles))

	// candle chartで出力する情報(時刻情報)をforで順次格納
	for i, candle := range df.Candles {
		s[i] = candle.Time
	}

	// 時刻情報を格納したスライスを出力
	return s
}

// Open情報の取得処理を定義
func (df *DataFrameCandle) Opens() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Open
	}
	return s
}

// Close情報の取得処理を定義
func (df *DataFrameCandle) Closes() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Close
	}
	return s
}

// High情報の取得処理を定義
func (df *DataFrameCandle) Highs() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.High
	}
	return s
}

// Low情報の取得処理を定義
func (df *DataFrameCandle) Lows() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Low
	}
	return s
}

// Volume情報の取得処理を定義
func (df *DataFrameCandle) Volumes() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Volume
	}
	return s
}

// 単純移動平均線(SMA)のデータ取得処理
func (df *DataFrameCandle) AddSma(period int) bool {
	if len(df.Candles) > period {
		df.Smas = append(df.Smas, Sma{
			Period: period,
			Values: talib.Sma(df.Closes(), period),
		})
		return true
	}
	return false
}
