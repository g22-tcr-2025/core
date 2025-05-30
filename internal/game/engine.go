package game

import (
	"clash-royale/internal/config"
	"clash-royale/internal/network"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

type Engine struct {
	Players               []*Player
	Tick                  int            // Timer
	Interrupt             chan *User     // User disconnected => Remain user
	TimerEnd              chan bool      // Timer end counting
	KingOrTroopsDestroyed chan []*Player // King Destoyed => Winner
	OnRequeue             func(u *User)  // Move user back to the queue
}

func NewEngine(u1, u2 *User, onRequeue func(u *User)) *Engine {
	p1 := Player{
		User:   u1,
		Level:  1.0,
		Mana:   5.0,
		EXP:    0.0,
		Heal:   false,
		Troops: randomTroop(u1.Metadata.Troops),
		Towers: copyTowers(u1.Metadata.Towers),
		Mutex:  sync.Mutex{},
	}

	p2 := Player{
		User:   u2,
		Level:  1.0,
		Mana:   5.0,
		EXP:    0.0,
		Heal:   false,
		Troops: randomTroop(u2.Metadata.Troops),
		Towers: copyTowers(u2.Metadata.Towers),
		Mutex:  sync.Mutex{},
	}

	return &Engine{
		Players:               []*Player{&p1, &p2},
		Tick:                  0,
		Interrupt:             make(chan *User),
		TimerEnd:              make(chan bool),
		KingOrTroopsDestroyed: make(chan []*Player),
		OnRequeue:             onRequeue,
	}
}

func randomTroop(troops []*Troop) []Troop {
	tempTroops := []Troop{}
	for _, t := range troops {
		if t.Name != "Queen" {
			tempTroops = append(tempTroops, *t)
		}
	}

	rand.Shuffle(len(tempTroops), func(i, j int) {
		tempTroops[i], tempTroops[j] = tempTroops[j], tempTroops[i]
	})

	result := make([]Troop, 3)
	for i := range 3 {
		result[i] = tempTroops[i]
	}
	return result
}

func copyTowers(towers []*Tower) []Tower {
	rs := make([]Tower, len(towers))
	for i := range len(towers) {
		rs[i] = *towers[i]
	}
	return rs
}

func (e *Engine) Start() {
	p1 := MatchData{
		PUsername:      e.Players[0].User.Metadata.Username,
		PLevelMetadata: e.Players[0].User.Metadata.Level,
		PEXPMetadata:   e.Players[0].User.Metadata.EXP,
		PLevel:         e.Players[0].Level,
		PEXP:           e.Players[0].EXP,
		PMana:          e.Players[0].Mana,
		PTroops:        e.Players[0].Troops,
		PTowers:        e.Players[0].Towers,
		OUsername:      e.Players[1].User.Metadata.Username,
		OLevelMetadata: e.Players[1].User.Metadata.Level,
		OEXPMetadata:   e.Players[1].User.Metadata.EXP,
		OLevel:         e.Players[1].Level,
		OEXP:           e.Players[1].EXP,
		OMana:          e.Players[1].Mana,
		OTroops:        e.Players[1].Troops,
		OTowers:        e.Players[1].Towers,
	}
	p2 := MatchData{
		PUsername:      e.Players[1].User.Metadata.Username,
		PLevelMetadata: e.Players[1].User.Metadata.Level,
		PEXPMetadata:   e.Players[1].User.Metadata.EXP,
		PLevel:         e.Players[1].Level,
		PEXP:           e.Players[1].EXP,
		PMana:          e.Players[1].Mana,
		PTroops:        e.Players[1].Troops,
		PTowers:        e.Players[1].Towers,
		OUsername:      e.Players[0].User.Metadata.Username,
		OLevelMetadata: e.Players[0].User.Metadata.Level,
		OEXPMetadata:   e.Players[0].User.Metadata.EXP,
		OLevel:         e.Players[0].Level,
		OEXP:           e.Players[0].EXP,
		OMana:          e.Players[0].Mana,
		OTroops:        e.Players[0].Troops,
		OTowers:        e.Players[0].Towers,
	}

	network.SendMessage(e.Players[0].User.Conn, network.Message{Type: config.MsgMatchStart, Data: p1})
	network.SendMessage(e.Players[1].User.Conn, network.Message{Type: config.MsgMatchStart, Data: p2})

	go runtime(e)
	go handleCommand(e)
	go handleGameEnd(e)
}

func runtime(e *Engine) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			e.Tick++

			// Regen mana
			for _, p := range e.Players {
				if p.Mana < 10 {
					p.Mutex.Lock()
					p.Mana++
					p.Mutex.Unlock()
				}
			}

			network.SendMessage(e.Players[0].User.Conn, network.Message{Type: config.MsgTick, Data: e.Tick})
			network.SendMessage(e.Players[0].User.Conn, network.Message{Type: config.MsgUpdatePlayerMnana, Data: e.Players[0].Mana})
			network.SendMessage(e.Players[0].User.Conn, network.Message{Type: config.MsgUpdateOpponentMana, Data: e.Players[1].Mana})

			// Opponent
			network.SendMessage(e.Players[1].User.Conn, network.Message{Type: config.MsgTick, Data: e.Tick})
			network.SendMessage(e.Players[1].User.Conn, network.Message{Type: config.MsgUpdatePlayerMnana, Data: e.Players[1].Mana})
			network.SendMessage(e.Players[1].User.Conn, network.Message{Type: config.MsgUpdateOpponentMana, Data: e.Players[0].Mana})

			// Check duration
			if e.Tick >= int(config.MatchDuration.Seconds()) {
				fmt.Println("ended timer")

				e.TimerEnd <- true
				return
			}
		case <-e.TimerEnd:
			fmt.Println("ended timer by interrupt")
			return
		}
	}
}

func handleCommand(e *Engine) {
	for {
		select {
		case msg := <-e.Players[0].User.Talk:
			var command Command
			json.Unmarshal(msg.Data.(json.RawMessage), &command)

			var rs network.Message
			var err error
			if command.TowerIndex == 5 && command.TroopIndex == 5 {
				// Healing
				rs, err = e.Players[0].Healing()
			} else {
				// Attack
				rs, err = e.Players[0].Attack(e.Players[1], &command)
			}

			if err != nil {
				network.SendMessage(e.Players[0].User.Conn, rs)
			} else {
				for i, p := range e.Players {
					var playerIndex int
					var opponentIndex int
					if i == 0 {
						playerIndex = 0
						opponentIndex = 1
					} else {
						playerIndex = 1
						opponentIndex = 0
					}

					matchData := MatchData{
						PUsername:      e.Players[playerIndex].User.Metadata.Username,
						PLevelMetadata: e.Players[playerIndex].User.Metadata.Level,
						PEXPMetadata:   e.Players[playerIndex].User.Metadata.EXP,
						PLevel:         e.Players[playerIndex].Level,
						PEXP:           e.Players[playerIndex].EXP,
						PMana:          e.Players[playerIndex].Mana,
						PTroops:        e.Players[playerIndex].Troops,
						PTowers:        e.Players[playerIndex].Towers,
						OUsername:      e.Players[opponentIndex].User.Metadata.Username,
						OLevelMetadata: e.Players[opponentIndex].User.Metadata.Level,
						OEXPMetadata:   e.Players[opponentIndex].User.Metadata.EXP,
						OLevel:         e.Players[opponentIndex].Level,
						OEXP:           e.Players[opponentIndex].EXP,
						OMana:          e.Players[opponentIndex].Mana,
						OTroops:        e.Players[opponentIndex].Troops,
						OTowers:        e.Players[opponentIndex].Towers,
					}

					network.SendMessage(p.User.Conn, network.Message{Type: config.MsgMatchUpdate, Data: matchData})
					network.SendMessage(p.User.Conn, rs)

					// All troops destroyed
					troopCount := 0
					for _, troop := range e.Players[playerIndex].Troops {
						if troop.HP <= 0 {
							troopCount++
						}
					}
					if troopCount >= 3 {
						e.KingOrTroopsDestroyed <- []*Player{e.Players[opponentIndex], e.Players[playerIndex]}
						return
					}

					// King destroyed
					if e.Players[opponentIndex].Towers[2].HP <= 0 {
						e.KingOrTroopsDestroyed <- []*Player{e.Players[playerIndex], e.Players[opponentIndex]}
						return
					}
				}
			}
		case msg := <-e.Players[1].User.Talk:
			var command Command
			json.Unmarshal(msg.Data.(json.RawMessage), &command)

			var rs network.Message
			var err error

			if command.TowerIndex == 5 && command.TroopIndex == 5 {
				// Healing
				rs, err = e.Players[1].Healing()
			} else {
				// Attack
				rs, err = e.Players[1].Attack(e.Players[0], &command)
			}
			if err != nil {
				network.SendMessage(e.Players[1].User.Conn, rs)
			} else {
				for i, p := range e.Players {
					var playerIndex int
					var opponentIndex int
					if i == 0 {
						playerIndex = 0
						opponentIndex = 1
					} else {
						playerIndex = 1
						opponentIndex = 0
					}

					matchData := MatchData{
						PUsername:      e.Players[playerIndex].User.Metadata.Username,
						PLevelMetadata: e.Players[playerIndex].User.Metadata.Level,
						PEXPMetadata:   e.Players[playerIndex].User.Metadata.EXP,
						PLevel:         e.Players[playerIndex].Level,
						PEXP:           e.Players[playerIndex].EXP,
						PMana:          e.Players[playerIndex].Mana,
						PTroops:        e.Players[playerIndex].Troops,
						PTowers:        e.Players[playerIndex].Towers,
						OUsername:      e.Players[opponentIndex].User.Metadata.Username,
						OLevelMetadata: e.Players[opponentIndex].User.Metadata.Level,
						OEXPMetadata:   e.Players[opponentIndex].User.Metadata.EXP,
						OLevel:         e.Players[opponentIndex].Level,
						OEXP:           e.Players[opponentIndex].EXP,
						OMana:          e.Players[opponentIndex].Mana,
						OTroops:        e.Players[opponentIndex].Troops,
						OTowers:        e.Players[opponentIndex].Towers,
					}

					network.SendMessage(p.User.Conn, network.Message{Type: config.MsgMatchUpdate, Data: matchData})
					network.SendMessage(p.User.Conn, rs)

					// All troops destroyed
					troopCount := 0
					for _, troop := range e.Players[playerIndex].Troops {
						if troop.HP <= 0 {
							troopCount++
						}
					}
					if troopCount >= 3 {
						e.KingOrTroopsDestroyed <- []*Player{e.Players[opponentIndex], e.Players[playerIndex]}
						return
					}

					// King destroyed
					if e.Players[opponentIndex].Towers[2].HP <= 0 {
						e.KingOrTroopsDestroyed <- []*Player{e.Players[playerIndex], e.Players[opponentIndex]}
						return
					}
				}
			}
		case <-e.Players[0].User.Interrupt:
			fmt.Printf("Player %s disconnected\n", e.Players[0].User.Metadata.Username)

			e.Interrupt <- e.Players[1].User
			return
		case <-e.Players[1].User.Interrupt:
			fmt.Printf("Player %s disconnected\n", e.Players[1].User.Metadata.Username)

			e.Interrupt <- e.Players[0].User
			return
		}
	}
}

func handleGameEnd(e *Engine) {
	select {
	case user := <-e.Interrupt:
		e.TimerEnd <- true

		for i := range 5 {
			network.SendMessage(user.Conn, network.Message{Type: config.MsgError, Data: []string{
				"The opponen disconnected!",
				"You will be moved to the queue in",
				fmt.Sprintf("%d seconds...", 5-i),
			}})

			time.Sleep(1 * time.Second)
		}
		if e.OnRequeue != nil {
			e.OnRequeue(user)
		}
		network.SendMessage(user.Conn, network.Message{Type: config.MsgError, Data: []string{
			"Waiting for other player...",
		}})
	case data := <-e.KingOrTroopsDestroyed:
		fmt.Println(data)
		e.TimerEnd <- true

		winner := data[0]
		loser := data[1]

		if winner != nil {
			winner.User.Metadata.EXP += 30.0
			loser.User.Metadata.EXP += 10.0
			doesUpgradeLevel(winner)
			doesUpgradeLevel(loser)

			winner.User.Metadata.SaveAll()
			loser.User.Metadata.SaveAll()

			go func() {
				for i := range 5 {
					network.SendMessage(winner.User.Conn, network.Message{Type: config.MsgMatchEnd, Data: []string{
						fmt.Sprintf("VICTORY! (Next round in %d seconds...)", 5-i),
						fmt.Sprintf("Winner: %s", winner.User.Metadata.Username),
						fmt.Sprintf("+%d EXP => Level: %d", 30, int(winner.User.Metadata.Level)),
					}})
					time.Sleep(1 * time.Second)
				}
				e.OnRequeue(winner.User)
			}()

			go func() {
				for i := range 5 {
					network.SendMessage(loser.User.Conn, network.Message{Type: config.MsgMatchEnd, Data: []string{
						fmt.Sprintf("DEFEAT! (Next round in %d seconds...)", 5-i),
						fmt.Sprintf("Winner: %s", winner.User.Metadata.Username),
						fmt.Sprintf("+%d EXP => Level: %d", 10, int(loser.User.Metadata.Level)),
					}})
					time.Sleep(1 * time.Second)
				}
				e.OnRequeue(loser.User)
			}()
		} else {
			// Draw
			player1 := e.Players[0]
			player2 := e.Players[1]

			player1.User.Metadata.EXP += 10
			player2.User.Metadata.EXP += 10

			doesUpgradeLevel(player1)
			doesUpgradeLevel(player2)

			player1.User.Metadata.SaveAll()
			player2.User.Metadata.SaveAll()

			go func() {
				for i := range 5 {
					network.SendMessage(player1.User.Conn, network.Message{Type: config.MsgMatchEnd, Data: []string{
						fmt.Sprintf("DRAW! (Next round in %d seconds...)", 5-i),
						fmt.Sprintf("+%d EXP => Level: %d", 10, int(player1.User.Metadata.Level)),
					}})
					time.Sleep(1 * time.Second)
				}
				e.OnRequeue(player1.User)
			}()

			go func() {
				for i := range 5 {
					network.SendMessage(player2.User.Conn, network.Message{Type: config.MsgMatchEnd, Data: []string{
						fmt.Sprintf("DRAW! (Next round in %d seconds...)", 5-i),
						fmt.Sprintf("+%d EXP => Level: %d", 10, int(player2.User.Metadata.Level)),
					}})
					time.Sleep(1 * time.Second)
				}
				e.OnRequeue(player2.User)
			}()
		}
	case <-e.TimerEnd:
		winner, loser := whoWin(e.Players[0], e.Players[1])

		if winner != nil {
			winner.User.Metadata.EXP += 30.0
			loser.User.Metadata.EXP += 10.0
			doesUpgradeLevel(winner)
			doesUpgradeLevel(loser)

			winner.User.Metadata.SaveAll()
			loser.User.Metadata.SaveAll()

			go func() {
				for i := range 5 {
					network.SendMessage(winner.User.Conn, network.Message{Type: config.MsgMatchEnd, Data: []string{
						fmt.Sprintf("VICTORY! (Next round in %d seconds...)", 5-i),
						fmt.Sprintf("Winner: %s", winner.User.Metadata.Username),
						fmt.Sprintf("+%d EXP => Level: %d", 30, int(winner.User.Metadata.Level)),
					}})
					time.Sleep(1 * time.Second)
				}
				if e.OnRequeue != nil {
					e.OnRequeue(winner.User)
				}
			}()

			go func() {
				for i := range 5 {
					network.SendMessage(loser.User.Conn, network.Message{Type: config.MsgMatchEnd, Data: []string{
						fmt.Sprintf("DEFEAT! (Next round in %d seconds...)", 5-i),
						fmt.Sprintf("Winner: %s", winner.User.Metadata.Username),
						fmt.Sprintf("+%d EXP => Level: %d", 10, int(loser.User.Metadata.Level)),
					}})
					time.Sleep(1 * time.Second)
				}
				if e.OnRequeue != nil {
					e.OnRequeue(loser.User)
				}
			}()
		} else {
			// Draw
			player1 := e.Players[0]
			player2 := e.Players[1]

			player1.User.Metadata.EXP += 10
			player2.User.Metadata.EXP += 10

			doesUpgradeLevel(player1)
			doesUpgradeLevel(player2)

			player1.User.Metadata.SaveAll()
			player2.User.Metadata.SaveAll()

			go func() {
				for i := range 5 {
					network.SendMessage(player1.User.Conn, network.Message{Type: config.MsgMatchEnd, Data: []string{
						fmt.Sprintf("DRAW! (Next round in %d seconds...)", 5-i),
						fmt.Sprintf("+%d EXP => Level: %d", 10, int(player1.User.Metadata.Level)),
					}})
					time.Sleep(1 * time.Second)
				}
				if e.OnRequeue != nil {
					e.OnRequeue(player1.User)
				}
			}()

			go func() {
				for i := range 5 {
					network.SendMessage(player2.User.Conn, network.Message{Type: config.MsgMatchEnd, Data: []string{
						fmt.Sprintf("DRAW! (Next round in %d seconds...)", 5-i),
						fmt.Sprintf("+%d EXP => Level: %d", 10, int(player2.User.Metadata.Level)),
					}})
					time.Sleep(1 * time.Second)
				}
				if e.OnRequeue != nil {
					e.OnRequeue(player2.User)
				}
			}()
		}
	}
}

func whoWin(p1, p2 *Player) (winner, loser *Player) {
	towerDestroyed1 := 0
	towerDestroyed2 := 0

	for i := range 3 {
		t1 := p1.Towers[i]
		t2 := p2.Towers[i]

		if t1.HP <= 0 {
			towerDestroyed1++
		}
		if t2.HP <= 0 {
			towerDestroyed2++
		}
	}
	if towerDestroyed1 < towerDestroyed2 {
		return p1, p2
	} else if towerDestroyed1 > towerDestroyed2 {
		return p2, p1
	}

	return nil, nil
}

func doesUpgradeLevel(p *Player) {
	baseEXP := 100.0

	currentLevel := p.User.Metadata.Level
	currentEXP := p.User.Metadata.EXP

	requiredEXP := baseEXP * math.Pow(1.1, currentLevel-1)
	remainEXP := currentEXP - requiredEXP

	if remainEXP >= 0 {
		// Can upgrade
		p.User.Metadata.Level++
		p.User.Metadata.EXP = remainEXP

		for _, troop := range p.User.Metadata.Troops {
			troop.ATK *= 1.1
			troop.DEF *= 1.1
			troop.EXP *= 1.1
			troop.HP *= 1.1
			troop.Mana *= 1.1
		}

		for _, tower := range p.User.Metadata.Towers {
			tower.ATK *= 1.1
			tower.DEF *= 1.1
			tower.EXP *= 1.1
			tower.HP *= 1.1
			tower.Crit *= 1.1
		}
	}
}
