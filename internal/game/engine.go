package game

import (
	"bufio"
	"clash-royale/internal/network"
	"log"
	"time"
)

type Engine struct {
	Players []*Player
	Tick    int
	Talk    chan network.Message
}

func NewEngine(p1, p2 *Player) *Engine {
	return &Engine{
		Players: []*Player{p1, p2},
		Tick:    0,
		Talk:    make(chan network.Message, 10), // Ignore spam
	}
}

func (e *Engine) ListenPlayer(p *Player) error {
	defer p.Conn.Close()
	reader := bufio.NewReader(p.Conn)
	for {
		msg, err := network.ReceiveMessage(reader)
		if err != nil {
			log.Printf("[%s] disconnected\n", p.Data.Metadata.Username)
			return err
		}

		e.Talk <- msg
	}
}

func (e *Engine) Start(duration time.Duration) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		e.Tick++

		// Regen mana
		// for _, p := range e.Players {

		// }

		network.SendMessage(e.Players[0].Conn, network.Message{Type: "demo", Data: 123})
		network.SendMessage(e.Players[1].Conn, network.Message{Type: "demo", Data: 123})

		// Check duration
		if e.Tick >= int(duration.Seconds()) {
			break
		}
	}
}
