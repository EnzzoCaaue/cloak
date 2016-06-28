package command

import (
	"bufio"
	"errors"
	"log"
	"os"
	"sync"
)

var (
	commands = &cloakaCommands{
		make(map[string]command),
		&sync.RWMutex{},
	}
)

type cloakaCommands struct {
	list map[string]command
	rw   *sync.RWMutex
}

type command interface {
	exec()
}

func (c *cloakaCommands) Add(arg string, cmd command) error {
	c.rw.Lock()
	defer c.rw.Unlock()
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
			continue
		}
		if v, ok := commands.list[string(cmd)]; ok {
			v.exec()
			continue
		}
		log.Println("Command not found")
	}
}
