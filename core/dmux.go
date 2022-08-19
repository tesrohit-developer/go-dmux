package core

import (
	"encoding/json"
	"errors"
	"fmt"
	source "github.com/flipkart-incubator/go-dmux/kafka"
	"github.com/tesrohit-developer/go-dmux/plugins"
	"log"
	"math"
	"sync"
	"time"
)

// ControlSignal type defined to build constants used in passing ControlSignals
// to DefaultDmux
type ControlSignal uint8

const (
	//Resize is signal used to trigger resize of a Running Channel
	Resize ControlSignal = 1
	//Stop is signenal used to Stop Dmux
	Stop ControlSignal = 2

	//Sucess response code for ControlSignal Action
	Sucess uint8 = 1
	//Failed response code for ControlSignal Action
	Failed uint8 = 2
)

type Sideline struct {
	Retries           int         `json:"retries"`
	ConsumerGroupName string      `json:"consumerGroupName"`
	ClusterName       string      `json:"clusterName"`
	ConnectionType    string      `json:"type"`
	SidelineMeta      interface{} `json:"sidelineMeta"`
}

//DmuxConf holds configuration parameters for Dmux
type DmuxConf struct {
	Size            int             `json:"size"`
	SourceQSize     int             `json:"source_queue_size"`
	SinkQSize       int             `json:"sink_queue_size"`
	DistributorType DistributorType `json:"distributor_type"`
	BatchSize       int             `json:"batch_size"`
	Version         int             `json:"version"`
	Sideline        Sideline        `json:"sideline"`
}

// ControlMsg is the struct passed to Dmux control Channel to enable it
// perform Admin operations such as Resize
type ControlMsg struct {
	signal ControlSignal
	meta   interface{}
}

//ResizeMeta is the struct used to define resize value which is used when
// Dmux is resizing
type ResizeMeta struct {
	newSize int
}

//ResponseMsg is used for running Dmux instnace to response to client
type ResponseMsg struct {
	signal ControlSignal
	status uint8
}

//Sink is interface that implements OutputSink of Dmux operation
type Sink interface {
	// Clone method is expected to return instance of Sink. If Sink is Stateless
	// this can return selfRefrence back. If Sink is Stateful, its good idea to
	// create new instnace of Sink.
	Clone() Sink

	// Consume method gets The interface.
	//TODO currently this method does not return error, need to solve for error
	// handling
	Consume(msg interface{}, retries int) error

	//BatchConsume method is invoked in batch_size is configured
	BatchConsume(msg []interface{}, version int)
}

//Source is interface that implements input Source to the Dmux
type Source interface {
	//Generate method takes output channel to which it writes data. The
	//implementation can write to to this using multiple goroutines
	//This method is not expected to return, its run in a separate goroutine
	Generate(out chan<- interface{})
	//Method used to trigger GracefulStop of Source
	Stop()
}

//Distributor interface abstracts the Logic to distribute the load from Source
// to Sink. Client can choose to use HashDistributor or RoundRobinDistributor or
//write their own distribution Logic
type Distributor interface {
	//Distribute method take incoming data interface and number of outbound channels
	//to return the index of channel to be selected for Distribution of this message
	Distribute(data interface{}, size int) int
}

//Dmux struct which enables Size based Dmultiplexing for
//Source to Sink connections.
//TODO restrict size to be powers of 2 for better optimization in modulo
type Dmux struct {
	size                   int
	batchSize              int
	sourceQSize, sinkQSize int
	control                chan ControlMsg
	response               chan ResponseMsg
	err                    chan error
	distribute             Distributor
	version                int
	sideline               Sideline
}

const defaultSourceQSize int = 1
const defaultSinkQSize int = 100
const defaultBatchSize int = 1
const defaultVersion int = 1

//GetDmux is public method used to Get instance of a Dmux struct
func GetDmux(conf DmuxConf, d Distributor) *Dmux {
	control := make(chan ControlMsg)
	response := make(chan ResponseMsg)
	err := make(chan error)
	sourceQSize := defaultSourceQSize
	sinkQSize := defaultSinkQSize
	batchSize := defaultBatchSize
	version := defaultVersion

	if conf.SourceQSize > 0 {
		sourceQSize = conf.SourceQSize
	}

	if conf.SinkQSize > 0 {
		sinkQSize = conf.SinkQSize
	}

	if conf.BatchSize > 0 {
		batchSize = conf.BatchSize
	}

	if conf.Version > 0 {
		version = conf.Version
	}

	output := &Dmux{conf.Size, batchSize, sourceQSize, sinkQSize,
		control, response, err, d, version, conf.Sideline}
	return output
}

//Connect method holds Dmux logic used to Connect Source to Sink
func (d *Dmux) Connect(source Source, sink Sink, sidelinePlugin interface{}) {
	go d.run(source, sink, sidelinePlugin)
}

//Await method added to enable testing when using bounded source
func (d *Dmux) Await(duration time.Duration) {

	select {
	case e := <-d.err:
		if e != nil {
			panic(e)
		}
	case <-time.After(duration):
		fmt.Println("Timedout")
	}

}

//Join used to sleep the main routine forever
func (d *Dmux) Join() {
	e := <-d.err
	if e != nil {
		panic(e)
	}
}

//Resize method is used to Resize a running Dmux
func (d *Dmux) Resize(size int) {
	d.control <- getResizeMsg(size)
	<-d.response
}

// Stop is used to GracefulStop running Dmux
func (d *Dmux) Stop() {
	d.control <- getStopMsg()
	<-d.response
}

func getResizeMsg(size int) ControlMsg {
	c := ControlMsg{
		signal: Resize,
		meta:   ResizeMeta{size},
	}
	return c
}

func getStopMsg() ControlMsg {
	c := ControlMsg{
		signal: Stop,
	}
	return c
}

func (d *Dmux) run(source Source, sink Sink, sidelinePlugin interface{}) {

	ch, wg := setup(d.size, d.sinkQSize, d.batchSize, sink, d.version, d.sideline, sidelinePlugin)
	in := make(chan interface{}, d.sourceQSize)
	//start source
	//TODO handle panic recovery if in channel is closed for shutdown
	go source.Generate(in)

	for {
		select {
		case data := <-in:
			i := d.distribute.Distribute(data, len(ch))
			// fmt.Printf("writing to channel %d len %d", i, len(ch[i]))
			ch[i] <- data
		case ctrl := <-d.control:
			if ctrl.signal == Resize {
				fmt.Println("processing resize")
				shutdown(ch, wg)
				resizeMeta := ctrl.meta.(ResizeMeta)
				ch, wg = setup(resizeMeta.newSize, d.sinkQSize, d.batchSize, sink, d.version, d.sideline, sidelinePlugin)
				d.response <- ResponseMsg{ctrl.signal, Sucess}
			} else if ctrl.signal == Stop {
				fmt.Println("processing stop")
				source.Stop()
				shutdown(ch, wg)
				close(in)
				d.response <- ResponseMsg{ctrl.signal, Sucess}
				d.err <- nil
				return
			}
		}
	}

}

func shutdown(ch []chan interface{}, wg *sync.WaitGroup) {
	for _, c := range ch {
		close(c)
	}
	wg.Wait()
}

func setup(size, qsize, batchSize int, sink Sink, version int, sideline Sideline, sidelinePlugin interface{}) ([]chan interface{}, *sync.WaitGroup) {
	if version == 1 && batchSize == 1 {
		if sidelinePlugin != nil {
			log.Printf("Calling simpleSetupWithSideline")
			return simpleSetupWithSideline(size, qsize, sink, sideline, sidelinePlugin)
		} else {
			log.Printf("Calling simpleSetup")
			return simpleSetup(size, qsize, sink)
		}
	} else {
		log.Printf("Calling batchSetup")
		return batchSetup(size, qsize, batchSize, sink, version)
	}
}

// Batching has been implemented a bit differently to provide more value to end Client.
// Typical batching would involve queing up till batch size and flush.
// This is inefficient at Client consumption when using modulo hashing (HashDistributor).
// Client of Batch Sink has to run single threaded code to ensure ordering is not lost at consumption.
// The implementation below avoids this shortcoming.
// In this implementation we create total = size * batchSize number of channels.
// And we create numer of BatchConsumer routine = size.
// Each BatchConsumer routine will consume (index, index+batchSize); index = index of consumer
// which is 0 to size-1.
// BatchConsumer will update its batch array index from one entry each of respective channel index. (This provides
// ability for consumer to consume in parallel) and then flush the batch.
// Close of any channel in a BatchConsumer will stop the BatchConsumer.
func batchSetup(sz, qsz, batchsz int, sink Sink, version int) ([]chan interface{}, *sync.WaitGroup) {
	size := sz * batchsz // create double nuber of channels

	wg := new(sync.WaitGroup)
	wg.Add(size)
	ch := make([]chan interface{}, size)

	//create channels
	for i := 0; i < size; i++ {
		ch[i] = make(chan interface{}, qsz)
	}

	//assign batch of channels per batch consumer
	for i := 0; i < size; i += batchsz {
		//async runner does batching and call to sink.BatchConsume
		go func(index int) {
			sk := sink.Clone()
			batch := make([]interface{}, batchsz) //reuse one time created array

			for { // run till any channel closes check failed case which breaks from this loop
				// fmt.Println("Running consumer ", index)
				//batch
				failed := false
				j := index
				for z := 0; z < batchsz; z++ {
					// fmt.Printf("Iterating count %d consumer %d j %d channelLen %d \n", z, index, j, len(ch[j]))
					msg, more := <-ch[j]  // more == false if channel is closed
					if !failed && !more { // failed = true if any one channel is closed
						failed = true
						break
					}
					batch[z] = msg
					j++
				}
				// fmt.Println("Failed ", failed)
				if failed { //ack all waitGroups and return to break from the infinite for loop
					for z := 0; z < batchsz; z++ {
						wg.Done()
					}
					return
				}
				// fmt.Println("flusing ", batch)
				//flush batched message
				sk.BatchConsume(batch, version)
			}

		}(i)
	}
	return ch, wg
}

func simpleSetupWithSideline(size, qsize int, sink Sink, sideline Sideline, sidelinePlugin interface{}) ([]chan interface{}, *sync.WaitGroup) {
	wg := new(sync.WaitGroup)
	wg.Add(size)
	ch := make([]chan interface{}, size)
	for i := 0; i < size; i++ {
		ch[i] = make(chan interface{}, qsize)
		go func(index int) {
			sk := sink.Clone()
			for msg := range ch[index] {
				val := msg.(source.KafkaMsg)
				var checkBytes []byte
				var check plugins.CheckMessageSidelineResponse
				var checkErr = errors.New("")
				for {
					var retryMessage = false
					log.Printf("Checking if the message is already sidelined %d, %d", val.GetRawMsg().Partition, val.GetRawMsg().Offset)
					checkBytes, checkErr = sidelinePlugin.(plugins.CheckMessageSidelineImpl).
						CheckMessageSideline([]byte(string(val.GetRawMsg().Key) + sideline.ConsumerGroupName + sideline.ClusterName))
					err := json.Unmarshal(checkBytes, &check)
					if err != nil {
						errors.New("error in serde of CheckMessageSidelineResponse")
						continue
					}
					if checkErr != nil {
						continue
					}

					log.Printf("Message if already sidelined %t %d %d", check.IsMessageSidelined, val.GetRawMsg().Partition, val.GetRawMsg().Offset)
					if check.IsMessageSidelined {
						var version int32
						version = check.Version
						sidelineMetaByteArray, sidelineMetaByteArrayErr := json.Marshal(sideline.SidelineMeta)
						if sidelineMetaByteArrayErr != nil {
							log.Printf("error in serde of SidelineMeta")
							continue
						}
						kafkaSidelineMessage := plugins.SidelineMessage{
							GroupId:           string(val.GetRawMsg().Key),
							Partition:         val.GetRawMsg().Partition,
							EntityId:          string(val.GetRawMsg().Key) + sideline.ConsumerGroupName + sideline.ClusterName,
							Offset:            val.GetRawMsg().Offset,
							ConsumerGroupName: sideline.ConsumerGroupName,
							ClusterName:       sideline.ClusterName,
							Message:           val.GetRawMsg().Value,
							Version:           version,
							ConnectionType:    sideline.ConnectionType,
							SidelineMeta:      sidelineMetaByteArray,
						}
						for {
							log.Printf("Sidelining the message %d, %d", val.GetRawMsg().Partition, val.GetRawMsg().Offset)
							sidelineByteArray, err := json.Marshal(kafkaSidelineMessage)
							if err != nil {
								log.Printf("error in serde of kafkaSidelineMessage")
								continue
							}
							sidelineMessageResponse := sidelinePlugin.(plugins.CheckMessageSidelineImpl).SidelineMessage(sidelineByteArray)
							if !sidelineMessageResponse.Success {
								log.Printf(sidelineMessageResponse.ErrorMessage)
								if sidelineMessageResponse.ConcurrentModificationError {
									retryMessage = true
									break
								}
								continue
							}
							break
						}
						if retryMessage {
							continue
						}
						break
					} else {
						consumeError := sk.Consume(msg, sideline.Retries)
						if consumeError != nil && consumeError.Error() == "exceeded retries" {
							sidelineMetaByteArray, sidelineMetaByteArrayErr := json.Marshal(sideline.SidelineMeta)
							if sidelineMetaByteArrayErr != nil {
								log.Printf("error in serde of SidelineMeta")
								continue
							}
							kafkaSidelineMessage := plugins.SidelineMessage{
								GroupId:           string(val.GetRawMsg().Key),
								Partition:         val.GetRawMsg().Partition,
								EntityId:          string(val.GetRawMsg().Key) + sideline.ConsumerGroupName + sideline.ClusterName,
								Offset:            val.GetRawMsg().Offset,
								ConsumerGroupName: sideline.ConsumerGroupName,
								ClusterName:       sideline.ClusterName,
								Message:           val.GetRawMsg().Value,
								Version:           0,
								ConnectionType:    sideline.ConnectionType,
								SidelineMeta:      sidelineMetaByteArray,
							}
							for {
								log.Printf("Sidelining the message as exceeded retries %d, %d", val.GetRawMsg().Partition, val.GetRawMsg().Offset)
								sidelineByteArray, err := json.Marshal(kafkaSidelineMessage)
								if err != nil {
									errors.New("error in serde of kafkaSidelineMessage")
									continue
								}
								sidelineMessageResponse := sidelinePlugin.(plugins.CheckMessageSidelineImpl).SidelineMessage(sidelineByteArray)
								if !sidelineMessageResponse.Success {
									log.Printf(sidelineMessageResponse.ErrorMessage)
									if sidelineMessageResponse.ConcurrentModificationError {
										retryMessage = true
										break
									}
									continue
								}
								break
							}
						}
						if retryMessage {
							continue
						}
						break
					}
				}
			}
			wg.Done()
		}(i)
	}
	return ch, wg
}

func simpleSetup(size, qsize int, sink Sink) ([]chan interface{}, *sync.WaitGroup) {
	wg := new(sync.WaitGroup)
	wg.Add(size)
	ch := make([]chan interface{}, size)
	for i := 0; i < size; i++ {
		ch[i] = make(chan interface{}, qsize)
		go func(index int) {
			sk := sink.Clone()
			for msg := range ch[index] {
				sk.Consume(msg, math.MaxInt32)
			}
			wg.Done()
		}(i)
	}
	return ch, wg
}
