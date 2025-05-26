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
	queue     []*game.User
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

	data.EnsureMetadata(loginData.Username)
	userMetadata, _ := data.LoadMetadata(loginData.Username)

	user := game.User{
		Conn:     conn,
		Metadata: *userMetadata,
	}

	m.mutext.Lock()
	m.queue = append(m.queue, &user)
	if len(m.queue) >= 2 {
		u1 := m.queue[0]
		u2 := m.queue[1]

		m.queue = m.queue[2:]

		// Make new engine
		engine := game.NewEngine(u1, u2)

		go engine.ListenUser(u1)
		go engine.ListenUser(u2)

		engine.Start(config.MatchDuration)
	}
	m.mutext.Unlock()
}
