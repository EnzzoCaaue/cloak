package command

import (
	"log"

	"github.com/Cloakaac/cloak/util"
	"github.com/raggaer/pigo"
)

type reloadMap struct{}

func init() {
	commands.Add("reload map", &reloadMap{})
}

func (r *reloadMap) exec() {
	util.ParseMap(pigo.Config.String("datapack"))
	log.Println("Map loaded")
}
