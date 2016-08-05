package command

import (
	"log"

	"github.com/Cloakaac/cloak/util"
	"github.com/spf13/viper"
)

type reloadMap struct{}

func init() {
	commands.Add("reload map", &reloadMap{})
}

func (r *reloadMap) exec() {
	util.ParseMap(viper.GetString("datapack"))
	log.Println("Map loaded")
}
