module github.com/tesrohit-developer/go-dmux

go 1.12

replace (
	github.com/go-dmux => ./
	github.com/tesrohit-developer/go-dmux/configs => ./configs
	github.com/tesrohit-developer/go-dmux/plugins => ./plugins
)

require (
	github.com/DataDog/zstd v1.5.2 // indirect
	github.com/Shopify/sarama v1.20.1
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/prometheus/client_golang v1.12.1
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/samuel/go-zookeeper v0.0.0-20180130194729-c4fab1ac1bec
	github.com/stretchr/testify v1.4.0
	golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3 // indirect
	golang.org/x/tools v0.0.0-20190524140312-2c0ae7006135 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0-20170531160350-a96e63847dc3
	honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc // indirect
)
