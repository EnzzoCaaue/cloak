package util

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/Cloakaac/cloak/otmap"
	"github.com/raggaer/pigo"
)

// House holds all information about a game house
type House struct {
	ID     uint32 `xml:"houseid,attr"`
	Name   string `xml:"name,attr"`
	EntryX uint16 `xml:"entryx,attr"`
	EntryY uint16 `xml:"entryy,attr"`
	EntryZ uint16 `xml:"entryz,attr"`
	Size   int    `xml:"size,attr"`
	TownID int    `xml:"townid,attr"`
}

// HouseList holds the house array
type HouseList struct {
	XMLName xml.Name `xml:"houses"`
	Houses  []*House `xml:"house"`
}

// ServerHouses contains the whole house list of the server
type ServerHouses struct {
	List *HouseList
	rw   *sync.RWMutex
}

type ServerTowns struct {
	List []otmap.Town
	rw   *sync.RWMutex
}

// Get returns a town by its name
func (s *ServerTowns) Get(name string) *otmap.Town {
	s.rw.RLock()
	defer s.rw.RUnlock()
	for _, town := range s.List {
		if town.Name == name {
			return &town
		}
	}
	return nil
}

// GetByID returns a town by its ID
func (s *ServerTowns) GetByID(id uint32) string {
	s.rw.RLock()
	defer s.rw.RUnlock()
	for _, town := range s.List {
		if town.ID == id {
			return town.Name
		}
	}
	return ""
}

// GetList returns the whole town list
func (s *ServerTowns) GetList() []otmap.Town {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.List
}

// Exists checks if a town is valid
func (s *ServerTowns) Exists(name string) bool {
	s.rw.RLock()
	defer s.rw.RUnlock()
	for _, town := range s.List {
		if town.Name == name {
			return true
		}
	}
	return false
}

// GetHouse gets a house by its ID
func (s *ServerHouses) GetHouse(id uint32) *House {
	s.rw.RLock()
	defer s.rw.RUnlock()
	for _, house := range s.List.Houses {
		if house.ID == id {
			return house
		}
	}
	return nil
}

// GetHouseByName gets a house by its name
func (s *ServerHouses) GetHouseByName(name string) *House {
	s.rw.RLock()
	defer s.rw.RUnlock()
	for _, house := range s.List.Houses {
		if house.Name == name {
			return house
		}
	}
	return nil
}

func parseHouses(path, houseFile string) error {
	Houses.rw.Lock()
	defer Houses.rw.Unlock()
	b, err := ioutil.ReadFile(path + "/data/world/" + houseFile)
	if err != nil {
		return err
	}
	return xml.Unmarshal(b, Houses.List)
}

// ParseMap loads and parses the given OTBM file
func ParseMap(path string) {
	serverMap := &otmap.Map{}
	serverMap.Initialize()
	otbLoader := &otmap.OtbLoader{}
	otbLoader.Load(path + "/data/items/items.otb")
	if err := serverMap.ReadOTBM(path+"/data/world/forgotten.otbm", otbLoader, false); err != nil {
		log.Fatal(err)
	}
	if err := parseHouses(path, serverMap.HouseFile); err != nil {
		log.Fatal(err)
	}
	Towns.rw.Lock()
	defer Towns.rw.Unlock()
	Towns.List = serverMap.Towns
	tileColor := color.RGBA{192, 192, 192, 255}
	wallColor := color.RGBA{255, 0, 0, 255}
	doorColor := color.RGBA{255, 255, 0, 255}
	backgroundColor := color.RGBA{0, 0, 0, 255}
	for _, h := range serverMap.Houses {
		houseData := Houses.GetHouse(h.ID)
		houseImage := image.NewRGBA(image.Rect(int(houseData.EntryX)-32, int(houseData.EntryY)-32, int(houseData.EntryX)+32, int(houseData.EntryY)+32))
		draw.Draw(houseImage, houseImage.Bounds(), &image.Uniform{
			backgroundColor,
		}, image.ZP, draw.Src)
		houseTiles := make(map[otmap.Position]bool, len(h.Tiles))
		for _, tile := range h.Tiles {
			pos := tile.Position()
			if pos.Z != uint8(houseData.EntryZ) {
				continue
			}
			houseImage.Set(int(pos.X), int(pos.Y), tileColor)
			houseTiles[pos] = true
		}
		for pos := range houseTiles {
			if houseImage.At(int(pos.X)+1, int(pos.Y)+1) != tileColor {
				houseImage.Set(int(pos.X)+1, int(pos.Y)+1, wallColor)
			}
			if houseImage.At(int(pos.X)+1, int(pos.Y)-1) != tileColor {
				houseImage.Set(int(pos.X)+1, int(pos.Y)-1, wallColor)
			}
			if houseImage.At(int(pos.X)-1, int(pos.Y)+1) != tileColor {
				houseImage.Set(int(pos.X)-1, int(pos.Y)+1, wallColor)
			}
			if houseImage.At(int(pos.X)-1, int(pos.Y)-1) != tileColor {
				houseImage.Set(int(pos.X)-1, int(pos.Y)-1, wallColor)
			}
		}
		houseImage.Set(int(houseData.EntryX), int(houseData.EntryY), doorColor)
		imgFile, _ := os.Create(fmt.Sprintf("%v/%v/%v.png", pigo.Config.String("template"), "public/houses", houseData.Name))
		png.Encode(imgFile, houseImage)
		imgFile.Close()
	}
}
