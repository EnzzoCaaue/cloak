package models

import (
	"github.com/raggaer/pigo"
)

// Player struct for database players
type Player struct {
	ID          int64
	AccountID   int64
	Name        string
	Vocation    int
	Gender      int
	Level       int
	Health      int
	HealthMax   int
	Mana        int
	ManaMax     int
	LookBody    int
	LookFeet    int
	LookHead    int
	LookLegs    int
	LookType    int
	LookAddons  int
	MagicLevel  int
	Soul        int
	Town        *Town
	Stamina     int
	SkillFist   int
	SkillClub   int
	SkillSword  int
	SkillAxe    int
	SkillDist   int
	SkillShield int
	SkillFish   int
	Experience  int
	Balance     int
	Premdays    int
	LastLogin   int64
	GuildName   string
	GuildRank   string
	GuildNick   string
	Online      int
	Cloaka      *CloakaPlayer
}

// CloakaPlayer struct for cloaka_players
type CloakaPlayer struct {
	ID       int64
	PlayerID int64
	Comment  string
	Deleted  int
	Hide     int
}

// HighscorePlayer contains players for highscores page
type HighscorePlayer struct {
	Name string
	Value int
	Place int
}

// GetTopPlayers gets sidebar top players
func GetTopPlayers(limit int) ([]*Player, error) {
	rows, err := pigo.Database.Query("SELECT name, level FROM players ORDER BY level DESC LIMIT ?", limit)
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
	row := pigo.Database.QueryRow("SELECT a.looktype, a.lookbody, a.lookhead, a.looklegs, a.lookfeet, a.lookaddons, a.id, a.name, a.level, a.vocation, a.sex, a.lastlogin, b.premdays, c.name, g.deleted FROM players a, accounts b, cloaka_towns c, cloaka_players g WHERE a.account_id = b.id AND a.town_id = c.id AND g.player_id = a.id AND a.name = ?", player.Name)
	row.Scan(&player.LookType, &player.LookBody, &player.LookHead, &player.LookLegs, &player.LookFeet, &player.LookAddons, &player.ID, &player.Name, &player.Level, &player.Vocation, &player.Gender, &player.LastLogin, &player.Premdays, &player.Town.Name, &player.Cloaka.Deleted)
	return player
}

// GetGuild gets a character guild
func (player *Player) GetGuild() {
	row := pigo.Database.QueryRow("SELECT a.name, b.name FROM guilds a, guild_ranks b, guild_membership c WHERE a.id = c.guild_id AND c.player_id = ? AND b.id = c.rank_id", player.ID)
	row.Scan(&player.GuildName, &player.GuildRank)
}

// Save saves a player into a database
func (player *Player) Save() error {
	result, err := pigo.Database.Exec(`INSERT INTO
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
	_, err = pigo.Database.Exec("INSERT INTO cloaka_players (player_id) VALUES (?)", player.ID)
	return err
}

// Exists checks if a character name is already in use
func (player *Player) Exists() bool {
	row := pigo.Database.QueryRow("SELECT EXISTS(SELECT 1 FROM players WHERE name = ?)", player.Name)
	exists := false
	row.Scan(&exists)
	return exists
}

// GetCharacters gets all account characters
func (account *CloakaAccount) GetCharacters() ([]*Player, error) {
	rows, err := pigo.Database.Query("SELECT a.id, a.name, a.vocation, a.level, a.lastlogin, a.balance, a.sex, b.deleted, c.name FROM players a, cloaka_players b, cloaka_towns c WHERE a.account_id = ? AND b.player_id = a.id AND a.town_id = c.town_id ORDER BY a.id DESC", account.Account.ID)
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

// Delete deletes a character
func (player *Player) Delete(del int64) error {
	_, err := pigo.Database.Exec("UPDATE players a, cloaka_players b SET a.deletion = ?, b.deleted = 1 WHERE b.player_id = a.id AND a.id = ?", del, player.ID)
	return err
}

// GetDeaths returns a slice with a character deaths
func (player *Player) GetDeaths() ([]*Death, error) {
	rows, err := pigo.Database.Query("SELECT a.time, a.level, a.killed_by, a.is_player, a.mostdamage_by, a.mostdamage_is_player, a.unjustified, a.mostdamage_unjustified FROM player_deaths a, players b WHERE a.player_id = b.id AND b.name = ?", player.Name)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	deaths := []*Death{}
	for rows.Next() {
		death := &Death{}
		rows.Scan(&death.Time, &death.Level, &death.KilledBy, &death.IsPlayer, &death.MostDamageBy, &death.MostDamageIsPlayer, &death.Unjustified, &death.MostDamageUnjustified)
		deaths = append(deaths, death)
	}
	return deaths, nil
}

// SearchPlayers searchs for player with name LIKE
func SearchPlayers(name string) ([]*Player, error) {
	rows, err := pigo.Database.Query("SELECT name FROM players WHERE name LIKE ?", "%"+name+"%")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	players := []*Player{}
	for rows.Next() {
		player := NewPlayer()
		rows.Scan(&player.Name)
		players = append(players, player)
	}
	return players, nil
}

// IsInGuild checks if a player is in a guild
func (player *Player) IsInGuild() bool {
	row := pigo.Database.QueryRow("SELECT EXISTS(SELECT 1 FROM guild_membership WHERE player_id = ?)", player.ID)
	exists := false
	row.Scan(&exists)
	return exists
}

// GetHighscores retuns highscores list by its type
func GetHighscores(index int, query string) ([]*HighscorePlayer, error) {
	rows, err := pigo.Database.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	highscores := []*HighscorePlayer{}
	for rows.Next() {
		h := &HighscorePlayer{}
		rows.Scan(&h.Name, &h.Value)
		h.Place = index + 1
		highscores = append(highscores, h)
		index++
	}
	return highscores, nil
}