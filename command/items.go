package command

import (
	"log"

	"github.com/Cloakaac/cloak/util"
	"github.com/spf13/viper"
)

type reloadItems struct{}

func init() {
	commands.Add("reload items", &reloadItems{})
}

func (r *reloadItems) exec() {
	util.ParseItems(viper.GetString("datapack"))
	log.Println("Items loaded")
}
