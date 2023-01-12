package shared

import (
	"context"

	"github.com/MrWong99/adventofcode/proto"
)

type GRPCClient struct {
	client proto.CalculateClient
}

func (c *GRPCClient) Calculate(input string) (result string, err error) {
	var resp *proto.ResultResponse
	resp, err = c.client.Calculate(context.Background(), &proto.CalcRequest{Input: input})
	if err != nil {
		return
	}
	result = resp.Value
	return
}

type GRPCServer struct {
	proto.UnimplementedCalculateServer
	Impl CalcService
}

func (s *GRPCServer) Calculate(ctx context.Context, req *proto.CalcRequest) (resp *proto.ResultResponse, err error) {
	var val string
	val, err = s.Impl.Calculate(req.Input)
	if err == nil {
		resp = &proto.ResultResponse{Value: val}
	}
	return
}
