package plugins

type KafkaSidelineMessage struct {
	GroupId           string
	Partition         int32
	EntityId          string
	Offset            int64
	ConsumerGroupName string
	ClusterName       string
	Message           []byte
	Version           int32
}
