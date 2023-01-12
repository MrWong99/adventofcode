package main

import (
	"strings"

	"github.com/MrWong99/adventofcode/pluginmanager/shared"
	"github.com/hashicorp/go-plugin"
)

type CalcService struct{}

func (CalcService) Calculate(input string) (string, error) {
	return strings.ToUpper(input), nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"kv": &shared.CalcGRPCPlugin{Impl: &CalcService{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
