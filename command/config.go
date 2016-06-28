package command

import (
	"github.com/Cloakaac/cloak/util"
	"github.com/raggaer/pigo"
	"log"
)

type reloadConfig struct{}

type reloadConfigLUA struct{}

func init() {
	commands.Add("reload config", &reloadConfig{})
	commands.Add("reload config lua", &reloadConfigLUA{})
}

func (r *reloadConfig) exec() {
	pigo.ParseConfig("config.json")
	log.Println("Config loaded")
}

func (r *reloadConfigLUA) exec() {
	util.ParseConfig(pigo.Config.String("datapack"))
	log.Println("Config LUA loaded")
}
