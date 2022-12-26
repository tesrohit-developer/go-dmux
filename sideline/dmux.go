package sideline

import (
	"fmt"
	c "github.com/flipkart-incubator/go-dmux/config"
	"github.com/flipkart-incubator/go-dmux/logging"
	"github.com/flipkart-incubator/go-dmux/metrics"
	"log"
)

//

type DmuxCustom struct {
}

func (d *DmuxCustom) DmuxStart(path string) {
	//fmt.Println(checkMessageSideline.SidelineMessage())

	dconf := c.DMuxConfigSetting{
		FilePath: path,
	}
	conf := dconf.GetDmuxConf()

	dmuxLogging := new(logging.DMuxLogging)
	dmuxLogging.Start(conf.Logging)

	c := c.Controller{config: conf}
	go c.start()

	log.Printf("config: %v", conf)

	//start showing metrics at the endpoint
	metrics.Start(conf.MetricPort)

	for _, item := range conf.DMuxItems {
		fmt.Println(item.ConnType)
		/*go func(connType ConnectionType, connConf interface{}, logDebug bool) {
			connType.Start(connConf, logDebug)
		}(item.ConnType, item.Connection, dmuxLogging.EnableDebug)*/
	}

	//main thread halts. TODO make changes to listen to kill and reboot
	select {}
}
