package main

import (
	"fmt"
	"gotrading/bitflyer"
	"gotrading/config"
	"gotrading/utils"
	"time"
)

// func main() {
// 	str_prm := "\n===== [ TEST ] ======"
// 	log.Println(str_prm)
// 	fmt.Println(config.Config.ApiKey)
// 	fmt.Println(config.Config.ApiSecret)
// }

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)
	ticker, _ := apiClient.GetTicker("BTC_USD")
	fmt.Println(ticker)
	fmt.Println(ticker.GetMidPrice())
	fmt.Println(ticker.DateTime())
	fmt.Println(ticker.TruncateDateTime(time.Hour))
}
