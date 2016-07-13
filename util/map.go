package util

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sync"

	"github.com/Cloakaac/cloak/otmap"
	"github.com/raggaer/pigo"
)

const (
	tileSize = 4
)

var (
	drawnPixels     = 0
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

// ServerTowns contains the whole town list of the server
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
			var houseTopX uint16
			var houseTopY uint16
			var houseImageX uint16
			var houseImageY uint16
			houseImageY = houseData.EntryY
			houseImageX = houseData.EntryX
			houseTopX = houseData.EntryX
			houseTopY = houseData.EntryY
			for _, tile := range h.Tiles {
				pos := tile.Position()
				if uint16(pos.Z) != houseData.EntryZ {
					continue
				}
				if pos.X < houseTopX || houseTopX == 0 {
					houseTopX = pos.X
				}
				if pos.Y < houseTopY || houseTopY == 0 {
					houseTopY = pos.Y
				}
				if pos.X > houseImageX || houseImageX == 0 {
					houseImageX = pos.X
				}
				if pos.Y > houseImageY || houseImageY == 0 {
					houseImageY = pos.Y
				}
			}
			houseImage := image.NewRGBA(image.Rect(0, 0, int((houseImageX-houseTopX)+3), int((houseImageY-houseTopY)+3)))
			houseOffset := 1
			for _, tile := range h.Tiles {
				pos := tile.Position()
				if uint16(pos.Z) != houseData.EntryZ {
					continue
				}
				x := int(pos.X-houseTopX) + houseOffset
				y := int(pos.Y-houseTopY) + houseOffset
				houseImage.Set(x, y, tileColor)
			}
			for x := 0; x < houseImage.Bounds().Dx(); x++ {
				for y := 0; y < houseImage.Bounds().Dy(); y++ {
					if houseImage.At(x, y) == tileColor {
						continue
					}
					if houseImage.At(x+1, y+1) == tileColor {
						houseImage.Set(x, y, wallColor)
					}
					if houseImage.At(x-1, y-1) == tileColor {
						houseImage.Set(x, y, wallColor)
					}
					if houseImage.At(x+1, y) == tileColor {
						houseImage.Set(x, y, wallColor)
					}
					if houseImage.At(x, y+1) == tileColor {
						houseImage.Set(x, y, wallColor)
					}
					if houseImage.At(x-1, y) == tileColor {
						houseImage.Set(x, y, wallColor)
					}
					if houseImage.At(x, y-1) == tileColor {
						houseImage.Set(x, y, wallColor)
					}
					if houseImage.At(x+1, y-1) == tileColor {
						houseImage.Set(x, y, wallColor)
					}
					if houseImage.At(x-1, y+1) == tileColor {
						houseImage.Set(x, y, wallColor)
					}
				}
			}
			x := int(houseData.EntryX-houseTopX) + houseOffset
			y := int(houseData.EntryY-houseTopY) + houseOffset
			houseImage.Set(x, y, doorColor)
			houseFile, err := os.Create(fmt.Sprintf("%v/%v/%v.png", pigo.Config.String("template"), "public/houses", houseData.Name))
			if err != nil {
				log.Fatal(err)
			}
			png.Encode(houseFile, houseImage)
			houseFile.Close()
			waitHouses.Done()
		}(h)
	}
	waitHouses.Wait()
}

/*if houseImage.At(int(pos.X)+1, int(pos.Y)+1) != tileColor {
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
}*/
