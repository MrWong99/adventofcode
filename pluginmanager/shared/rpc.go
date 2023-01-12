package shared

import "net/rpc"

type RPCClient struct {
	client *rpc.Client
}

func (c *RPCClient) Calculate(input string) (resp string, err error) {
	err = c.client.Call("Plugin.Calculate", input, &resp)
	return
}

type RPCServer struct {
	Impl CalcService
}

func (s *RPCServer) Calculate(input string, resp *string) error {
	val, err := s.Impl.Calculate(input)
	if err == nil {
		*resp = val
	}
	return err
}
