package controllers

import (
	"fmt"
	"gotrading/config"
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("./app/views/google.html"))

// Viewを表示する関数を定義
func viewChartHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "google.html", nil)
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
