package main

import (
	"github.com/tesrohit-developer/go-dmux/metrics"
	"log"
	"os"

	"github.com/tesrohit-developer/go-dmux/logging"
)

//

// **************** Bootstrap ***********

func main() {
	args := os.Args[1:]
	sz := len(args)

	var path string

	if sz == 1 {
		path = args[0]
	}

	dconf := DMuxConfigSetting{
		FilePath: path,
	}
	conf := dconf.GetDmuxConf()

	dmuxLogging := new(logging.DMuxLogging)
	dmuxLogging.Start(conf.Logging)

	c := Controller{config: conf}
	go c.start()

	log.Printf("config: %v", conf)

	//start showing metrics at the endpoint
	metrics.Start(conf.MetricPort)

	for _, item := range conf.DMuxItems {
		go func(connType ConnectionType, connConf interface{}, logDebug bool, sidelineEnabled bool) {
			connType.Start(connConf, logDebug, sidelineEnabled)
		}(item.ConnType, item.Connection, dmuxLogging.EnableDebug, item.SidelineEnabled)
	}

	//main thread halts. TODO make changes to listen to kill and reboot
	select {}
}
