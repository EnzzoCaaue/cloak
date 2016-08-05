package models

import (
	"github.com/raggaer/pigo"
)

// AdminInformation holds all admin dashboard information
type AdminInformation struct {
	Accounts  int
	Players   int
	Sorcerers int
	Druids    int
	Paladins  int
	Knights   int
	Males     int
	Females   int
}

// Record holds information about a daemon
type Record struct {
	Total int
	At    int64
}

// GetAdminInformation returns dashboard information
func GetAdminInformation() *AdminInformation {
	row := pigo.Database.QueryRow("SELECT (SELECT COUNT(1) FROM accounts) as accounts, (SELECT COUNT(1) FROM players) players, (SELECT COUNT(1) FROM players WHERE vocation IN(1, 5)) as sorcerers, (SELECT COUNT(1) FROM players WHERE vocation IN (2, 6)) as druids, (SELECT COUNT(1) FROM players WHERE vocation IN (3, 7)) as paladins, (SELECT COUNT(1) FROM players WHERE vocation IN (4, 8)) as knights, (SELECT COUNT(1) FROM players WHERE sex = 1) as males, (SELECT COUNT(1) FROM players WHERE sex = 0) as females")
	info := &AdminInformation{}
	row.Scan(&info.Accounts, &info.Players, &info.Sorcerers, &info.Druids, &info.Paladins, &info.Knights, &info.Males, &info.Females)
	return info
}

// ClearOnlineLogs clears the online daemon records
func ClearOnlineLogs() error {
	_, err := pigo.Database.Exec("DELETE FROM cloaka_online_records")
	return err
}

// GetOnlineCount returns the current online number of players
func GetOnlineCount() int {
	row := pigo.Database.QueryRow("SELECT COUNT(1) FROM players_online")
	total := 0
	row.Scan(&total)
	return total
}

// AddOnlineRecord adds an online player record
func AddOnlineRecord(total int, at int64) error {
	_, err := pigo.Database.Exec("INSERT INTO cloaka_online_records (total, at) VALUES (?, ?)", total, at)
	return err
}

// GetOnlineRecords returns online players records
func GetOnlineRecords(limit int) ([]*Record, error) {
	rows, err := pigo.Database.Query("SELECT total, at FROM cloaka_online_records ORDER BY at DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	records := []*Record{}
	for rows.Next() {
		record := &Record{}
		rows.Scan(&record.Total, &record.At)
		records = append(records, record)
	}
	return records, nil
}
