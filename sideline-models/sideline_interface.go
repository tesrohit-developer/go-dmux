package sideline_models

type CheckMessageSideline interface {
	CheckMessageSideline(key []byte) ([]byte, error)
	SidelineMessage(msg []byte) SidelineMessageResponse
	InitialisePlugin(conf []byte) error
}
