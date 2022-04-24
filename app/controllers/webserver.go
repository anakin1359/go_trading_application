package controllers

import (
	"fmt"
	"gotrading/app/models"
	"gotrading/config"
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("./app/views/google.html"))

// Viewを表示する関数を定義
func viewChartHandler(w http.ResponseWriter, r *http.Request) {
	// dfcandle.goで定義したGetAllCandle関数に渡す引数を定義
	limit := 100
	duration := "1m" // 1s or 1m or 1h
	durationTime := config.Config.Durations[duration]

	// GetAllCandle関数に上記で定義した引数を渡して得られたデータをdfに格納
	df, _ := models.GetAllCandle(config.Config.ProductCode, durationTime, limit)

	err := templates.ExecuteTemplate(w, "google.html", df.Candles)
	// エラーの場合はInternalServerErrorを表示
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func StartWebServer() error {
	// /chart/ にアクセスされた時にviewChartHandlerを呼び出す
	http.HandleFunc("/chart/", viewChartHandler)

	// ListenAndServeでConfigで定義したPortに接続させる
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}
