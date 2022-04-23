package bitflyer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// https://lightning.bitflyer.com/docs?lang=ja
const baseURL = "https://api.bitflyer.com/v1/"

type APIClient struct {
	key        string
	secret     string
	httpClient *http.Client
}

// configで定義したapi-key, api-secret を参照するコンストラクタ
func New(key, secret string) *APIClient {
	apiClient := &APIClient{key, secret, &http.Client{}}
	return apiClient
}

// API Requestを実行する時のHeaderを定義
func (api APIClient) header(method, endpoint string, body []byte) map[string]string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	log.Println(timestamp)

	message := timestamp + method + endpoint + string(body)

	// Header情報は bitFlyer Lightning Document を参照
	mac := hmac.New(sha256.New, []byte(api.secret))
	mac.Write([]byte(message))
	sign := hex.EncodeToString(mac.Sum(nil))
	return map[string]string{
		"ACCESS-KEY":       api.key,
		"ACCESS-TIMESTAMP": timestamp,
		"ACCESS-SIGN":      sign,
		"Content-Type":     "application/json",
	}
}

// HTTP Request の定義
func (api *APIClient) doRequest(method, urlPath string, query map[string]string, data []byte) (body []byte, err error) {

	// エンドポイントの正当性を判定 https://api.bitflyer.com/v1/
	baseURL, err := url.Parse(baseURL)
	if err != nil {
		return
	}

	// 引数として渡されたエンドポイントの後に続くメソッドの正当性を判定
	apiURL, err := url.Parse(urlPath)
	if err != nil {
		return
	}

	// エンドポイントとメソッドを繋ぎ合わせて実行したいAPIのエンドポイントを定義
	endpoint := baseURL.ResolveReference(apiURL).String()

	// ログ出力処理を定義
	log.Printf("action=doRequest endpoint=%s", endpoint)

	// http requestの実行処理を定義
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(data))
	if err != nil {
		return
	}

	// http request実行後にクエリが渡された場合は取得する必要がある（Key, Value形式で取得）
	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}

	// RawQueryを実行する場合はエンコードする処理が必要
	req.URL.RawQuery = q.Encode()

	// header情報が存在する場合はその情報をheader情報に追加する
	for key, value := range api.header(method, req.URL.RequestURI(), data) {
		req.Header.Add(key, value)
	}

	// http responseの定義(エラーの場合はレスポンスがないためnilとエラーを返す)
	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// 処理終了後に必ず実行
	defer resp.Body.Close()

	// ioutilで返却された値のBodyの情報を読み込む
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 何も存在しなければbodyの情報とnilを返却する
	return body, nil
}

// GET /v1/me/getbalance -API Response Sample-
// [
//   {
//     "currency_code": "JPY",
//     "amount": 1024078,
//     "available": 508000
//   },
//   {
//     "currency_code": "BTC",
//     "amount": 10.24,
//     "available": 4.12
//   },
//   {
//     "currency_code": "ETH",
//     "amount": 20.48,
//     "available": 16.38
//   }
// ]
type Balance struct {
	CurrentCode string  `json:"currency_code"`
	Amount      float64 `json:"amount"`    // いくら保有しているか
	Available   float64 `json:"available"` // いくら使用するか
}

//  /v1/me/getbalance にリクエストする処理を定義
func (api *APIClient) GetBalance() ([]Balance, error) {
	url := "me/getbalance"
	resp, err := api.doRequest("GET", url, map[string]string{}, nil)
	log.Printf("url=%s resp=%s", url, string(resp))
	if err != nil {
		log.Printf("action=GetBalance err=%s", err.Error())
		return nil, err
	}

	var balance []Balance
	err = json.Unmarshal(resp, &balance)
	if err != nil {
		log.Printf("action=GetBalance err=%s", err.Error())
		return nil, err
	}

	return balance, nil
}
