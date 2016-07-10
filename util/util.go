package util

import (
	"fmt"
	"sync"
	"time"

	"github.com/Cloakaac/cloak/otmap"
)

const (
	novoc = iota
	sorcerer
	druid
	paladin
	knight
	masterSorcerer
	elderDruid
	royalPaladin
	eliteKnight
)

const (
	female = iota
	male
)

var (
	// Mode stores the AAC mode
	Mode       int
	genderList = map[string]int{
		"Male":   male,
		"Female": female,
	}
	vocationList = map[string]int{
		"Sorcerer":        sorcerer,
		"Druid":           druid,
		"Paladin":         paladin,
		"Knight":          knight,
		"Master Sorcerer": masterSorcerer,
		"Elder Druid":     elderDruid,
		"Royal Paladin":   royalPaladin,
		"Elite Knight":    eliteKnight,
	}
	// Monsters holds all server monsters
	Monsters = &ServerMonsters{
		m:  make(map[string]*Monster),
		rw: &sync.RWMutex{},
	}
	// Config contains the whole parsed config lua file
	Config = &ConfigLUA{
		make(map[string]interface{}),
		&sync.RWMutex{},
	}
	// Stages contains the server experience stages
	Stages = &ServerStages{
		&stageDefinition{},
		&sync.RWMutex{},
	}
	// Items contains the server items xml parsed
	Items = &ServerItems{
		make(map[int]ItemDefinition),
		&sync.RWMutex{},
	}
	// Houses contains the server houses.xml
	Houses = &ServerHouses{
		&HouseList{},
		&sync.RWMutex{},
	}
	// Towns contains the server town list
	Towns = &ServerTowns{
		[]otmap.Town{},
		&sync.RWMutex{},
	}
)

// SetMode sets the AAC run mode DEBUG(0) RELEASE(1)
func SetMode(mode int) {
	Mode = mode
}

// Vocation gets the vocation id from a given string
func Vocation(voc string) int {
	return vocationList[voc]
}

// Gender gets the gender id from a given string
func Gender(gender string) int {
	return genderList[gender]
}

// GetGender gets the gender string from a id
func GetGender(gender int) string {
	if gender == 0 {
		return "Female"
	}
	return "Male"
}

// GetVocation gets the vocation string from a id
func GetVocation(voc int) string {
	for i, v := range vocationList {
		if v == voc {
			return i
		}
	}
	return "Sorcerer"
}

// UnixToString converts a int64 to a string date
func UnixToString(unix int64) string {
	if unix == 0 {
		return "Never"
	}
	timeDate := time.Unix(unix, 0)
	timeString := fmt.Sprintf("%v %v %v, %v:%v:%v", timeDate.Month().String()[:3], timeDate.Day(), timeDate.Year(), timeDate.Hour(), timeDate.Minute(), timeDate.Second())
	return timeString
}
