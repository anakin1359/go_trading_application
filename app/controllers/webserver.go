package controllers

import (
	"encoding/json"
	"fmt"
	"gotrading/app/models"
	"gotrading/config"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
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

// api通信時にエラーが発生した場合にレスポンスとして受け取るJSONを定義
type JSONError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// api通信時にエラーが発生した場合のレスポンス処理を定義
func APIError(w http.ResponseWriter, errMessage string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonError, err := json.Marshal(JSONError{Error: errMessage, Code: code})
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonError)
}

// apiのエンドポイントを判定するための正規表現を定義
var apiValidPath = regexp.MustCompile("^/api/candle/$")

func apiMakeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// エンドポイントがマッチングするかを判定する処理を定義
		m := apiValidPath.FindStringSubmatch(r.URL.Path)

		// 判定した結果、マッチングしたものが見当たらない(0)の場合は上記で定義したAPIError関数を使用してエラーを返す
		if len(m) == 0 {
			APIError(w, "Not found", http.StatusNotFound)
		}

		// 判定した結果、マッチングしたものが見つかった場合はResponseWriter(w)とRequest(r)を関数に返却する
		fn(w, r)
	}
}

// api通信を実行する関数の大元(この関数がレスポンスを返却する)
func apiCandleHandler(w http.ResponseWriter, r *http.Request) {
	// browserからproduct_codeを選択できるようにするための設定(browserからproduct_codeを送信)
	productCode := r.URL.Query().Get("product_code")

	// product_codeが見つからなかった場合のエラーハンドリング
	if productCode == "" {
		APIError(w, "No product_code param", http.StatusBadRequest)
		return
	}

	// browserからlimitを取得できるようにするための設定
	strLimit := r.URL.Query().Get("limit")

	// limitに格納されたデータをASCII => Integerに変換
	limit, err := strconv.Atoi(strLimit)

	// limitに何もデータが入っていない場合、またはerrが空でない場合、またはlimitが0より小さい場合、またはlimitが1000以上の場合
	if strLimit == "" || err != nil || limit < 0 || limit > 1000 {
		// 上記のいずれかに当てはまる場合はlimitの数値を1000に変換
		limit = 1000
	}

	// browserからduration(時刻形式)を取得できるようにするための設定
	duration := r.URL.Query().Get("duration")

	// durationに何もデータが入っていない場合は、デフォルトで1mを設定
	if duration == "" {
		duration = "1m"
	}
	durationTime := config.Config.Durations[duration]

	// 想定外の処理を全て拾った後に、各項目のデフォルト値を「df」に代入
	df, _ := models.GetAllCandle(productCode, durationTime, limit)

	// 「df」を使用して構造体をJSONに変換
	js, err := json.Marshal(df)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// JSON形式への変換が正常に行われた場合はHeaderに"Content-Type", "application/json" を設定
	w.Header().Set("Content-Type", "application/json")

	// ResponseWriterにJSONに変換した情報を返す
	w.Write(js)

}

func StartWebServer() error {
	// /api/candle/ にアクセスされた時にapiMakeHandler関数を実行(引数として上記で定義したapiCandleHandler関数を指定)
	http.HandleFunc("/api/candle/", apiMakeHandler(apiCandleHandler))

	// /chart/ にアクセスされた時にviewChartHandlerを呼び出す
	http.HandleFunc("/chart/", viewChartHandler)

	// ListenAndServeでConfigで定義したPortに接続させる
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}
