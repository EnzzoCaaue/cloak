package command

import (
	"github.com/Cloakaac/cloak/util"
	"github.com/raggaer/pigo"
	"log"
)

type reloadItems struct{}

func init() {
	commands.Add("reload items", &reloadItems{})
}

func (r *reloadItems) exec() {
	util.ParseItems(pigo.Config.String("datapack"))
	log.Println("Items loaded")
}
