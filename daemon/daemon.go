package daemon

import (
	"errors"
	"log"
	"sync"
	"time"
)

var (
	// List holds all cloaka daemons
	List = &cloakaDaemons{
		make(map[string]*cloakaDaemon),
		&sync.RWMutex{},
	}
)

type cloakaDaemon struct {
	dm       daemon
	duration *time.Ticker
	stop     chan bool
}

type cloakaDaemons struct {
	list map[string]*cloakaDaemon
	rw   *sync.RWMutex
}

type daemon interface {
	tick()
}

func (c *cloakaDaemons) Add(key string, duration time.Duration, dm daemon) error {
	log.Println(key)
	c.rw.Lock()
	defer c.rw.Unlock()
	if _, ok := c.list[key]; !ok {
		daemon := &cloakaDaemon{
			dm,
			time.NewTicker(duration),
			make(chan bool, 1),
		}
		c.list[key] = daemon
		go runDaemon(daemon)
		return nil
	}
	return errors.New("Daemon service already exists")
}

func (c *cloakaDaemons) Stop(key string) error {
	c.rw.Lock()
	defer c.rw.Unlock()
	if daemon, ok := c.list[key]; ok {
		daemon.stop <- true
		return nil
	}
	return errors.New("Daemon service not found")
}

func runDaemon(dm *cloakaDaemon) {
	defer func() {
		dm.duration.Stop()
		close(dm.stop)
	}()
	for {
		select {
		case <-dm.duration.C:
			dm.dm.tick()
		case s := <-dm.stop:
			if s {
				return
			}
		}
	}
}
