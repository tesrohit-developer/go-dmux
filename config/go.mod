module github.com/flipkart-incubator/go-dmux/config

go 1.12

require github.com/flipkart-incubator/go-dmux/connection v0.0.0

require (
	github.com/DataDog/zstd v1.5.2 // indirect
	github.com/Shopify/toxiproxy v2.1.4+incompatible // indirect
	github.com/eapache/go-resiliency v1.3.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/flipkart-incubator/go-dmux/logging v0.0.0
	github.com/frankban/quicktest v1.14.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
)

replace github.com/flipkart-incubator/go-dmux/connection v0.0.0 => ../connection

replace github.com/flipkart-incubator/go-dmux/logging v0.0.0 => ../logging

replace github.com/flipkart-incubator/go-dmux/offset_monitor v0.0.0 => ../offset_monitor

replace github.com/flipkart-incubator/go-dmux/http v0.0.0 => ../http

replace github.com/flipkart-incubator/go-dmux/kafka v0.0.0 => ../kafka

replace github.com/flipkart-incubator/go-dmux/kafka/consumer-group v0.0.0 => ../kafka/consumer-group

replace github.com/flipkart-incubator/go-dmux/sideline v0.0.0 => ../sideline

replace github.com/flipkart-incubator/go-dmux/metrics v0.0.0 => ../metrics

replace github.com/flipkart-incubator/go-dmux/config v0.0.0 => ../config

replace github.com/flipkart-incubator/go-dmux/core v0.0.0 => ../core
