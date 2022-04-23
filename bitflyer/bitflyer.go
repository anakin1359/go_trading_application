package bitflyer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
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

// /v1/me/getbalance パラメータ定義
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

// GET /v1/ticker -API Response Sample-
// {
// 	"product_code": "BTC_JPY",
// 	"state": "RUNNING",
// 	"timestamp": "2015-07-08T02:50:59.97",
// 	"tick_id": 3579,
// 	"best_bid": 30000,
// 	"best_ask": 36640,
// 	"best_bid_size": 0.1,
// 	"best_ask_size": 5,
// 	"total_bid_depth": 15.13,
// 	"total_ask_depth": 20,
// 	"market_bid_size": 0,
// 	"market_ask_size": 0,
// 	"ltp": 31690,
// 	"volume": 16819.26,
// 	"volume_by_product": 6819.26
//   }

// /v1/ticker パラメータ定義
type Ticker struct {
	ProductCode     string  `json:"product_code"`
	State           string  `json:"state"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	MarketBidSize   float64 `json:"market_bid_size"`
	MarketAskSize   float64 `json:"market_ask_size"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

// structで定義した内容を利用して売りと買いの中間値を出力するメソッドを定義する
func (t *Ticker) GetMidPrice() float64 {
	return (t.BestBid + t.BestAsk) / 2
}

// 時刻形式をDB対応の形式(RFC3339)に変換させるメソッド
func (t *Ticker) DateTime() time.Time {
	dateTime, err := time.Parse(time.RFC3339, t.Timestamp)
	if err != nil {
		fmt.Print("ここでエラーになっているが一旦スキップする")
		log.Printf("action=DateTime, err=%s", err.Error())
	}
	return dateTime
}

// 時刻単位を調節するメソッド（ex: hh:mm:ss => hh:mm）
func (t *Ticker) TruncateDateTime(duration time.Duration) time.Time {
	return t.DateTime().Truncate(duration)
}

// /v1/ticker にリクエストする処理を定義
func (api *APIClient) GetTicker(productCode string) (*Ticker, error) {
	url := "ticker"
	resp, err := api.doRequest("GET", url, map[string]string{"product_code": productCode}, nil)
	if err != nil {
		return nil, err
	}

	var ticker Ticker
	err = json.Unmarshal(resp, &ticker)
	if err != nil {
		return nil, err
	}

	return &ticker, nil
}

type JsonRPC2 struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Result  interface{} `json:"result,omitempty"`
	Id      *int        `json:"id,omitempty"`
}

type SubscribeParams struct {
	Channel string `json:"channel"`
}

// リアルタイム通信を行うAPIを定義
func (api *APIClient) GetRealTimeTicker(symbol string, ch chan<- Ticker) {
	u := url.URL{Scheme: "wss", Host: "ws.lightstream.bitflyer.com", Path: "/json-rpc"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	channel := fmt.Sprintf("lightning_ticker_%s", symbol)
	if err := c.WriteJSON(&JsonRPC2{Version: "2.0", Method: "subscribe", Params: &SubscribeParams{channel}}); err != nil {
		log.Fatal("subscribe:", err)
		return
	}

OUTER:
	for {
		message := new(JsonRPC2)
		if err := c.ReadJSON(message); err != nil {
			log.Println("read:", err)
			return
		}

		if message.Method == "channelMessage" {
			switch v := message.Params.(type) {
			case map[string]interface{}:
				for key, binary := range v {
					if key == "message" {
						marshaTic, err := json.Marshal(binary)
						if err != nil {
							continue OUTER
						}
						var ticker Ticker
						if err := json.Unmarshal(marshaTic, &ticker); err != nil {
							continue OUTER
						}
						ch <- ticker
					}
				}
			}
		}
	}
}

type Order struct {
	ID                     int     `json:"id"`
	ChildOrderAcceptanceID string  `json:"child_order_acceptance_id"`
	ProductCode            string  `json:"product_code"`
	ChildOrderType         string  `json:"child_order_type"`
	Side                   string  `json:"side"`
	Price                  float64 `json:"price"`
	Size                   float64 `json:"size"`
	MinuteToExpires        int     `json:"minute_to_expire"`
	TimeInForce            string  `json:"time_in_force"`
	Status                 string  `json:"status"`
	ErrorMessage           string  `json:"error_message"`
	AveragePrice           float64 `json:"average_price"`
	ChildOrderState        string  `json:"child_order_state"`
	ExpireDate             string  `json:"expire_date"`
	ChildOrderDate         string  `json:"child_order_date"`
	OutstandingSize        float64 `json:"outstanding_size"`
	CancelSize             float64 `json:"cancel_size"`
	ExecutedSize           float64 `json:"executed_size"`
	TotalCommission        float64 `json:"total_commission"`
	Count                  int     `json:"count"`
	Before                 int     `json:"before"`
	After                  int     `json:"after"`
}

// オーダーした時に返されるデータのレスポンスの型を定義
type ResponseSendChildOrder struct {
	ChildOrderAcceptanceID string `json:"child_order_acceptance_id"`
}

// リクエスト(注文)の処理を定義
func (api *APIClient) SendOrder(order *Order) (*ResponseSendChildOrder, error) {
	data, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	url := "me/sendchildorder"
	resp, err := api.doRequest("POST", url, map[string]string{}, data)
	if err != nil {
		return nil, err
	}

	var response ResponseSendChildOrder
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// リクエスト(参照)の処理を定義
func (api *APIClient) ListOrder(query map[string]string) ([]Order, error) {
	resp, err := api.doRequest("GET", "me/getchildorders", query, nil)
	if err != nil {
		return nil, err
	}
	var responseListOrder []Order
	err = json.Unmarshal(resp, &responseListOrder)
	if err != nil {
		return nil, err
	}
	return responseListOrder, nil
}
