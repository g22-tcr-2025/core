package game

import (
	"clash-royale/internal/config"
	"clash-royale/internal/network"
	"fmt"
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
	// CASE 0: troop out of index
	if command.TroopIndex < 0 || command.TroopIndex > 2 {
		return network.Message{Type: config.MsgError, Data: "🤖 invaid index"}
	}
	troop := &p.Troops[command.TroopIndex]

	// CASE 1: Not enough mana
	if p.Mana <= 0 || p.Mana < troop.Mana {
		return network.Message{Type: config.MsgError, Data: "You don't have enough mana!"}
	}

	// CASE 2: Troop destroyed
	if troop.HP <= 0 {
		return network.Message{Type: config.MsgError, Data: fmt.Sprintf("🤖 %s destroyed!", troop.Name)}
	}

	// CASE 3: Wrong target tower
	switch command.TowerIndex {
	case 0:
		if !((o.Towers[0].HP > 0) && (o.Towers[0].HP <= o.User.Metadata.Towers[0].HP) && (o.Towers[1].HP == o.User.Metadata.Towers[1].HP) && (o.Towers[2].HP == o.User.Metadata.Towers[2].HP)) {
			return network.Message{Type: config.MsgError, Data: "🏰 invaid index"}
		}
	case 1:
		if !((o.Towers[0].HP <= 0) && (o.Towers[1].HP > 0) && (o.Towers[1].HP <= o.User.Metadata.Towers[1].HP) && (o.Towers[2].HP == o.User.Metadata.Towers[2].HP)) {
			return network.Message{Type: config.MsgError, Data: "🏰 invaid index"}
		}
	case 2:
		if !((o.Towers[0].HP <= 0) && (o.Towers[1].HP <= 0) && (o.Towers[2].HP > 0) && (o.Towers[2].HP <= o.User.Metadata.Towers[2].HP)) {
			return network.Message{Type: config.MsgError, Data: "🏰 invaid index"}
		}
	default:
		return network.Message{Type: config.MsgError, Data: "🏰 not found!"}
	}

	tower := &o.Towers[command.TowerIndex]

	p.Mana -= troop.Mana

	crit := tower.HasCrit()
	// fmt.Println(crit)

	var dmgToTroop float64
	if crit {
		dmgToTroop = max(tower.ATK*1.2-troop.DEF, 0.0)
	} else {
		dmgToTroop = max(tower.ATK-troop.DEF, 0.0)
	}
	troop.HP = max(troop.HP-dmgToTroop, 0.0)

	dmgToTower := max(troop.ATK-tower.DEF, 0.0)
	tower.HP = max(tower.HP-dmgToTower, 0.0)

	return network.Message{Type: config.MsgAttackResult, Data: CombatResult{
		Attacker:     p.User.Metadata.Username,
		Defender:     o.User.Metadata.Username,
		UsingTroop:   *troop,
		TargetTower:  *tower,
		DamgeToTroop: dmgToTroop,
		DamgeToTower: dmgToTower,
	}}
}
