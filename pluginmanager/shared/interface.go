package shared

import (
	"context"
	"net/rpc"

	"github.com/MrWong99/adventofcode/proto"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

const (
	GRPCPluginKey = "calc_grpc"
	RPCPluginKey  = "calc"
)

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "Advent",
	MagicCookieValue: "OfCode",
}

var PluginMap = map[string]plugin.Plugin{
	GRPCPluginKey: &CalcGRPCPlugin{},
	RPCPluginKey:  &CalcPlugin{},
}

type CalcService interface {
	Calculate(input string) (string, error)
}

type CalcPlugin struct {
	Impl CalcService
}

func (p *CalcPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (*CalcPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}

type CalcGRPCPlugin struct {
	plugin.Plugin
	Impl CalcService
}

func (p *CalcGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterCalculateServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *CalcGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewCalculateClient(c)}, nil
}
