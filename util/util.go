package util

const (
	novoc = iota
	sorcerer
	druid
	paladin
	knight
	masterSorcerer
	elderDruid
	royalPaladin
	eliteKnight
)

const (
	female = iota
	male
)

var (
	// Mode stores the AAC mode
	Mode int
	genderList = map[string]int{
		"Male":   male,
		"Female": female,
	}
	vocationList = map[string]int{
		"Sorcerer":        sorcerer,
		"Druid":           druid,
		"Paladin":         paladin,
		"Knight":          knight,
		"Master Sorcerer": masterSorcerer,
		"Elder Druid":     elderDruid,
		"Royal Paladin":   royalPaladin,
		"Elite Knight":    eliteKnight,
	}
)

// SetMode sets the AAC run mode DEBUG(0) RELEASE(1)
func SetMode(mode int) {
	Mode = mode
}

// Vocation gets the vocation id from a given string
func Vocation(voc string) int {
	return vocationList[voc]
}

// Gender gets the gender id from a given string
func Gender(gender string) int {
	return genderList[gender]
}