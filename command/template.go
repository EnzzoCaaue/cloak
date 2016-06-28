package command

import (
    "github.com/Cloakaac/cloak/template"
	"log"
)

type reloadTemplate struct {}

func init() {
    commands.Add("reload template", &reloadTemplate{})
}

func (r *reloadTemplate) exec() {
    template.Load()
    log.Println("Template loaded")
}