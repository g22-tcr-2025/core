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
		conn.Close()
		return
	}

	// After login successfully
	data.EnsureMetadata(loginData.Username)
	userMetadata, _ := data.LoadMetadata(loginData.Username)

	user := game.User{
		Conn:      conn,
		Metadata:  userMetadata,
		Talk:      make(chan network.Message),
		Interrupt: make(chan bool),
	}

	go user.ListenUser()

	tryMatch(&user, m)
}

func tryMatch(u *game.User, m *MatchMaker) {
	m.mutext.Lock()
	m.queue = append(m.queue, u)
	m.mutext.Unlock()

	if len(m.queue) >= 2 {
		m.mutext.Lock()
		u1 := m.queue[0]
		u2 := m.queue[1]
		m.queue = m.queue[2:]
		m.mutext.Unlock()

		go game.NewEngine(u1, u2, func(u *game.User) {
			tryMatch(u, m)
		}).Start()
	}
}
