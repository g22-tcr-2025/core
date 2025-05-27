package game

import (
	"clash-royale/internal/config"
	"clash-royale/internal/network"
	"sync"
)

type Player struct {
	User   User    `json:"user"`
	Mana   float64 `json:"mana"`
	EXP    float64 `json:"exp"`
	Troops []Troop `json:"troops"`
	Towers []Tower `json:"towers"`
	Mutex  sync.Mutex
}

func (p *Player) Attack(o *Player, command *Command) network.Message {
	troop := p.Troops[command.TroopIndex]
	if p.Mana <= 0 || p.Mana < troop.Mana {
		return network.Message{Type: config.MsgError, Data: "You don't have enough mana!"}
	}
	p.Mana -= troop.Mana

	tower := o.Towers[command.TowerIndex]

	crit := tower.HasCrit()
	// fmt.Println(crit)

	var dmgToTroop float64
	if crit {
		dmgToTroop = max(tower.ATK*1.2-troop.DEF, 0.0)
	} else {
		dmgToTroop = max(tower.ATK-troop.DEF, 0.0)
	}
	troop.HP -= dmgToTroop

	dmgToTower := max(troop.ATK-tower.DEF, 0.0)
	tower.HP -= dmgToTower
	return network.Message{Type: config.MsgAttackResult, Data: CombatResult{
		Attacker:     p.User.Metadata.Username,
		Defender:     o.User.Metadata.Username,
		UsingTroop:   troop,
		TargetTower:  tower,
		DamgeToTroop: dmgToTroop,
		DamgeToTower: dmgToTower,
	}}
}
