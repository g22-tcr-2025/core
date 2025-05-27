package game

import (
	"clash-royale/internal/config"
	"clash-royale/internal/network"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Engine struct {
	Players   []*Player
	Tick      int
	End       chan bool
	OnRequeue func(u *User)
}

func NewEngine(u1, u2 *User, onRequeue func(u *User)) *Engine {
	p1 := Player{
		User:   *u1,
		Mana:   5.0,
		EXP:    0.0,
		Troops: randomTroop(u1.Metadata.Troops),
		Towers: u1.Metadata.Towers,
		Mutex:  sync.Mutex{},
	}

	p2 := Player{
		User:   *u2,
		Mana:   5.0,
		EXP:    0.0,
		Troops: randomTroop(u2.Metadata.Troops),
		Towers: u2.Metadata.Towers,
		Mutex:  sync.Mutex{},
	}

	return &Engine{
		Players:   []*Player{&p1, &p2},
		Tick:      0,
		End:       make(chan bool),
		OnRequeue: onRequeue,
	}
}

func randomTroop(troops []Troop) []Troop {
	rand.Shuffle(len(troops), func(i, j int) {
		troops[i], troops[j] = troops[j], troops[i]
	})
	return troops[:3]
}

func (e *Engine) Start() {
	p1 := MatchData{
		PUsername: e.Players[0].User.Metadata.Username,
		PLevel:    e.Players[0].User.Metadata.Level,
		PMana:     e.Players[0].Mana,
		PTroops:   e.Players[0].Troops,
		PTowers:   e.Players[0].Towers,
		OUsername: e.Players[1].User.Metadata.Username,
		OLevel:    e.Players[1].User.Metadata.Level,
		OMana:     e.Players[1].Mana,
		OTroops:   e.Players[1].Troops,
		OTowers:   e.Players[1].Towers,
	}
	p2 := MatchData{
		PUsername: e.Players[1].User.Metadata.Username,
		PLevel:    e.Players[1].User.Metadata.Level,
		PMana:     e.Players[1].Mana,
		PTroops:   e.Players[1].Troops,
		PTowers:   e.Players[1].Towers,
		OUsername: e.Players[0].User.Metadata.Username,
		OLevel:    e.Players[0].User.Metadata.Level,
		OMana:     e.Players[0].Mana,
		OTroops:   e.Players[0].Troops,
		OTowers:   e.Players[0].Towers,
	}

	network.SendMessage(e.Players[0].User.Conn, network.Message{Type: config.MsgMatchStart, Data: p1})
	network.SendMessage(e.Players[1].User.Conn, network.Message{Type: config.MsgMatchStart, Data: p2})

	go runtime(e)
	go handleCommand(e)
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

			network.SendMessage(e.Players[0].User.Conn, network.Message{Type: config.MsgUpdateMnana, Data: e.Players[0].Mana})
			network.SendMessage(e.Players[1].User.Conn, network.Message{Type: config.MsgUpdateMnana, Data: e.Players[1].Mana})

			// Check duration
			if e.Tick >= int(config.MatchDuration.Seconds()) {
				return
			}
		case <-e.End:
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

			rs := e.Players[0].Attack(e.Players[1], &command)
			for _, p := range e.Players {
				network.SendMessage(p.User.Conn, rs)
			}
		case msg := <-e.Players[1].User.Talk:
			var command Command
			json.Unmarshal(msg.Data.(json.RawMessage), &command)
			rs := e.Players[1].Attack(e.Players[0], &command)

			for _, p := range e.Players {
				network.SendMessage(p.User.Conn, rs)
			}
		case <-e.Players[0].User.Interrupt:
			fmt.Println("Player 1 disconnected")
			e.End <- true
			if e.OnRequeue != nil {
				e.OnRequeue(&e.Players[1].User)
			}
			return
		case <-e.Players[1].User.Interrupt:
			fmt.Println("Player 2 disconnected")
			e.End <- true
			if e.OnRequeue != nil {
				e.OnRequeue(&e.Players[0].User)
			}
			return
		}
	}
}
