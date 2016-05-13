package main

import (
	"github.com/Cloakaac/cloak/database"
	"github.com/Cloakaac/cloak/template"
	"github.com/Cloakaac/cloak/util"
	"github.com/Cloakaac/cloak/web"
	"log"
)

func main() {
	err := util.LoadConfig()
	if err != nil {
		log.Println("Cannot parse config.json")
		log.Fatal(err)
	}
	log.Println("Config values parsed")
	err = database.NewConnection(util.Parser.Database.User, util.Parser.Database.Password, util.Parser.Database.Database)
	if err != nil {
		log.Println("An error occured while connecting to MySQL database")
		log.Fatal(err)
	}
	log.Println("MySQL connection estabilished")
	err = template.NewRender(util.Parser.Template)
	if err != nil {
		log.Println("An error occured while parsing Cloaka template")
		log.Fatal(err)
	}
	log.Println("Template renderer registered")
	util.SetMode(util.Parser.Mode)
	if util.Parser.Mode == 0 {
		log.Println("Running Cloaka on DEBUG mode")
	} else {
		log.Println("Running Cloaka on RELEASE mode")
	}
	web.Start(util.Parser.Port, util.Parser.Template)
}
