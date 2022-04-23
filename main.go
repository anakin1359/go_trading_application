package main

import (
	"fmt"
	"gotrading/bitflyer"
	"gotrading/config"
	"gotrading/utils"
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
	fmt.Println(apiClient.GetBalance())
}
