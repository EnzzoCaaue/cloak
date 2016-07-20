package command

import (
	"log"

	"github.com/Cloakaac/cloak/util"
	"github.com/spf13/viper"
)

type reloadMonster struct{}

func init() {
	commands.Add("reload monsters", &reloadMonster{})
}

func (r *reloadMonster) exec() {
	util.ParseMonsters(viper.GetString("datapack"))
	log.Println("Monsters loaded")
}
