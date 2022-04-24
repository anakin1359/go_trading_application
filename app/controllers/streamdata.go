package controllers

import (
	"gotrading/app/models"
	"gotrading/bitflyer"
	"gotrading/config"
	"log"
)

// bitFlyerから取得したデータをストリーミングする関数を定義
func StreamIngestionData() {
	var tickerChannel = make(chan bitflyer.Ticker)
	apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)
	go apiClient.GetRealTimeTicker(config.Config.ProductCode, tickerChannel)
	for ticker := range tickerChannel {
		log.Printf("action=StreamIngestionData, %v", ticker)
		for _, duration := range config.Config.Durations {
			isCreated := models.CreateCandleWithDuration(ticker, ticker.ProductCode, duration)
			if isCreated == true && duration == config.Config.TradeDuration {
				// TODO
			}
		}
	}
}
