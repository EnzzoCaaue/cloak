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
	"math/rand"
	"os"
	"sync"

	"github.com/Cloakaac/cloak/otmap"
	"github.com/raggaer/pigo"
)

var (
	tileColor       = color.RGBA{192, 192, 192, 255}
	wallColor       = color.RGBA{255, 0, 0, 255}
	doorColor       = color.RGBA{255, 255, 0, 255}
	backgroundColor = color.RGBA{0, 0, 0, 255}
	otbPath         = "/data/items/items.otb"
	otbmPath        = "/data/world/"
	otbmExtension   = ".otbm"
)

// House holds all information about a game house
type House struct {
	ID     uint32 `xml:"houseid,attr"`
	Name   string `xml:"name,attr"`
	EntryX uint16 `xml:"entryx,attr"`
	EntryY uint16 `xml:"entryy,attr"`
	EntryZ uint16 `xml:"entryz,attr"`
	Size   int    `xml:"size,attr"`
	TownID uint32 `xml:"townid,attr"`
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

// GetRandom returns a random town
func (s *ServerTowns) GetRandom() otmap.Town {
	s.rw.RLock()
	defer s.rw.RUnlock()
	rng := rand.Intn(len(s.List))
	return s.List[rng]
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

// GetList returns the full list of houses by its town
func (s *ServerHouses) GetList(id uint32) []*House {
	list := []*House{}
	s.rw.RLock()
	defer s.rw.RUnlock()
	for _, house := range s.List.Houses {
		if house.TownID == id {
			list = append(list, house)
		}
	}
	return list
}

func parseHouses(path, houseFile string) error {
	Houses.rw.Lock()
	defer Houses.rw.Unlock()
	b, err := ioutil.ReadFile(path + otbmPath + houseFile)
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
	otbLoader.Load(path + otbPath)
	if err := serverMap.ReadOTBM(fmt.Sprintf("%v%v%v%v", path, otbmPath, Config.String("mapName"), otbmExtension), otbLoader, false); err != nil {
		log.Fatal(err)
	}
	if err := parseHouses(path, serverMap.HouseFile); err != nil {
		log.Fatal(err)
	}
	Towns.rw.Lock()
	defer Towns.rw.Unlock()
	Towns.List = serverMap.Towns
	waitHouses := &sync.WaitGroup{}
	waitHouses.Add(len(serverMap.Houses))
	for _, h := range serverMap.Houses {
		go func(h *otmap.House) {
			houseData := Houses.GetHouse(h.ID)
			houseImage := image.NewRGBA(image.Rect(int(houseData.EntryX)-32, int(houseData.EntryY)-32, int(houseData.EntryX)+32, int(houseData.EntryY)+32))
			draw.Draw(houseImage, houseImage.Bounds(), &image.Uniform{
				backgroundColor,
			}, image.ZP, draw.Src)
			houseTiles := make([]otmap.Position, len(h.Tiles))
			for _, tile := range h.Tiles {
				pos := tile.Position()
				if pos.Z != uint8(houseData.EntryZ) {
					continue
				}
				houseImage.Set(int(pos.X), int(pos.Y), tileColor)
				houseTiles = append(houseTiles, pos)
			}
			drawWalls(houseTiles, houseImage)
			houseImage.Set(int(houseData.EntryX), int(houseData.EntryY), doorColor)
			imgFile, err := os.Create(fmt.Sprintf("%v/%v/%v.png", pigo.Config.String("template"), "public/houses", houseData.Name))
			if err != nil {
				log.Fatal(err)
			}
			png.Encode(imgFile, houseImage)
			imgFile.Close()
			waitHouses.Done()
		}(h)
	}
	waitHouses.Wait()
}

func drawWalls(tiles []otmap.Position, houseImage *image.RGBA) {
	for _, pos := range tiles {
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
		if houseImage.At(int(pos.X)+1, int(pos.Y)) != tileColor {
			houseImage.Set(int(pos.X)+1, int(pos.Y), wallColor)
		}
		if houseImage.At(int(pos.X), int(pos.Y)+1) != tileColor {
			houseImage.Set(int(pos.X), int(pos.Y)+1, wallColor)
		}
		if houseImage.At(int(pos.X)-1, int(pos.Y)) != tileColor {
			houseImage.Set(int(pos.X)-1, int(pos.Y), wallColor)
		}
		if houseImage.At(int(pos.X), int(pos.Y)-1) != tileColor {
			houseImage.Set(int(pos.X), int(pos.Y)-1, wallColor)
		}
	}
}
