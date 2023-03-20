package sideline

import (
	"encoding/json"
	"github.com/flipkart-incubator/go-dmux/sideline_impls"
)

type CheckMessageSidelineImpl struct {
}

func sidelineInitExample() {
	custom := sideline_impls.DmuxCustom{}
	sidelineImpl := &CheckMessageSidelineImpl{}
	path := "" // config path
	custom.DmuxStart(path, sidelineImpl)
}

func (c *CheckMessageSidelineImpl) CheckMessageSideline(key []byte) ([]byte, error) {
	checkMessageSidelineResponse := CheckMessageSidelineResponse{}
	return json.Marshal(checkMessageSidelineResponse)
}

func (c *CheckMessageSidelineImpl) SidelineMessage(msg []byte) SidelineMessageResponse {
	sidelineMessageResponse := SidelineMessageResponse{}
	return sidelineMessageResponse
}

func (c *CheckMessageSidelineImpl) InitialisePlugin(confBytes []byte) error {
	return nil
}
