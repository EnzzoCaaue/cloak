package main

import (
	"flag"
	"github.com/Cloakaac/cloak/database"
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/template"
	"github.com/Cloakaac/cloak/util"
	"github.com/Cloakaac/cloak/web"
	"log"
)

func main() {
	mode := flag.Int("mode", 0, "AAC usage mode DEBUG(0) RELEASE(1)")
	port := flag.String("port", "8080", "AAC web server port")
	user := flag.String("user", "root", "MySQL user")
	password := flag.String("password", "admin", "MySQL password")
	db := flag.String("database", "cloaka", "MySQL database")
	tpl := flag.String("template", "G:/Workspace/Go/src/github.com/Cloakaac/template-default", "AAC template path")
	flag.Parse()
	err := database.NewConnection(*user, *password, *db)
	if err != nil {
		log.Println("An error occured while connecting to MySQL database")
		log.Fatal(err)
	}
	log.Println("MySQL connection estabilished")
	err = template.NewRender(*tpl)
	if err != nil {
		log.Println("An error occured while parsing Cloaka template")
		log.Fatal(err)
	}
	models.GetConfig()
	log.Println("Config values parsed")
	log.Println("Template renderer registered")
	util.SetMode(*mode)
	if *mode == 0 {
		log.Println("Running Cloaka on DEBUG mode")
	} else {
		log.Println("Running Cloaka on RELEASE mode")
	}
	web.Start(*port, *tpl)
}
