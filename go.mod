module github.com/tesrohit-developer/go-dmux

go 1.17

replace (
	github.com/go-dmux => ./
	github.com/tesrohit-developer/go-dmux/plugins => ./plugins
)

require (
	github.com/Shopify/sarama v1.20.1
	github.com/go-dmux v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.7.3
	github.com/samuel/go-zookeeper v0.0.0-20180130194729-c4fab1ac1bec
	github.com/stretchr/testify v1.4.0
	github.com/tesrohit-developer/go-dmux/plugins v0.0.0-00010101000000-000000000000
	gopkg.in/natefinch/lumberjack.v2 v2.0.0-20170531160350-a96e63847dc3
)
