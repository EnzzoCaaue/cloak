package otmap

import "sync"

type House struct {
	ID      uint32
	DoorPos Position
	Tiles   []Tile
}

type Town struct {
	ID        uint32
	Name      string
	TemplePos Position
}

type Map struct {
	Width  uint16
	Height uint16

	Description string
	HouseFile   string
	SpawnFile   string

	Tiles     map[Position]Tile
	Houses    []*House
	Towns     []Town
	Waypoints map[Position]string
	rw        *sync.RWMutex
}

func (otMap *Map) Initialize() {
	otMap.Tiles = make(map[Position]Tile)
	otMap.Waypoints = make(map[Position]string)
	otMap.rw = &sync.RWMutex{}
}

func (otMap *Map) getHouse(id uint32) *House {
	otMap.rw.RLock()
	defer otMap.rw.RUnlock()
	for i := range otMap.Houses {
		house := otMap.Houses[i]
		if house.ID == id {
			return house
		}
	}
	return nil
}

func (otMap *Map) addHouse(house *House) error {
	otMap.rw.Lock()
	defer otMap.rw.Unlock()
	otMap.Houses = append(otMap.Houses, house)
	return nil
}

// GetTile returns a tile
func (otMap *Map) GetTile(x uint16, y uint16, z uint8) Tile {
	return otMap.Tiles[Position{x, y, z}]
}
