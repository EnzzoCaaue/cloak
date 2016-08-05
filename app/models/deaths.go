package models

// Death is the struct for character deaths
type Death struct {
	Time                  int64
	Level                 int
	KilledBy              string
	IsPlayer              int
	MostDamageBy          string
	MostDamageIsPlayer    int
	Unjustified           int
	MostDamageUnjustified int
}
