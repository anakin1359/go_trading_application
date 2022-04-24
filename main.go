package main

import (
	"gotrading/app/controllers"
	"gotrading/config"
	"gotrading/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	controllers.StreamIngestionData()
	controllers.StartWebServer()
}

// bitFlyerでの自動売買処理は実行時以外は無効化
// func main() {
// 	utils.LoggingSettings(config.Config.LogFile)
// 	apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)

// 	order := &bitflyer.Order{
// 		ProductCode:     config.Config.ProductCode,
// 		ChildOrderType:  "MARKET", // 成行 => 指値の場合はLIMIT
// 		Side:            "BUY",    // 購入 => 売却の場合はSELL
// 		Size:            0.01,     // Bitcoinの数量
// 		MinuteToExpires: 1,        // 分
// 		TimeInForce:     "GTC",    // キャンセルするまで有効な注文
// 	}
// 	res, _ := apiClient.SendOrder(order)
// 	fmt.Println(res.ChildOrderAcceptanceID)

// 	i := "JRF20181012-144016-140584" // i := "JRF2022XXXX-XXXXXX-XXXXXX"
// 	params := map[string]string{
// 		"product_code":              config.Config.ProductCode,
// 		"child_order_acceptance_id": i,
// 	}
// 	r, _ := apiClient.ListOrder(params)
// 	fmt.Println(r)
// }
