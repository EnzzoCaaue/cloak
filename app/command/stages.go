package command

import (
	"log"

	"github.com/Cloakaac/cloak/util"
	"github.com/spf13/viper"
)

type reloadStages struct{}

func init() {
	commands.Add("reload stages", &reloadStages{})
}

func (r *reloadStages) exec() {
	util.ParseStages(viper.GetString("datapack"))
	log.Println("Stages loaded")
}
