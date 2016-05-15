package models

import (
	"github.com/Cloakaac/cloak/database"
)

// Town struct for towns database
type Town struct {
	ID     int64
	Name   string
	TownID int64
}

// GetTowns gets all towns
func GetTowns() ([]*Town, error) {
	rows, err := database.Connection.Query("SELECT name, town_id FROM cloaka_towns ORDER BY id DESC")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	towns := []*Town{}
	for rows.Next() {
		town := &Town{}
		rows.Scan(&town.Name, &town.TownID)
		towns = append(towns, town)
	}
	return towns, err
}

// NewTown creates a new town struct
func NewTown(name string) *Town {
	return &Town{
		-1,
		name,
		-1,
	}
}

// Get gets town information from database
func (t *Town) Get() *Town {
	row := database.Connection.QueryRow("SELECT id, town_id FROM cloaka_towns WHERE name = ?", t.Name)
	row.Scan(&t.ID, &t.TownID)
	return t
}

// Exists checks if a town exists
func (t *Town) Exists() bool {
	row := database.Connection.QueryRow("SELECT EXISTS(SELECT 1 FROM cloaka_towns WHERE name = ?)", t.Name)
	exists := false
	row.Scan(&exists)
	return exists
}

// GetTownByName gets a town by its name
func GetTownByName(name string) *Town {
	town := NewTown(name)
	if !town.Exists() {
		return nil
	}
	return town.Get()
}