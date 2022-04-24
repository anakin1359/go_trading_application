## tools
---
| ツール名 | URL | 詳細 |
|:---|:---|:---|
|bitFlyer | https://developers.google.com/chart |bitFlyer | bitflyer top-page |
|bitFlyer Lightning | https://lightning.bitflyer.com/trade | bitflyer trader chart |
|bitFlyer API Document |https://lightning.bitflyer.com/docs?lang=ja | bitFlyerのAPIドキュメント |
|Google Charts | https://developers.google.com/chart | グラフ描画 ([Candlestick Charts](https://developers.google.com/chart/interactive/docs/gallery/candlestickchart)を使用) |

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
## browser access
---
```
http://localhost:8080/chart/
```

## sqlite exec
---
```
$ go run main.go
&{0 {stockdata.sql 0xc0000c40c0} 0 {0 0} [0xc000116000] map[] 0 1 0xc00008e0c0 false map[0xc000116000:map[0xc000116000:true]] map[] 0 0 0 0 <nil> 0 0 0 0 0x4deae0}
2022/04/24 21:04:44 bitflyer.go:244: connecting to wss://ws.lightstream.bitflyer.com/json-rpc
2022/04/24 21:04:52 streamdata.go:16: action=StreamIngestionData, {BTC_JPY 2022-04-24T12:04:52.7816596Z 8151627 5.089392e+06 5.09192e+06 0.1 0.02 454.42344649 736.80110316 5.088659e+06 533.16792366 533.16792366}
2022/04/24 21:04:52 streamdata.go:16: action=StreamIngestionData, {BTC_JPY 2022-04-24T12:04:53.2556747Z 8151636 5.089758e+06 5.092459e+06 0.006 0.02 452.51944649 736.059523 5.088659e+06 533.16792366 533.16792366}
2022/04/24 21:04:53 streamdata.go:16: action=StreamIngestionData, {BTC_JPY 2022-04-24T12:04:53.8921046Z 8151649 5.089955e+06 5.092459e+06 0.011 0.02 452.53044649 732.3163666 5.088659e+06 533.16792366 533.16792366}
```
```
$ ls -lha ./stockdata.sql
-rw-r--r-- 1 username 197609 36K  4月 24 02:02 ./stockdata.sql
```
```
$ which sqlite3
/c/Users/${username}/anaconda3/Library/bin/sqlite3
or
/usr/bin/sqlite3
or
/usr/local/bin/sqlite3
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