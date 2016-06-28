package command

import (
	"github.com/Cloakaac/cloak/util"
	"github.com/raggaer/pigo"
	"log"
)

type reloadMonster struct{}

func init() {
	commands.Add("reload monsters", &reloadMonster{})
}

func (r *reloadMonster) exec() {
	util.ParseMonsters(pigo.Config.String("datapack"))
	log.Println("Monsters loaded")
}
