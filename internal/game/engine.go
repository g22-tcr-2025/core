package game

import (
	"clash-royale/internal/config"
	"clash-royale/internal/network"
	"fmt"
	"log"
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
	log.Println("Match started")
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

			network.SendMessage(e.Players[0].User.Conn, network.Message{Type: config.MsgStateUpdate, Data: e.Players[0].Mana})
			network.SendMessage(e.Players[1].User.Conn, network.Message{Type: config.MsgStateUpdate, Data: e.Players[1].Mana})

			// Check duration
			if e.Tick >= int(config.MatchDuration.Seconds()) {
				e.End <- true
				return
			}
		case <-e.End:
			log.Println("Match ended")
			return
		}
	}
}

func handleCommand(e *Engine) {
	for {
		select {
		case msg := <-e.Players[0].User.Talk:
			fmt.Println(msg)
		case msg := <-e.Players[1].User.Talk:
			fmt.Println(msg)
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
