package util

import (
    "log"
    "io/ioutil"
    "encoding/xml"
)

type monster struct {
    XMLName xml.Name `xml:"monster"`
    Name string `xml:"name,attr"`
    NameDescription string `xml:"nameDescription,attr"`
    Race string `xml:"race,attr"`
    Experience int `xml:"experience,attr"`
    Speed int `xml:"speed,attr"`
    ManaCost int `xml:"manacost,attr"`
    Health monsterHealth `xml:"health"`
    Look monsterLook `xml:"look"`
    Voices []monsterSentence `xml:"voices>voice"`
    Loot []monsterItem `xml:"loot>item"`
}

type monsterItem struct {
    ID int `xml:"id,attr"`
    CountMax int `xml:"countmax,attr"`
    Chance int `xml:"chance,attr"`
}

type monsterSentence struct {
    Sentence string `xml:"sentence,attr"`
}

type monsterLook struct {
    Type int `xml:"type,attr"`
    Head int `xml:"head,attr"`
    Body int `xml:"body,attr"`
    Legs int `xml:"legs,attr"`
    Feet int `xml:"feet,attr"`
    Corpse int `xml:"corpse,attr"`
}

type monsterHealth struct {
    Now int `xml:"now,attr"`
    Max int `xml:"max,attr"`
}

type monsterDef struct {
    Name string `xml:"name,attr"`
    File string `xml:"file,attr"`
}

type monsterDefinition struct {
    XMLName xml.Name `xml:"monsters"`
    Monsters []monsterDef `xml:"monster"`
}

// ParseMonsters parses monsters.xml
func ParseMonsters(path string) {
    monsters = []*monster{}
    b, err := ioutil.ReadFile(path+"/data/monster/monsters.xml")
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
    b, err := ioutil.ReadFile(path+"/data/monster/"+file)
    if err != nil {
        log.Println("Error while parsing monster:", name)
        return
    }
    m := &monster{}
    err = xml.Unmarshal(b, &m)
    if err != nil {
        log.Fatal(err)
    }
    monsters = append(monsters, m)
}