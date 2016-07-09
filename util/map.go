package util

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
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

func parseHouses(path, houseFile string) error {
	Houses.rw.Lock()
	defer Houses.rw.Unlock()
	Houses.List = &HouseList{}
	b, err := ioutil.ReadFile(path + "/data/world/" + houseFile)
	if err != nil {
		return err
	}
	return xml.Unmarshal(b, Houses.List)
}

// ParseMap loads and parses the given OTBM file
func ParseMap(path string) error {
	serverMap := &otmap.Map{}
	serverMap.Initialize()
	otbLoader := &otmap.OtbLoader{}
	otbLoader.Load(path + "/data/items/items.otb")
	if err := serverMap.ReadOTBM(path+"/data/world/forgotten.otbm", otbLoader); err != nil {
		return err
	}
	parseHouses(path, serverMap.HouseFile)
	for _, h := range serverMap.Houses {
		houseData := Houses.GetHouse(h.ID)
		houseImage := image.NewRGBA(image.Rect(int(houseData.EntryX)-20, int(houseData.EntryY)-20, int(houseData.EntryX)+20, int(houseData.EntryY)+20))
		houseImage.Set(int(houseData.EntryX), int(houseData.EntryY), color.RGBA{
			255,
			0,
			0,
			255,
		})
		for _, tile := range h.Tiles {
			pos := tile.Position()
			houseImage.Set(int(pos.X), int(pos.Y), color.RGBA{
				216,
				200,
				10,
				255,
			})
		}
		imgFile, _ := os.Create(fmt.Sprintf("%v/%v/%v.png", pigo.Config.String("template"), "public/houses", houseData.Name))
		defer imgFile.Close()
		png.Encode(imgFile, houseImage)
	}
	return nil
}
