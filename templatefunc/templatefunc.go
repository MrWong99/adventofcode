package templatefunc

import (
	"strconv"

	"github.com/MrWong99/adventofcode/db"
	"github.com/MrWong99/adventofcode/pluginmanager"
)

func AddPlugin(name, cmd string) (err error) {
	if err = pluginmanager.RegisterCalculator(&pluginmanager.CalculatorPlugin{Name: name, Cmd: cmd}); err != nil {
		return
	}
	if err = db.AddPlugin(&db.Plugin{Name: name, Cmd: cmd}); err != nil {
		pluginmanager.DeleteCalculator(name)
	}
	return
}

func DeletePlugin(id string) error {
	uid, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		return err
	}
	plugin, err := db.DeletePlugin(uint(uid))
	if err != nil {
		return err
	}
	pluginmanager.DeleteCalculator(plugin.Name)
	return nil
}

func Calculate(id, input string) (string, error) {
	uid, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		return "", err
	}
	plugin, err := db.FindPlugin(uint(uid))
	if err != nil {
		return "", err
	}
	return pluginmanager.Calculate(plugin.Name, input)
}
