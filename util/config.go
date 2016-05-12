package util

// Parser stores the database config values
var Parser = &Config{}

// Config is the struct that stores all the database config values
type Config struct {
	General struct {
		Path       string
		Port       string
		Secret     string
		Logosize   int
		Daemons    bool
		DaemonsInt int
	}
	Mysql struct {
		User     string
		Password string
		Database string
	}
	Captcha struct {
		Key    string
		Secret string
	}
	Highscores struct {
		Per int
	}
	Character struct {
		Max int
	}
	Register struct {
		Premdays   int
		Level      int
		Health     int
		Healthmax  int
		Mana       int
		Manamax    int
		Experience int
		Cap        int
		Maglevel   int
		Stamina    int
	}
	Skills struct {
		Fist   int
		Club   int
		Sword  int
		Axe    int
		Dist   int
		Shield int
		Fish   int
	}
	Spawn struct {
		Name string
		Town int
		Posx int
		Posy int
		Posz int
	}
	Malelooktype struct {
		Looktype   int
		Lookhead   int
		Lookbody   int
		Looklegs   int
		Lookfeet   int
		Lookaddons int
	}
	Femalelooktype struct {
		Looktype   int
		Lookhead   int
		Lookbody   int
		Looklegs   int
		Lookfeet   int
		Lookaddons int
	}
	Style struct {
		Template string
	}
	Paypal struct {
		Public      string
		Private     string
		Type        string
		Currency    string
		Points      int
		Min         int
		Max         int
		Promo       int
		Description string
	}
	Paygol struct {
		Serviceid   int
		Type        string
		Currency    string
		Points      int
		Min         int
		Max         int
		Promo       int
		Description string
	}
}

// SetTemplate sets the AAC template path
func SetTemplate(path string) {
	Parser.Style.Template = path
}