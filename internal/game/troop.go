package game

type Troop struct {
	Name string  `json:"name"`
	HP   float64 `json:"hp"`
	ATK  float64 `json:"atk"`
	DEF  float64 `json:"def"`
	Mana float64 `json:"mana"`
	EXP  float64 `json:"exp"`
}

func (t *Troop) Attack(target *Tower) {
}
