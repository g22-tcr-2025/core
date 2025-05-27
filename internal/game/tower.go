package game

import "math/rand"

type Tower struct {
	Type string  `json:"type"`
	HP   float64 `json:"hp"`
	ATK  float64 `json:"atk"`
	DEF  float64 `json:"def"`
	Crit float64 `json:"crit"`
	EXP  float64 `json:"exp"`
}

func (t *Tower) HasCrit() bool {
	return rand.Float64() < t.Crit
}
