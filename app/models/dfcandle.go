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
	Emas        []Ema         `json:"emas,omitempty"`
	BBands      *BBands       `json:"bbands,omitempty"` // ポインタで設定(空のJSONを返却することがあるため)
}

// 単純移動平均線(simple-moving-average)
type Sma struct {
	Period int       `json:"period,omitempty"` // Period: 終値, omitempty: jsonに変換した際に数値が0（または値が何も入っていない状態）の場合は省略する
	Values []float64 `json:"values,omitempty"`
}

// 指数平滑移動平均(exponential-moving-average)
type Ema struct {
	Period int       `json:"period,omitempty"`
	Values []float64 `json:"values,omitempty"`
}

// ボリンジャーバンド(bollinger-band)
type BBands struct {
	N    int       `json:"n,omitempty"`    // 日数
	K    float64   `json:"k,omitempty"`    // 標準偏差(Σ)
	Up   []float64 `json:"up,omitempty"`   // 移動平均線 +Σ
	Mid  []float64 `json:"mid,omitempty"`  // 移動平均線
	Down []float64 `json:"down,omitempty"` // 移動平均線 -Σ
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

// 指数平滑移動平均(EMA)のデータ取得処理
func (df *DataFrameCandle) AddEma(period int) bool {
	if len(df.Candles) > period {
		df.Emas = append(df.Emas, Ema{
			Period: period,
			Values: talib.Ema(df.Closes(), period),
		})
		return true
	}
	return false
}

// ボリンジャーバンド(bollinger-band)のデータ取得処理
func (df *DataFrameCandle) AddBBands(n int, k float64) bool {
	if n <= len(df.Closes()) {
		up, mid, down := talib.BBands(df.Closes(), n, k, k, 0)
		df.BBands = &BBands{
			N:    n,
			K:    k,
			Up:   up,
			Mid:  mid,
			Down: down,
		}
		return true
	}
	return false
}
