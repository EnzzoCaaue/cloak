package models

import (
	"github.com/Cloakaac/cloak/database"
)

// Player struct for database players
type Player struct {
	ID          int64
	AccountID   int64
	Name        string
	Vocation    int
	Gender         int
	Level       int
	Health int
	HealthMax int
	Mana int
	ManaMax int
	LookBody int
	LookFeet int
	LookHead int
	LookLegs int
	LookType int
	LookAddons int
	MagicLevel int
	Soul int
	Town *Town
	Stamina int	
	SkillFist int
	SkillClub int
	SkillSword int
	SkillAxe int
	SkillDist int
	SkillShield int
	SkillFish int
	Experience int
	Balance int
	Premdays int
	LastLogin int64
	GuildName string
	GuildRank string
	Cloaka *CloakaPlayer
}

// CloakaPlayer struct for cloaka_players
type CloakaPlayer struct {
	ID int64
	PlayerID int64
	Comment string
	Deleted int
	Hide int
}

// GetTopPlayers gets sidebar top players
func GetTopPlayers(limit int) ([]*Player, error) {
	rows, err := database.Connection.Query("SELECT name, level FROM players ORDER BY level DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	players := []*Player{}
	for rows.Next() {
		player := &Player{}
		rows.Scan(&player.Name, &player.Level)
		players = append(players, player)
	}
	return players, nil
}

// NewPlayer returns a new player struct
func NewPlayer() *Player {
	player := &Player{}
	player.Cloaka = &CloakaPlayer{}
	player.Town = &Town{}
	return player
}

// GetPlayerByName gets a character by its name
func GetPlayerByName(name string) *Player {
	player := NewPlayer()
	player.Name = name
	if !player.Exists() {
		return nil
	}
	row := database.Connection.QueryRow("SELECT a.id, a.name, a.level, a.vocation, a.sex, a.lastlogin, b.premdays, c.name, d.name, e.name FROM players a, accounts b, cloaka_towns c, guilds d, guild_ranks e, guild_membership f WHERE d.id = f.guild_id AND f.player_id = a.id AND e.id = f.rank_id AND a.account_id = b.id AND a.town_id = c.id AND a.name = ?", player.Name)
	row.Scan(&player.ID, &player.Name, &player.Level, &player.Vocation, &player.Gender, &player.LastLogin, &player.Premdays, &player.Town.Name, &player.GuildName, &player.GuildRank)
	return player
}

// Save saves a player into a database
func (player *Player) Save() error {
	result, err := database.Connection.Exec(`INSERT INTO
	players
	(name, account_id, level, vocation, health, healthmax, experience, lookbody, lookfeet, lookhead, looklegs, looktype, lookaddons, maglevel, mana, manamax, soul, town_id, posx, posy, posz, conditions, cap, sex, stamina, skill_fist, skill_club, skill_sword, skill_axe, skill_dist, skill_shielding, skill_fishing)
	VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		player.Name,
		player.AccountID,
		player.Level,
		player.Vocation,
		player.Health,
		player.HealthMax,
		player.Experience,
		player.LookBody,
		player.LookFeet,
		player.LookHead,
		player.LookLegs,
		player.LookType,
		player.LookAddons,
		player.MagicLevel,
		player.Mana,
		player.ManaMax,
		100,
		player.Town.ID,
		0,
		0,
		0,
		"",
		100,
		player.Gender,
		player.Stamina,
		player.SkillFist,
		player.SkillClub,
		player.SkillSword,
		player.SkillAxe,
		player.SkillDist,
		player.SkillShield,
		player.SkillFish,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	player.ID = id
	_, err = database.Connection.Exec("INSERT INTO cloaka_players (player_id) VALUES (?)", player.ID)
	return err
}

// Exists checks if a character name is already in use
func (player *Player) Exists() bool {
	row := database.Connection.QueryRow("SELECT EXISTS(SELECT 1 FROM players WHERE name = ?)", player.Name)
	exists := false
	row.Scan(&exists)
	return exists
}

// GetCharacters gets all account characters
func (account *CloakaAccount) GetCharacters() ([]*Player, error) {
	rows, err := database.Connection.Query("SELECT a.id, a.name, a.vocation, a.level, a.lastlogin, a.balance, a.sex, b.deleted, c.name FROM players a, cloaka_players b, cloaka_towns c WHERE a.account_id = ? AND b.player_id = a.id AND a.town_id = c.town_id ORDER BY a.id DESC", account.Account.ID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	characters := []*Player{}
	for rows.Next() {
		player := NewPlayer()
		rows.Scan(&player.ID, &player.Name, &player.Vocation, &player.Level, &player.LastLogin, &player.Balance, &player.Gender, &player.Cloaka.Deleted, &player.Town.Name)
		characters = append(characters, player)
	}
	return characters, nil
}