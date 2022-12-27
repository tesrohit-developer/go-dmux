module github.com/flipkart-incubator/go-dmux/connection

go 1.12

require github.com/flipkart-incubator/go-dmux/offset_monitor v0.0.0

require github.com/flipkart-incubator/go-dmux/core v0.0.0

require github.com/flipkart-incubator/go-dmux/http v0.0.0

require (
	github.com/Shopify/sarama v1.20.1
	github.com/flipkart-incubator/go-dmux/kafka v0.0.0
)

replace github.com/flipkart-incubator/go-dmux/connection v0.0.0 => ../connection

replace github.com/flipkart-incubator/go-dmux/logging v0.0.0 => ../logging

replace github.com/flipkart-incubator/go-dmux/offset_monitor v0.0.0 => ../offset_monitor

replace github.com/flipkart-incubator/go-dmux/http v0.0.0 => ../http

replace github.com/flipkart-incubator/go-dmux/kafka v0.0.0 => ../kafka

replace github.com/flipkart-incubator/go-dmux/kafka/consumer-group v0.0.0 => ../kafka/consumer-group

replace github.com/flipkart-incubator/go-dmux/sideline v0.0.0 => ../sideline

replace github.com/flipkart-incubator/go-dmux/sideline-models v0.0.0 => ../sideline-models

replace github.com/flipkart-incubator/go-dmux/metrics v0.0.0 => ../metrics

replace github.com/flipkart-incubator/go-dmux/config v0.0.0 => ../config

replace github.com/flipkart-incubator/go-dmux/core v0.0.0 => ../core
