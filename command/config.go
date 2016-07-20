package command

import (
	"log"

	"github.com/Cloakaac/cloak/util"
	"github.com/raggaer/pigo"
	"github.com/spf13/viper"
)

type reloadConfig struct{}

type reloadConfigLUA struct{}

func init() {
	commands.Add("reload config", &reloadConfig{})
	commands.Add("reload config lua", &reloadConfigLUA{})
}

func (r *reloadConfig) exec() {
	pigo.ParseConfig("config")
	log.Println("Config loaded")
}

func (r *reloadConfigLUA) exec() {
	util.ParseConfig(viper.GetString("datapack"))
	log.Println("Config LUA loaded")
}
