package command

import (
	"github.com/Cloakaac/cloak/util"
	"github.com/raggaer/pigo"
	"log"
)

type reloadStages struct{}

func init() {
	commands.Add("reload stages", &reloadStages{})
}

func (r *reloadStages) exec() {
	util.ParseStages(pigo.Config.String("datapack"))
	log.Println("Stages loaded")
}
