package command

import (
	"os"
)

type exitCloaka struct{}

func init() {
	commands.Add("exit", &exitCloaka{})
}

func (r *exitCloaka) exec() {
	os.Exit(0)
}
