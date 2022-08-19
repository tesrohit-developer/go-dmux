package plugins

import (
	"fmt"
	gplugin "github.com/hashicorp/go-plugin"
	"net/rpc"
)

// CheckMessageSidelineImpl is the interface that we're exposing as a plugin.
type CheckMessageSidelineImpl interface {
	CheckMessageSideline(key []byte) ([]byte, error)
	SidelineMessage(msg []byte) SidelineMessageResponse
}

// Here is an implementation that talks over RPC
type CheckMessageSidelineRPC struct {
	Client *rpc.Client
}

func (g *CheckMessageSidelineRPC) CheckMessageSideline(key []byte) ([]byte, error) {
	var resp []byte
	fmt.Println("Checking from dmux plugin")
	err := g.Client.Call("Plugin.CheckMessageSideline", key, &resp)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return resp, nil
}

func (g *CheckMessageSidelineRPC) SidelineMessage(msg []byte) SidelineMessageResponse {
	var resp SidelineMessageResponse
	err := g.Client.Call("Plugin.SidelineMessage", msg, &resp)
	if err != nil {
		fmt.Println(err.Error())
		sidelineMessageResponse := SidelineMessageResponse{UnknownError: true}
		return sidelineMessageResponse
	}
	return resp
}

// Here is the RPC server that CheckMessageSidelineRPC talks to, conforming to
// the requirements of net/rpc
type CheckMessageSidelineRPCServer struct {
	// This is the real implementation
	Impl CheckMessageSidelineImpl
}

func (s *CheckMessageSidelineRPCServer) CheckMessageSideline(key []byte, resp *[]byte) error {
	var err error
	*resp, err = s.Impl.CheckMessageSideline(key)
	return err
}

func (s *CheckMessageSidelineRPCServer) SidelineMessage(msg []byte, sidelineMessageResponse *SidelineMessageResponse) SidelineMessageResponse {
	*sidelineMessageResponse = s.Impl.SidelineMessage(msg)
	return *sidelineMessageResponse
}

// Dummy implementation of a plugin.Plugin interface for use in PluginMap.
// At runtime, a real implementation from a plugin implementation overwrides
// this.
type CheckMessageSidelineImplPlugin struct{}

func (CheckMessageSidelineImplPlugin) Server(*gplugin.MuxBroker) (interface{}, error) {
	return &CheckMessageSidelineRPCServer{}, nil
	//return interface{}, nil
}

func (CheckMessageSidelineImplPlugin) Client(b *gplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &CheckMessageSidelineRPC{Client: c}, nil
	//return interface{}, nil
}
