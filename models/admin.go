package models

import (
	"github.com/raggaer/pigo"
)

// AdminInformation holds all admin dashboard information
type AdminInformation struct {
    Accounts int
    Players int
    Sorcerers int
    Druids int
    Paladins int
    Knights int
    Males int
    Females int
}

// GetAdminInformation returns dashboard information
func GetAdminInformation() *AdminInformation {
    row := pigo.Database.QueryRow("SELECT (SELECT COUNT(*) FROM accounts) as accounts, (SELECT COUNT(*) FROM players) players, (SELECT COUNT(*) FROM players WHERE vocation IN(1, 5)) as sorcerers, (SELECT COUNT(*) FROM players WHERE vocation IN (2, 6)) as druids, (sELECT COUNT(*) FROM players WHERE vocation IN (3, 7)) as paladins, (SELECT COUNT(*) FROM players WHERE vocation IN (4, 8)) as knights, (SELECT COUNT(*) FROM players WHERE sex = 1) as males, (SELECT COUNT(*) FROM players WHERE sex = 0) as females")
    info := &AdminInformation{}
    row.Scan(&info.Accounts, &info.Players, &info.Sorcerers, &info.Druids, &info.Paladins, &info.Knights, &info.Males, &info.Females)
    return info
}