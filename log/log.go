package log

import (
	"os"

	"github.com/hashicorp/go-hclog"
)

var Log = hclog.New(&hclog.LoggerOptions{
	Name:            "adventofcode-main",
	Level:           hclog.Info,
	Output:          os.Stdout,
	IncludeLocation: true,
	Color:           hclog.AutoColor,
})
