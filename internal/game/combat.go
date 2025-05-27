package game

type Command struct {
	TroopIndex int `json:"troop_index"`
	TowerIndex int `json:"tower_index"`
}

func CalculateDamage(atk, def float64, crit bool) float64 {
	dmg := atk
	if crit {
		dmg *= 1.2
	}
	dmg -= def
	if dmg < 0.0 {
		return 0.0
	}
	return dmg
}
