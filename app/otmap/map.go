package otmap

import "sync"

// House used to represent server houses
type House struct {
	ID      uint32
	DoorPos Position
	Tiles   []Tile
}

// Town used to represent server towns
type Town struct {
	ID        uint32
	Name      string
	TemplePos Position
}

// Map used to represent the whole server map
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

// Initialize creates an empty map
func (otMap *Map) Initialize() {
	otMap.rw = &sync.RWMutex{}
	otMap.rw.Lock()
	defer otMap.rw.Unlock()
	otMap.Tiles = make(map[Position]Tile)
	otMap.Waypoints = make(map[Position]string)
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
	otMap.rw.RLock()
	defer otMap.rw.RUnlock()
	return otMap.Tiles[Position{x, y, z}]
}
