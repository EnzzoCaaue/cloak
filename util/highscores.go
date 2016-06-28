package util

import (
    "fmt"
)

// GetHighscoreQuery returns a highscore query
func GetHighscoreQuery(page int, highscoreType string, per int) (int, string, string) {
    query := "SELECT name, level FROM players ORDER BY level DESC LIMIT 0,10"
    skillName := "Experience"
    pageIndex := page * per
    switch highscoreType {
    case "level":
        skillName = "Experience"
        query = fmt.Sprintf("SELECT name, level FROM players ORDER BY level DESC LIMIT %v, %v", 
            pageIndex,
            per,
        )
    case "magic":
        skillName = "Magic Level"
        query = fmt.Sprintf("SELECT name, maglevel FROM players ORDER BY maglevel DESC LIMIT %v, %v",
            pageIndex,
            per,
        )
    case "shield":
        skillName = "Shielding"
        query = fmt.Sprintf("SELECT name, skill_shielding FROM players ORDER BY skill_shielding DESC LIMIT %v, %v",
            pageIndex,
            per,
        ) 
    case "distance":
        skillName = "Distance Fighting"
        query = fmt.Sprintf("SELECT name, skill_dist FROM players ORDER BY skill_dist DESC LIMIT %v, %v",
            pageIndex,
            per,
        )     
    case "sword":
        skillName = "Sword Fighting"
        query = fmt.Sprintf("SELECT name, skill_sword FROM players ORDER BY skill_sword DESC LIMIT %v, %v",
            pageIndex,
            per,
        )   
    case "axe":
        skillName = "Axe Fighting"
        query = fmt.Sprintf("SELECT name, skill_axe FROM players ORDER BY skill_axe DESC LIMIT %v, %v",
            pageIndex,
            per,
        )     
    case "fist":
        skillName = "Fist Fighting"
        query = fmt.Sprintf("SELECT name, skill_fist FROM players ORDER BY skill_fist DESC LIMIT %v, %v",
            pageIndex,
            per,
        )       
    case "club":
        skillName = "Club Fighting"
        query = fmt.Sprintf("SELECT name, skill_club FROM players ORDER BY skill_club DESC LIMIT %v, %v",
            pageIndex,
            per,
        )
    case "fish":
        skillName = "Fishing"
        query = fmt.Sprintf("SELECT name, skill_fishing FROM players ORDER BY skill_fishing DESC LIMIT %v, %v",
            pageIndex,
            per,
        )      
    }
    return pageIndex, query, skillName
}