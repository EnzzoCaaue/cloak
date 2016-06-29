package daemon

import (
    "sync"
	"errors"
	"time"
)

var (
    daemons = &cloakaDaemons{
        make(map[string]*cloakaDaemon),
        &sync.RWMutex{},
    }
)

type cloakaDaemon struct {
    dm daemon
    duration *time.Ticker
}

type cloakaDaemons struct {
    list map[string]*cloakaDaemon
    rw *sync.RWMutex
}

type daemon interface{
    tick()
    name()
}

// RunDaemons runs all daemon tickers
func RunDaemons() {
    daemons.rw.RLock()
    defer daemons.rw.RUnlock()
    for _, d := range daemons.list {
        go runDaemon(d)
    }
}

func (c *cloakaDaemons) Add(key string, duration time.Duration, dm daemon) error {
    c.rw.Lock()
    defer c.rw.Unlock()
    if _, ok := c.list[key]; !ok {
        c.list[key] = &cloakaDaemon{
            dm,
            time.NewTicker(duration),
        }
        return nil
    }
    return errors.New("Daemon service already exists")
}

func runDaemon(dm *cloakaDaemon) {
    for {
        select {
        case <- dm.duration.C:
            dm.dm.name()
            dm.dm.tick()    
        }
    }
}