package configs

type ConnectionType string

const (
	//KafkaHTTP key to define kafka to generic http sink
	KafkaHTTP ConnectionType = "kafka_http"
	//KafkaFoxtrot key to define kafka to foxtrot http sink
	KafkaFoxtrot ConnectionType = "kafka_foxtrot"
)

/*func getSidelinePlugin() interface{} {
	sidelineImpls := plugins.NewManager("sideline_plugin",
		"sideline-*", "", &plugins.CheckMessageSidelineImplPlugin{})
	// defer sidelineImpls.Dispose()
	// Initialize sidelineImpls manager
	err := sidelineImpls.Init()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Launch all greeters binaries
	sidelineImpls.Launch()
	p, err := sidelineImpls.GetInterface("em")
	if err != nil {
		log.Fatal(err.Error())
	}
	return p
}*/

/*//Start invokes Run of the respective connection in a go routine
func (c ConnectionType) Start(conf interface{}, enableDebug bool) {
	switch c {
	case KafkaHTTP:
		connObj := &connection.KafkaHTTPConn{
			EnableDebugLog: enableDebug,
			Conf:           conf,
			SidelinePlugin: getSidelinePlugin(),
		}
		log.Println("Starting ", KafkaHTTP)
		connObj.Run()
	case KafkaFoxtrot:
		connObj := &connection.KafkaFoxtrotConn{
			EnableDebugLog: enableDebug,
			Conf:           conf,
		}
		log.Println("Starting ", KafkaFoxtrot)
		connObj.Run()
	default:
		panic("Invalid Connection Type")

	}
}*/
