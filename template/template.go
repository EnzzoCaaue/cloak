package template

import (
	"fmt"
	"html/template"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Cloakaac/cloak/util"
	"github.com/yaimko/yaimko"
)

// Load loads the AAC template
func Load() {
	yaimko.Template.Delims("[[", "]]")
	m := template.FuncMap{
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
			return (twofactor != 0 && len(key) > 0)
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
		"luaUnixToDate": func(unix string) string {
			b, err := strconv.ParseInt(unix, 10, 64)
			if err != nil {
				return "Never"
			}
			timeDate := time.Unix(b, 0)
			timeString := fmt.Sprintf("%v %v %v, %v:%v:%v", timeDate.Month().String()[:3], timeDate.Day(), timeDate.Year(), timeDate.Hour(), timeDate.Minute(), timeDate.Second())
			return timeString
		},
		"currentMenu": func(active string, key string) bool {
			return active == key
		},
		"isPremium": func(days int) bool {
			return days > 0
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
			return time <= 0
		},
		"maskedText": func(text string) string {
			return strings.Repeat("*", len(text))
		},
		"isEven": func(number int) bool {
			return number%2 == 0
		},
		"isCurrentSkill": func(current string, skill string) string {
			if current == skill {
				return "active"
			}
			return ""
		},
		"isNotCurrentCharacter": func(current string, name string) bool {
			return current != name
		},
		"getCaptchaKey": func() string {
			return yaimko.Config.String("captcha.public")
		},
		"parseComment": func(comment string) []string {
			return strings.Split(comment, "\n")
		},
		"showMostDamage": func(killer string, mostdamager string) bool {
			return killer != mostdamager
		},
		"longToShort": func(msg string) string {
			if len(msg) > 30 {
				return msg[:30]
			}
			return msg
		},
		"getTopPlayers": func(limit int) string {
			return ""
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
		"csrfField": func(token string) template.HTML {
			return template.HTML(fmt.Sprintf("<input type='hidden' name='_csrf' value = '%v'>", token))
		},
	}
	yaimko.Template.Funcs(m)
	if err := yaimko.LoadTemplateFiles(); err != nil {
		yaimko.ERROR.Fatal(err)
	}
}
