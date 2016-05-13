package util

import (
	"io/ioutil"
	"encoding/json"
)

// Parser stores the database config values
var Parser = &Config{}

// Config is the main config structure of the AAC
type Config struct {
	Captcha *cptcha
	Port string
	Mode int
	Database *connection
	Template string
	Routes []*luaRoute
	Register *register
}

type register struct {
	Level int
	Premdays int
	Stamina int
	Experience int
	Health int
	Healthmax int
	Mana int
	Manamax int
	Male *registerMale
	Female *registerFemale
	Skills *registerSkills
}

type registerSkills struct {
	Axe int
	Sword int
	Club int
	Dist int
	Fish int
	Fist int
	Shield int
}

type registerMale struct {
	Lookbody int
	Lookfeet int
	Lookhead int
	Looktype int
	Lookaddons int
}

type registerFemale struct {
	Lookbody int
	Lookfeet int
	Lookhead int
	Looktype int
	Lookaddons int
}

type cptcha struct {
	Public string
	Secret string
}

type connection struct {
	User string
	Password string
	Database string
}

type luaRoute struct {
	Path string
    Method string
    File string
    Mode string
}

// LoadConfig loads and parsers a config.json file
func LoadConfig() error {
	f, err := ioutil.ReadFile("config.json")
    if err != nil {
        return err
    }
    err = json.Unmarshal(f, Parser)
    if err != nil {
        return err
    }
    return nil
}