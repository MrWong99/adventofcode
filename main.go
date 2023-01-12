package main

import (
	"flag"
	"os"
	"strings"

	"github.com/MrWong99/adventofcode/pluginmanager"
	"github.com/hashicorp/go-hclog"
)

var log = hclog.New(&hclog.LoggerOptions{
	Name:            "main",
	Level:           hclog.Info,
	Output:          os.Stdout,
	IncludeLocation: true,
	Color:           hclog.AutoColor,
})

var CalcsToLoad *string = flag.String(
	"calculator-plugins", "", "an optional list of calculator plugins in the format 'calc1Name=calc1Command,calc2Name=calc2Command,...'",
)

var InputFile *string = flag.String(
	"input-file", "", "when provided: the file in given path will be loaded and used as input",
)

var CalcToUse *string = flag.String(
	"calculate-with-plugin", "", "when provided: the result will be calculated via CLI and the entire execution terminates (no server)",
)

func main() {
	os.Exit(runWithExit())
}

func runWithExit() int {
	flag.Parse()
	manager := pluginmanager.NewManager()
	defer manager.Close()

	if *CalcsToLoad == "" {
		log.Info("No calculation-plugins provided")
	} else {
		for _, calc := range strings.Split(*CalcsToLoad, ",") {
			parts := strings.SplitN(calc, "=", 2)
			if len(parts) != 2 {
				log.Warn("The calculator-plugin is not valid", "plugin", calc)
				continue
			}
			if err := loadPlugin(manager, parts[0], parts[1]); err != nil {
				log.Warn("Could not load plugin", "name", parts[0], "cmd", parts[1], "error", err)
			}
		}
	}

	if *CalcToUse == "" {
		log.Info("No CLI calculation requested. Starting server now...")
	} else {
		log.Info("CLI calculation with plugin requested", "plugin", *CalcToUse)
		input := ""
		if *InputFile == "" {
			log.Warn("No '-input-file' provided. Using empty string as calculation input")
		} else {
			content, err := os.ReadFile(*InputFile)
			if err != nil {
				log.Error("Could not read input file", "-input-file", *InputFile, "error", err)
				return 1
			}
			input = string(content)
			log.Info("Successfully loaded input file", "-input-file", *InputFile)
			log.Debug("Input file content", "-input-file", *InputFile, "content", input)
		}
		result, err := manager.Calculate(*CalcToUse, input)
		if err == nil {
			log.Info("Calculation successful", "result", result)
		} else {
			log.Error("Calculation failed", "error", err)
			return 2
		}
	}
	return 0
}

func loadPlugin(manager *pluginmanager.Manager, name, cmd string) error {
	return manager.RegisterCalculator(&pluginmanager.Calculator{Name: name, Cmd: cmd})
}