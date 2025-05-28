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

func (p *Player) Attack(o *Player, command *Command) (network.Message, error) {
	// CASE 0: troop out of index
	if command.TroopIndex < 0 || command.TroopIndex > 2 {
		return network.Message{Type: config.MsgError, Data: []string{"🤖 invaid index"}}, fmt.Errorf("troop out of index")
	}
	troop := &p.Troops[command.TroopIndex]

	// CASE 1: Troop destroyed
	if troop.HP <= 0 {
		return network.Message{Type: config.MsgError, Data: []string{fmt.Sprintf("🤖 %s destroyed!", troop.Name)}}, fmt.Errorf("troop destroyed")
	}

	// CASE 2: Not enough mana
	if p.Mana <= 0 || p.Mana < troop.Mana {
		return network.Message{Type: config.MsgError, Data: []string{"You don't have enough mana!"}}, fmt.Errorf("not enough mana")
	}

	// CASE 3: Wrong target tower
	switch command.TowerIndex {
	case 0:
		if !((o.Towers[0].HP > 0) && (o.Towers[0].HP <= o.User.Metadata.Towers[0].HP) && (o.Towers[1].HP == o.User.Metadata.Towers[1].HP) && (o.Towers[2].HP == o.User.Metadata.Towers[2].HP)) {
			return network.Message{Type: config.MsgError, Data: []string{"🏰 invaid index"}}, fmt.Errorf("wrong target tower")
		}
	case 1:
		if !((o.Towers[0].HP <= 0) && (o.Towers[1].HP > 0) && (o.Towers[1].HP <= o.User.Metadata.Towers[1].HP) && (o.Towers[2].HP == o.User.Metadata.Towers[2].HP)) {
			return network.Message{Type: config.MsgError, Data: []string{"🏰 invaid index"}}, fmt.Errorf("wrong target tower")
		}
	case 2:
		if !((o.Towers[0].HP <= 0) && (o.Towers[1].HP <= 0) && (o.Towers[2].HP > 0) && (o.Towers[2].HP <= o.User.Metadata.Towers[2].HP)) {
			return network.Message{Type: config.MsgError, Data: []string{"🏰 invaid index"}}, fmt.Errorf("wrong target tower")
		}
	default:
		return network.Message{Type: config.MsgError, Data: []string{"🏰 not found!"}}, fmt.Errorf("wrong target tower")
	}

	tower := &o.Towers[command.TowerIndex]

	p.Mana -= troop.Mana

	crit := tower.HasCrit()

	var dmgToTroop float64
	var dmgToTroopOrigin float64
	var dmgToTroopAddition float64
	var defenseDmgToTroop float64

	if crit {
		dmgToTroop = max(tower.ATK*1.2-troop.DEF, 0.0)
		dmgToTroopOrigin = tower.ATK
		dmgToTroopAddition = tower.ATK * 0.2
		defenseDmgToTroop = troop.DEF
	} else {
		dmgToTroop = max(tower.ATK-troop.DEF, 0.0)
		dmgToTroopOrigin = tower.ATK
		dmgToTroopAddition = 0
		defenseDmgToTroop = troop.DEF
	}
	troop.HP = max(troop.HP-dmgToTroop, 0.0)

	dmgToTower := max(troop.ATK-tower.DEF, 0.0)
	dmgToTowerOrigin := troop.ATK
	defenseDmgToTower := tower.DEF
	tower.HP = max(tower.HP-dmgToTower, 0.0)

	return network.Message{Type: config.MsgAttackResult, Data: CombatResult{
		Attacker:             p.User.Metadata.Username,
		Defender:             o.User.Metadata.Username,
		UsingTroop:           *troop,
		TargetTower:          *tower,
		DamgeToTroop:         dmgToTroop,
		DamgeToTroopOrigin:   dmgToTroopOrigin,
		DamgeToTroopAddition: dmgToTroopAddition,
		DefenseDamgeToTroop:  defenseDmgToTroop,
		DamgeToTower:         dmgToTower,
		DamgeToTowerOrigin:   dmgToTowerOrigin,
		DefenseDamgeToTower:  defenseDmgToTower,
	}}, nil
}
