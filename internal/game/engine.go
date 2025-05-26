package game

import (
	"bufio"
	"clash-royale/internal/network"
	"log"
	"math/rand"
	"time"
)

type Engine struct {
	Players []*Player
	Tick    int
	Talk    chan network.Message
}

func (e *Engine) ListenUser(u *User) error {
	defer u.Conn.Close()
	reader := bufio.NewReader(u.Conn)
	for {
		msg, err := network.ReceiveMessage(reader)
		if err != nil {
			log.Printf("[%s] disconnected\n", u.Metadata.Username)
			return err
		}

		e.Talk <- msg
	}
}

func NewEngine(u1, u2 *User) *Engine {
	p1 := Player{
		User:   *u1,
		Mana:   0.0,
		EXP:    0.0,
		Troops: randomTroop(u1.Metadata.Troops),
		Towers: u1.Metadata.Towers,
	}

	p2 := Player{
		User:   *u2,
		Mana:   0.0,
		EXP:    0.0,
		Troops: randomTroop(u2.Metadata.Troops),
		Towers: u2.Metadata.Towers,
	}

	return &Engine{
		Players: []*Player{&p1, &p2},
		Tick:    0,
		Talk:    make(chan network.Message, 10), // Ignore spam
	}
}

func randomTroop(troops []Troop) []Troop {
	rand.Shuffle(len(troops), func(i, j int) {
		troops[i], troops[j] = troops[j], troops[i]
	})
	return troops[:3]
}

func (e *Engine) Start(duration time.Duration) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		e.Tick++

		// Regen mana
		// for _, p := range e.Players {

		// }

		network.SendMessage(e.Players[0].User.Conn, network.Message{Type: "demo", Data: 123})
		network.SendMessage(e.Players[1].User.Conn, network.Message{Type: "demo", Data: 123})

		// Check duration
		if e.Tick >= int(duration.Seconds()) {
			break
		}
	}
}
