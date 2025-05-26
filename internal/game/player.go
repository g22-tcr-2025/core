package game

type Player struct {
	User   User    `json:"user"`
	Mana   float64 `json:"mana"`
	EXP    float64 `json:"exp"`
	Troops []Troop `json:"troops"`
	Towers []Tower `json:"towers"`
}
