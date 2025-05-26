package game

import (
	"bufio"
	"clash-royale/internal/config"
	"clash-royale/internal/network"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Engine struct {
	Players   []*Player
	Tick      int
	Talk      chan network.Message
	End       chan bool
	Interrupt chan bool
}

func (e *Engine) ListenUser(u *User) error {
	defer u.Conn.Close()
	reader := bufio.NewReader(u.Conn)
	for {
		msg, err := network.ReceiveMessage(reader)
		if err != nil {
			log.Printf("[%s] disconnected\n", u.Metadata.Username)
			e.Interrupt <- true
			return err
		}

		e.Talk <- msg
	}
}

func NewEngine(u1, u2 *User) *Engine {
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
		Talk:      make(chan network.Message, 10), // Ignore spam
		End:       make(chan bool),
		Interrupt: make(chan bool),
	}
}

func randomTroop(troops []Troop) []Troop {
	rand.Shuffle(len(troops), func(i, j int) {
		troops[i], troops[j] = troops[j], troops[i]
	})
	return troops[:3]
}

func (e *Engine) Start(duration time.Duration) {
	go runtime(duration, e)

	for {
		select {
		case <-e.Talk:
			e.Players[0].Mutex.Lock()
			e.Players[0].Mana -= 10
			e.Players[0].Mutex.Unlock()
			network.SendMessage(e.Players[0].User.Conn, network.Message{Type: config.MsgStateUpdate, Data: e.Players[0].Mana})
			network.SendMessage(e.Players[1].User.Conn, network.Message{Type: config.MsgStateUpdate, Data: e.Players[1].Mana})
		case <-e.End:
			network.SendMessage(e.Players[0].User.Conn, network.Message{Type: "demo", Data: "end"})
			network.SendMessage(e.Players[1].User.Conn, network.Message{Type: "demo", Data: "end"})
			return
		}
	}
}

func runtime(duration time.Duration, e *Engine) {
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
			if e.Tick >= int(duration.Seconds()) {
				e.End <- true
				return
			}
		case <-e.Interrupt:
			e.End <- true
			return
		}
	}
}
