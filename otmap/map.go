package otmap

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
	Houses    []House
	Towns     []Town
	Waypoints map[Position]string
}

func (otMap *Map) Initialize() {
	otMap.Tiles = make(map[Position]Tile)
	otMap.Waypoints = make(map[Position]string)
}

func (otMap *Map) getHouse(id uint32) *House {
	for i := range otMap.Houses {
		house := otMap.Houses[i]
		if house.ID == id {
			return &house
		}
	}

	return nil
}
