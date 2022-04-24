## project
---
```
.
|-- README.md
|-- app
|   `-- models
|       |-- base.go
|       `-- candle.go
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

## sqlite exec
---
```
$ go run main.go
&{0 {stockdata.sql 0xc00005a4e0} 0 {0 0} [0xc0000ae000] map[] 0 1 0xc000028120 false map[0xc0000ae000:map[0xc0000ae000:true]] map[] 0 0 0 0 <nil> 0 0 0 0 0x7ff6da9c8840}
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