package util

import (
	"encoding/xml"
    "sync"
	"io/ioutil"
    "log"
)

type ServerStages struct {
    v *stageDefinition
    rw *sync.RWMutex
}

type stageDefinition struct {
	XMLName  xml.Name     `xml:"stages"`
    Config stageConfig `xml:"config"`
	Stages []StageDef `xml:"stage"`
}

type stageConfig struct {
    Enabled bool `xml:"enabled,attr"`
}

type StageDef struct {
    MinLevel int `xml:"minlevel,attr"`
    MaxLevel int `xml:"maxlevel,attr"`
    Multiplier int `xml:"multiplier,attr"`
}

// ParseStages loads server stages with the given path
func ParseStages(path string) {
    Stages.rw.RLock()
    defer Stages.rw.RUnlock()
    Stages.v = &stageDefinition{}
    b, err := ioutil.ReadFile(path+ "/data/XML/stages.xml")
    if err != nil {
        log.Fatal(err)
    }
    err = xml.Unmarshal(b, Stages.v)
    if err != nil {
        log.Fatal(err)
    }
}

// IsEnabled checks if stages are enabled on the server
func (s *ServerStages) IsEnabled() bool {
    s.rw.RLock()
    defer s.rw.RUnlock()
    return s.v.Config.Enabled
}

// GetAll retrieves all stages
func (s *ServerStages) GetAll() []StageDef {
    s.rw.RLock()
    defer s.rw.RUnlock()
    return s.v.Stages
}

// Get returns an stage by its index
func (s *ServerStages) Get(index int) StageDef {
    s.rw.RLock()
    defer s.rw.RUnlock()
    return s.v.Stages[index]
}