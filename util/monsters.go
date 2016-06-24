package util

import (
	"encoding/xml"
	"io/ioutil"
	"log"
)

// Monster is the main monster representation
type Monster struct {
	XMLName         xml.Name          `xml:"monster"`
	Name            string            `xml:"name,attr"`
	NameDescription string            `xml:"nameDescription,attr"`
	Race            string            `xml:"race,attr"`
	Experience      int               `xml:"experience,attr"`
	Speed           int               `xml:"speed,attr"`
	ManaCost        int               `xml:"manacost,attr"`
	Health          MonsterHealth     `xml:"health"`
	Look            MonsterLook       `xml:"look"`
	Voices          []MonsterSentence `xml:"voices>voice"`
	Loot            []MonsterItem     `xml:"loot>item"`
}

// MonsterItem struct for monster loot
type MonsterItem struct {
	ID       int `xml:"id,attr"`
	CountMax int `xml:"countmax,attr"`
	Chance   int `xml:"chance,attr"`
}

// MonsterSentence struct for monster talks
type MonsterSentence struct {
	Sentence string `xml:"sentence,attr"`
}

// MonsterLook struct for monster looktype
type MonsterLook struct {
	Type   int `xml:"type,attr"`
	Head   int `xml:"head,attr"`
	Body   int `xml:"body,attr"`
	Legs   int `xml:"legs,attr"`
	Feet   int `xml:"feet,attr"`
	Corpse int `xml:"corpse,attr"`
}

// MonsterHealth struct for monster stats
type MonsterHealth struct {
	Now int `xml:"now,attr"`
	Max int `xml:"max,attr"`
}

type monsterDef struct {
	Name string `xml:"name,attr"`
	File string `xml:"file,attr"`
}

type monsterDefinition struct {
	XMLName  xml.Name     `xml:"monsters"`
	Monsters []monsterDef `xml:"monster"`
}

// ParseMonsters parses monsters.xml
func ParseMonsters(path string) {
	monsters = make(map[string]*Monster)
	b, err := ioutil.ReadFile(path + "/data/monster/monsters.xml")
	if err != nil {
		log.Fatal(err)
	}
	definitions := monsterDefinition{}
	err = xml.Unmarshal(b, &definitions)
	if err != nil {
		log.Fatal(err)
	}
	for _, monster := range definitions.Monsters {
		parseMonster(path, monster.Name, monster.File)
	}
}

func parseMonster(path, name, file string) {
	b, err := ioutil.ReadFile(path + "/data/monster/" + file)
	if err != nil {
		log.Println("Error while parsing monster:", name)
		return
	}
	m := &Monster{}
	err = xml.Unmarshal(b, &m)
	if err != nil {
		log.Fatal(err)
	}
	monsters[name] = m
}

// GetMonster returns a monster struct
func GetMonster(name string) *Monster {
	if v, ok := monsters[name]; ok {
		return v
	}
	return nil
}
