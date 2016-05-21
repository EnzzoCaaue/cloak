package models

import (
    "github.com/Cloakaac/cloak/database"
)

type Guild struct {
    ID int64
    Name string
    Owner *Player
    Creation int64
    Motd string
}

// NewGuild returns a new guild struct
func NewGuild() *Guild {
    guild := &Guild{
        -1,
        "",
        &Player{},
        0,
        "",
    }
    return guild
}

// GetGuildList gets the full list from database
func GetGuildList() ([]*Guild, error) {
    list := []*Guild{}
    rows, err := database.Connection.Query("SELECT a.name, a.creationdata, a.motd, b.name FROM guilds a, players b WHERE a.ownerid = b.id ORDER BY a.creationdata DESC")
    defer rows.Close()
    if err != nil {
        return nil, err
    }
    for rows.Next() {
        guild := NewGuild()
        rows.Scan(&guild.Name, &guild.Creation, &guild.Motd, &guild.Owner.Name)
        list = append(list, guild)
    }
    return list, nil
}

// GuildExists checks if a guild exists
func GuildExists(name string) bool {
    row := database.Connection.QueryRow("SELECT EXISTS(SELECT 1 FROM guilds WHERE name = ?)", name)
    exists := false
    row.Scan(&exists)
    return exists
}

// Create creates a guild
func (guild *Guild) Create() error {
    _, err := database.Connection.Exec("INSERT INTO guilds (name, ownerid, creationdata, motd) VALUES(?, ?, ?, ?)", guild.Name, guild.Owner.ID, guild.Creation, guild.Motd)
    return err
}