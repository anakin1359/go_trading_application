サイト## develop environment
---
| 項目 | 種別 |go Lang version |
|:---|:---|:---|
| OS | Windows10 |go version go1.18 windows/amd64|
| OS | Ubuntu 18.04.5 LTS (Bionic Beaver) |go version go1.18.1 linux/amd64|
<br>

## tools
---
| ツール名 | URL | 詳細 |
|:---|:---|:---|
| bitFlyer | https://developers.google.com/chart |bitFlyer | bitflyer top-page |
| bitFlyer Lightning | https://lightning.bitflyer.com/trade | bitflyer trader chart |
| bitFlyer API Document |https://lightning.bitflyer.com/docs?lang=ja | bitFlyerのAPIドキュメント |
| Google Charts (Candlestick) | https://developers.google.com/chart | ロウソク型チャート ([Candlestick Charts](https://developers.google.com/chart/interactive/docs/gallery/candlestickchart)を使用) |
| Google Charts (Combo) | https://developers.google.com/chart | 複数種類グラフ使用 ([Combo Charts](https://developers.google.com/chart/interactive/docs/gallery/combochart)を使用) |
| Google Charts (graph Control) | https://developers.google.com/chart | グラフ操作 ([Controls and Dashboards](https://developers.google.com/chart/interactive/docs/gallery/controls)を使用) |
| jQuery | https://jquery.com/ | jquery チュートリアルサイト |
| jQuery | https://www.w3schools.com/jquery/default.asp | jquery サンプルコード |


<br>

## project
---
```
.
|-- README.md
|-- app
|   |-- controllers
|   |   |-- streamdata.go
|   |   `-- webserver.go
|   |-- models
|   |   |-- base.go
|   |   |-- candle.go
|   |   `-- dfcandle.go
|   `-- views
|       `-- google.html
|-- bitflyer
|   `-- bitflyer.go
|-- config
|   `-- config.go
|-- config.ini
|-- go.mod
|-- go.sum
|-- gotrading.log
|-- main.go
|-- stockdata.sql
`-- utils
    `-- logging.go
```
<br>

## config.ini
---
```
touch ./config.ini; ls -lh $_
```
```
[bitflyer]
api_key = XXXXXXXXXXXXXXXXXXXXXX
api_secret = XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX

[gotrading]
log_file = gotrading.log
product_code = BTC_JPY // BTC_USD
trade_duration = 1m

[db]
name = stockdata.sql
driver = sqlite3

[web]
port = 8080
```
<br>

## browser access (chart)
---
```
http://localhost:8080/chart/
```
<br>

## browser access (ajax)
---
```
http://localhost:8080/api/candle/?product_code={product_code}&duration={1s/m/h}&limit={1-1000}
```
ex)
```
http://localhost:8080/api/candle/?product_code=BTC_JPY&duration=1m&limit=1
```
<br>

## sqlite exec
---
```
$ go run main.go
2022/04/25 01:51:01 bitflyer.go:244: connecting to wss://ws.lightstream.bitflyer.com/json-rpc
2022/04/25 01:51:02 streamdata.go:17: action=StreamIngestionData, {BTC_JPY 2022-04-24T16:50:59.8864462Z 8279633 5.068835e+06 5.070359e+06 0.26 0.021 446.02196243 833.55624702 5.068883e+06 630.37817952 630.37817952}
2022/04/25 01:51:02 streamdata.go:17: action=StreamIngestionData, {BTC_JPY 2022-04-24T16:51:00.4623748Z 8279650 5.068835e+06 5.07035e+06 0.26 0.02361368 449.01806243 834.1174607 5.068883e+06 630.37817952 630.37817952}
2022/04/25 01:51:03 streamdata.go:17: action=StreamIngestionData, {BTC_JPY 2022-04-24T16:51:01.0734443Z 8279655 5.068835e+06 5.07035e+06 0.26 0.02361368 449.01806243 833.6805607 5.068883e+06 630.37817952 630.37817952}
```
```
$ ls -lha ./stockdata.sql
-rw-r--r-- 1 username 197609 36K  4月 24 02:02 ./stockdata.sql
```
```
$ which sqlite3
/c/Users/${username}/anaconda3/Library/bin/sqlite3 or /usr/bin/sqlite3 or /usr/local/bin/sqlite3
```
```
$ sqlite3 stockdata.sql
SQLite version 3.36.0 2021-06-18 18:36:39
Enter ".help" for usage hints.
```
```
sqlite> .tables
BTC_JPY_1h0m0s  BTC_JPY_1m0s    BTC_JPY_1s      signal_events
```
```
sqlite> select * from BTC_JPY_1s;
2022-04-24T11:49:25Z|5093062.5|5093062.5|5093062.5|5093062.5|528.00569753
2022-04-24T11:49:26Z|5093062.5|5093062.5|5093062.5|5093062.5|528.00569753
2022-04-24T11:49:27Z|5093062.5|5093062.5|5093062.5|5093062.5|1056.01139506
```