package logic

import (
	"bufio"
	"clash-royale/internal/config"
	"clash-royale/internal/data"
	"clash-royale/internal/game"
	"clash-royale/internal/network"
	"encoding/json"
	"log"
	"net"
	"sync"
)

type MatchMaker struct {
	UserStore *data.UserStore
	queue     []*game.Player
	mutext    sync.Mutex
}

func (m *MatchMaker) HandleConnection(conn net.Conn) {
	log.Println("New player connected:", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	// LOGIN step
	msg, err := network.ReceiveMessage(reader)
	if err != nil {
		log.Println(conn.RemoteAddr(), "disconnected!")
		return
	}
	if msg.Type != config.MsgLogin {
		log.Println(conn.RemoteAddr(), "must login first!")
		return
	}
	var loginData game.LoginData
	json.Unmarshal(msg.Data.(json.RawMessage), &loginData)

	ok := m.UserStore.Validate(loginData)
	network.SendMessage(conn, network.Message{Type: config.MsgLoginResponse, Data: ok})
	if !ok {
		return
	}

	data.EnsurePlayerMetadata(loginData.Username)
	playerData, _ := data.LoadPlayerData(loginData.Username)

	player := game.Player{
		Conn: conn,
		Data: *playerData,
	}

	m.mutext.Lock()
	m.queue = append(m.queue, &player)
	if len(m.queue) >= 2 {
		p1 := m.queue[0]
		p2 := m.queue[1]

		m.queue = m.queue[2:]

		// Make new engine
		engine := game.NewEngine(p1, p2)

		go engine.ListenPlayer(p1)
		go engine.ListenPlayer(p2)

		engine.Start(config.MatchDuration)
	}
	m.mutext.Unlock()
}
