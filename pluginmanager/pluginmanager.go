package pluginmanager

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/MrWong99/adventofcode/pluginmanager/shared"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

var calculators = make(map[string]*CalculatorPlugin, 0)

func Calculate(calculatorName, input string) (string, error) {
	if calc, ok := calculators[calculatorName]; ok {
		return calc.service.Calculate(input)
	} else {
		return "", fmt.Errorf("calculator with name %q is not registered", calculatorName)
	}
}

func DeleteCalculator(name string) {
	delete(calculators, name)
}

func RegisterCalculator(calc *CalculatorPlugin) error {
	if calc.Name == "" {
		return errors.New("a calculator must have a non-empty name")
	}
	if calc.Cmd == "" {
		return errors.New("a calculator must have a non-empty cmd")
	}
	if calc.client == nil || calc.service == nil {
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig:  shared.Handshake,
			Plugins:          shared.PluginMap,
			Cmd:              exec.Command("sh", "-c", calc.Cmd),
			AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
			Logger: hclog.New(&hclog.LoggerOptions{
				Name:   "plugin-" + calc.Name,
				Level:  hclog.Info,
				Output: os.Stderr,
				Color:  hclog.AutoColor,
			}),
		})
		rpcClient, err := client.Client()
		if err != nil {
			client.Kill()
			return fmt.Errorf("could not start client: %v", err)
		}

		plugin, err := rpcClient.Dispense(shared.GRPCPluginKey)
		if err != nil {
			client.Kill()
			return fmt.Errorf("could not request grpc plugin: %v", err)
		}

		service, ok := plugin.(shared.CalcService)
		if !ok {
			client.Kill()
			return fmt.Errorf("plugin did not implement calc interface: %T", plugin)
		}
		calc.client = client
		calc.service = service
	}
	if c, ok := calculators[calc.Name]; ok {
		c.Close()
	}
	calculators[calc.Name] = calc
	return nil
}

func Close() {
	for _, c := range calculators {
		c.Close()
	}
	calculators = make(map[string]*CalculatorPlugin, 0)
}

type CalculatorPlugin struct {
	Name    string
	Cmd     string
	client  *plugin.Client
	service shared.CalcService
}

func (c *CalculatorPlugin) Calculate(input string) (string, error) {
	if c.service == nil {
		return "", fmt.Errorf("calculator %q is not yet started and needs to be registered", c.Name)
	}
	return c.service.Calculate(input)
}

func (c *CalculatorPlugin) Close() {
	if c.client != nil {
		c.client.Kill()
		c.client = nil
	}
	c.service = nil
}
