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