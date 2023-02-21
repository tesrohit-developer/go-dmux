package connection

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/flipkart-incubator/go-dmux/core"
	sink "github.com/flipkart-incubator/go-dmux/http"
	source "github.com/flipkart-incubator/go-dmux/kafka"
	"github.com/flipkart-incubator/go-dmux/offset_monitor"
	sideline_models "github.com/flipkart-incubator/go-dmux/sideline-models"
	"log"
	"os"
)

// **************** CONFIG ***********

//KafkaHTTPWithSidelineConn struct to abstract this connections Run
type KafkaHTTPWithSidelineConn struct {
	EnableDebugLog bool
	Conf           interface{}
	SidelineImpl   interface{}
}

func (c *KafkaHTTPWithSidelineConn) getConfiguration() *KafkaHTTPConnConfig {
	data, _ := json.Marshal(c.Conf)
	var config *KafkaHTTPConnConfig
	json.Unmarshal(data, &config)
	return config
}

//Run method to start this Connection from source to sink
func (c *KafkaHTTPWithSidelineConn) Run() {
	conf := c.getConfiguration()
	log.Println("starting go-dmux with conf", conf)
	if c.EnableDebugLog {
		// enable sarama logs if booted with debug logs
		log.Println("enabling sarama logs")
		sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags)
	}
	kafkaMsgFactory := getKafkaHTTPFactory()
	offMonitor := offset_monitor.GetOffMonitor(conf.OffsetMonitor)
	src := source.GetKafkaSource(conf.Source, kafkaMsgFactory, offMonitor)
	offsetTracker := source.GetKafkaOffsetTracker(conf.PendingAcks, src)
	hook := GetKafkaHook(offsetTracker, c.EnableDebugLog)
	sk := sink.GetHTTPSink(conf.Dmux.Size, conf.Sink)
	sk.RegisterHook(hook)
	src.RegisterHook(hook)

	//hash distribution
	h := GetKafkaMsgHasher()

	d := core.GetDistribution(conf.Dmux.DistributorType, h)

	dmux := core.GetDmux(conf.Dmux, d)
	dmux.ConnectWithSideline(src, sk, c.SidelineImpl.(sideline_models.CheckMessageSideline))
	dmux.Join()
}
