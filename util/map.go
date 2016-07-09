package util

import "github.com/Cloakaac/cloak/otmap"

// ServerMap contains all the map information parsed
type ServerMap struct {
	Map *otmap.Map
}

// GetTown returns a town by its name
func (s *ServerMap) GetTown(name string) {
    s.Map.Towns[0].
}

// ParseMap loads and parses the given OTBM file
func ParseMap(path string) error {
	Map.Map.Initialize()
	otbLoader := &otmap.OtbLoader{}
	otbLoader.Load(path + "/data/items/items.otb")
	if err := Map.Map.ReadOTBM(path+"/data/world/forgotten.otbm", otbLoader); err != nil {
		return err
	}
	return nil
}
