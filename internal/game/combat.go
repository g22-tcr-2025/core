package game

type Command struct {
	TroopIndex int `json:"troop_index"`
	TowerIndex int `json:"tower_index"`
}

type CombatResult struct {
	Attacker             string  `json:"attacker"`
	Defender             string  `json:"defender"`
	UsingTroop           Troop   `json:"using_troop"`
	TargetTower          Tower   `json:"target_tower"`
	DamgeToTroop         float64 `json:"dmg_to_troop"`
	DamgeToTroopOrigin   float64 `json:"dmg_to_troop_origin"`
	DamgeToTroopAddition float64 `json:"dmg_to_troop_addition"`
	DefenseDamgeToTroop  float64 `json:"defense_dmg_to_troop"`
	DamgeToTower         float64 `json:"dmg_to_tower"`
	DamgeToTowerOrigin   float64 `json:"dmg_to_tower_origin"`
	DefenseDamgeToTower  float64 `json:"defense_dmg_to_tower"`
}
