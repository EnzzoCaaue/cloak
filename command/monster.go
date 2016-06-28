package command

import (
    "github.com/raggaer/pigo"
    "github.com/Cloakaac/cloak/util"
	"log"
)

type reloadMonster struct {}

func init() {
    commands.Add("reload monsters", &reloadMonster{})
}

func (r *reloadMonster) exec() {
    util.ParseMonsters(pigo.Config.String("datapack"))
    log.Println("Monsters loaded")
}