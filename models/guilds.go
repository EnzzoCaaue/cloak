package models

import (
	"github.com/raggaer/pigo"
)

type Guild struct {
	ID          int64
	Name        string
	Owner       *Player
	Creation    int64
	Motd        string
	Info        *GuildInfo
	Members     []*Player
	Invitations []*Player
	Ranks       []*GuildRank
}

type GuildRank struct {
	Name  string
	Level int
}

type GuildInfo struct {
	Online  int
	Members int
	Top     int
	Low     int
}

// NewGuild returns a new guild struct
func NewGuild() *Guild {
	guild := &Guild{
		-1,
		"",
		&Player{},
		0,
		"",
		&GuildInfo{},
		[]*Player{},
		[]*Player{},
		[]*GuildRank{},
	}
	return guild
}

// GetGuildList gets the full list from database
func GetGuildList() ([]*Guild, error) {
	list := []*Guild{}
	rows, err := pigo.Database.Query("SELECT a.name, a.creationdata, a.motd, b.name, (SELECT COUNT(*) FROM guild_membership WHERE a.id = guild_id) AS members, (SELECT COUNT(*) FROM guild_membership c, players_online d WHERE c.player_id = d.player_id AND c.guild_id = a.id) AS onl, (SELECT MAX(f.level) FROM guild_membership e, players f WHERE f.id = e.player_id AND e.guild_id = a.id) as top, (SELECT MIN(f.level) FROM guild_membership e, players f WHERE f.id = e.player_id AND e.guild_id = a.id) as low FROM guilds a, players b WHERE a.ownerid = b.id ORDER BY a.creationdata DESC")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		guild := NewGuild()
		rows.Scan(&guild.Name, &guild.Creation, &guild.Motd, &guild.Owner.Name, &guild.Info.Members, &guild.Info.Online, &guild.Info.Top, &guild.Info.Low)
		list = append(list, guild)
	}
	return list, nil
}

// GuildExists checks if a guild exists
func GuildExists(name string) bool {
	row := pigo.Database.QueryRow("SELECT EXISTS(SELECT 1 FROM guilds WHERE name = ?)", name)
	exists := false
	row.Scan(&exists)
	return exists
}

// Create creates a guild
func (guild *Guild) Create() error {
	r, err := pigo.Database.Exec("INSERT INTO guilds (name, ownerid, creationdata, motd) VALUES(?, ?, ?, ?)", guild.Name, guild.Owner.ID, guild.Creation, guild.Motd)
	if err != nil {
		return err
	}
	guild.ID, _ = r.LastInsertId()
	_, err = pigo.Database.Exec("INSERT INTO guild_membership (player_id, guild_id, rank_id) VALUES (?, ?, (SELECT id FROM guild_ranks WHERE guild_id = ? AND level = 3))", guild.Owner.ID, guild.ID, guild.ID)
	return err
}

// GetGuildByName gets a guild by its name
func GetGuildByName(name string) (*Guild, error) {
	row := pigo.Database.QueryRow("SELECT a.ownerid, a.id, a.name, a.creationdata, a.motd, b.name, (SELECT COUNT(*) FROM guild_membership WHERE a.id = guild_id) AS members, (SELECT COUNT(*) FROM guild_membership c, players_online d WHERE c.player_id = d.player_id AND c.guild_id = a.id) AS onl, (SELECT MAX(f.level) FROM guild_membership e, players f WHERE f.id = e.player_id AND e.guild_id = a.id) as top, (SELECT MIN(f.level) FROM guild_membership e, players f WHERE f.id = e.player_id AND e.guild_id = a.id) as low FROM guilds a, players b WHERE a.ownerid = b.id AND a.name = ?", name)
	guild := NewGuild()
	row.Scan(&guild.Owner.ID, &guild.ID, &guild.Name, &guild.Creation, &guild.Motd, &guild.Owner.Name, &guild.Info.Members, &guild.Info.Online, &guild.Info.Top, &guild.Info.Low)
	rows, err := pigo.Database.Query("SELECT p.id, p.name, p.level, p.vocation, gm.nick, gr.name AS rank_name, IF(po.player_id IS NULL, 0, 1) as onl FROM players AS p LEFT JOIN guild_membership AS gm ON gm.player_id = p.id LEFT JOIN guild_ranks AS gr ON gr.id = gm.rank_id LEFT JOIN players_online AS po ON p.id = po.player_id WHERE gm.guild_id = ? ORDER BY gm.rank_id, p.name", guild.ID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		player := NewPlayer()
		rows.Scan(&player.ID, &player.Name, &player.Level, &player.Vocation, &player.GuildNick, &player.GuildRank, &player.Online)
		guild.Members = append(guild.Members, player)
	}
	rows.Close()
	rows, err = pigo.Database.Query("SELECT a.name FROM players a, guild_invites b WHERE a.id = b.player_id AND b.guild_id = ?", guild.ID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		player := NewPlayer()
		rows.Scan(&player.Name)
		guild.Invitations = append(guild.Invitations, player)
	}
	rows.Close()
	rows, err = pigo.Database.Query("SELECT name, level FROM guild_ranks WHERE guild_id = ? ORDER BY level DESC", guild.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		rank := &GuildRank{}
		rows.Scan(&rank.Name, &rank.Level)
		guild.Ranks = append(guild.Ranks, rank)
	}
	return guild, nil
}

// ChangeMotd changes a guild motd
func (guild *Guild) ChangeMotd(motd string) error {
	_, err := pigo.Database.Exec("UPDATE guilds SET motd = ? WHERE id = ?", motd, guild.ID)
	return err
}

// ChangeRanks changes a guild ranks
func (guild *Guild) ChangeRanks(third, second, first string) error {
	_, err := pigo.Database.Exec(`
	UPDATE guild_ranks SET name =
	( CASE WHEN level = 3 THEN ?
		   WHEN level = 2 THEN ?
		   WHEN level = 1 THEN ?
	 END )
	 WHERE guild_id = ?`, third, second, first, guild.ID)
	return err
}

// InvitePlayer invites a player to the guilds
func (guild *Guild) InvitePlayer(player int64) error {
	_, err := pigo.Database.Exec("INSERT INTO guild_invites (player_id, guild_id) VALUES (?, ?)", player, guild.ID)
	return err
}

func (guild *Guild) IsInvited(player int64) bool {
	row := pigo.Database.QueryRow("SELECT EXISTS(SELECT 1 FROM guild_invites WHERE player_id = ? AND guild_id = ?)", player, guild.ID)
	exists := false
	row.Scan(&exists)
	return exists
}
