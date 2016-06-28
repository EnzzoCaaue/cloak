package template

import (
	"fmt"
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/util"
	"github.com/raggaer/pigo"
	"html/template"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Load loads the AAC template
func Load() {
	pigo.LoadTemplate("[[", "]]", template.FuncMap{
		"gender": func(gender int) string {
			return util.GetGender(gender)
		},
		"vocation": func(vocation int) string {
			return util.GetVocation(vocation)
		},
		"rawHTML": func(msg string) template.HTML {
			return template.HTML(msg)
		},
		"accountSecurity": func(twofactor int, key string) bool {
			if twofactor != 0 && len(key) > 0 {
				return true
			}
			return false
		},
		"hasRecoveryKey": func(key string) bool {
			return len(key) > 0
		},
		"unixToNews": func(unix int64) string {
			t := time.Unix(unix, 0)
			return fmt.Sprintf("%v %v %v -", t.Month().String()[:3], t.Day(), t.Year())
		},
		"unixToDate": func(unix int64) string {
			if unix == 0 {
				return "Never"
			}
			timeDate := time.Unix(unix, 0)
			timeString := fmt.Sprintf("%v %v %v, %v:%v:%v", timeDate.Month().String()[:3], timeDate.Day(), timeDate.Year(), timeDate.Hour(), timeDate.Minute(), timeDate.Second())
			return timeString
		},
		"currentMenu": func(active string, key string) bool {
			if active == key {
				return true
			}
			return false
		},
		"isPremium": func(days int) bool {
			if days > 0 {
				return true
			}
			return false
		},
		"isAdmin": func(admin int) bool {
			if admin > 0 {
				return true
			}
			return false
		},
		"premiumDaysString": func(days int) string {
			daysStr := strconv.Itoa(days)
			str := "You have " + daysStr
			if days == 1 {
				str = str + " day left"
			} else {
				str = str + " days left"
			}
			return str
		},
		"urlEncode": func(name string) string {
			nameScape := url.QueryEscape(name)
			return nameScape
		},
		"isAlive": func(time int) bool {
			if time <= 0 {
				return true
			}
			return false
		},
		"maskedText": func(text string) string {
			i := len(text)
			s := ""
			for x := 0; x < i; x++ {
				s = s + "*"
			}
			return s
		},
		"isEven": func(number int) bool {
			if number%2 == 0 {
				return true
			}
			return false
		},
		"isLogged": func() bool {
			return true
		},
		"isCurrentSkill": func(current string, skill string) string {
			if current == skill {
				return "active"
			}
			return ""
		},
		"isNotCurrentCharacter": func(current string, name string) bool {
			if current == name {
				return false
			}
			return true
		},
		"getCaptchaKey": func() string {
			return pigo.Config.Key("captcha").String("public")
		},
		"parseComment": func(comment string) []string {
			return strings.Split(comment, "\n")
		},
		"showMostDamage": func(killer string, mostdamager string) bool {
			if killer != mostdamager {
				return true
			}
			return false
		},
		"longToShort": func(msg string) string {
			if len(msg) > 30 {
				return msg[:30]
			}
			return msg
		},
		"getTopPlayers": func(limit int) []*models.Player {
			if pigo.Cache.IsExpired("topPlayers") {
				players, err := models.GetTopPlayers(limit)
				if err != nil {
					return nil
				}
				pigo.Cache.Put("topPlayers", time.Minute, players)
				return players
			}
			return pigo.Cache.Get("topPlayers").([]*models.Player)
		},
		"coolIndex": func(index int) int {
			return index + 1
		},
		"bytesToMb": func(bytes uint64) float64 {
			return float64(bytes / 1000000)
		},
		"isMyPlayer": func(name string, names []string) bool {
			for _, n := range names {
				if name == n {
					return true
				}
			}
			return false
		},
	})
}
