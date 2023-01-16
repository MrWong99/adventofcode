package db

import (
	"os"

	"github.com/MrWong99/adventofcode/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Plugin struct {
	gorm.Model
	Name string
	Cmd  string
}

var db *gorm.DB

var cachedPlugins = make([]*Plugin, 0)

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("plugin.db"), &gorm.Config{})
	if err != nil {
		log.Log.Error("Could not initialize database", "error", err)
		os.Exit(1)
	}
	err = db.AutoMigrate(&Plugin{})
	if err != nil {
		log.Log.Error("Could not migrate database", "error", err)
		os.Exit(1)
	}
	cachedPlugins, err = readPlugins()
	if err != nil {
		log.Log.Error("Could not read database", "error", err)
		os.Exit(1)
	}
}

func readPlugins() ([]*Plugin, error) {
	var plugins []*Plugin
	res := db.Find(&plugins)
	return plugins, res.Error
}

func FindPlugin(id uint) (*Plugin, error) {
	for _, plugin := range cachedPlugins {
		if plugin.ID == id {
			return plugin, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func AllPlugins() []*Plugin {
	return cachedPlugins
}

func AddPlugin(plugin *Plugin) error {
	res := db.Create(plugin)
	if res.Error != nil {
		return res.Error
	}
	cachedPlugins = append(cachedPlugins, plugin)
	return nil
}

func DeletePlugin(id uint) (*Plugin, error) {
	plugin := &Plugin{}
	res := db.Clauses(clause.Returning{}).Where("ID = ?", id).Delete(plugin)
	if res.Error != nil {
		return plugin, res.Error
	}
	index := 0
	for i, p := range cachedPlugins {
		if p.ID == id {
			index = i
			break
		}
	}
	cachedPlugins[index] = cachedPlugins[len(cachedPlugins)-1]
	cachedPlugins = cachedPlugins[:len(cachedPlugins)-1]
	return plugin, nil
}
