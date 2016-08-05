package util

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
)

const (
	stagesPath      = "data/XML/stages"
	stagesExtension = "xml"
)

// ServerStages holds all the server stages
type ServerStages struct {
	v  *stageDefinition
	rw *sync.RWMutex
}

type stageDefinition struct {
	XMLName xml.Name    `xml:"stages"`
	Config  stageConfig `xml:"config"`
	Stages  []StageDef  `xml:"stage"`
}

type stageConfig struct {
	Enabled bool `xml:"enabled,attr"`
}

// StageDef represents an experience stage
type StageDef struct {
	MinLevel   int `xml:"minlevel,attr"`
	MaxLevel   int `xml:"maxlevel,attr"`
	Multiplier int `xml:"multiplier,attr"`
}

// ParseStages loads server stages with the given path
func ParseStages(path string) {
	Stages.rw.Lock()
	defer Stages.rw.Unlock()
	Stages.v = &stageDefinition{}
	b, err := ioutil.ReadFile(fmt.Sprintf("%v/%v.%v", path, stagesPath, stagesExtension))
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
