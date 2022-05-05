module github.com/tesrohit-developer/go-dmux

go 1.12

replace (
	github.com/go-dmux => ./
	github.com/tesrohit-developer/go-dmux/plugins => ./plugins
	honnef.co/go/tools => github.com/dominikh/go-tools v0.0.0-20190102054323-c2f93a96b099
)

require (
	github.com/Shopify/sarama v1.20.1
	github.com/go-dmux v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.7.3
	github.com/samuel/go-zookeeper v0.0.0-20180130194729-c4fab1ac1bec
	github.com/stretchr/testify v1.4.0
	github.com/tesrohit-developer/go-dmux/plugins v0.0.0-00010101000000-000000000000
	golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3 // indirect
	golang.org/x/tools v0.0.0-20190524140312-2c0ae7006135 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0-20170531160350-a96e63847dc3
	honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc // indirect
)
