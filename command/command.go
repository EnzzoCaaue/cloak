package command

import (
	"bufio"
	"errors"
	"os"
	"log"
)

var (
    commands = &cloakaCommands{
		make(map[string]command),
	}
)

type cloakaCommands struct {
	list map[string]command
}

type command interface {
	exec()
}

func (c *cloakaCommands) Add(arg string, cmd command) error {
	if _, ok := c.list[arg]; !ok {
		c.list[arg] = cmd
		return nil
	}
	return errors.New("Command already exists")
}

// ConsoleWatch watchs the console stdin
func ConsoleWatch() {
	reader := bufio.NewReader(os.Stdin)
	for {
		cmd, _, err := reader.ReadLine()
		if err != nil {
			log.Println("Error while reading command from stdin")
		}
		if v, ok := commands.list[string(cmd)]; ok {
			v.exec()
			continue
		}
		log.Println("Command not found")
	}
}