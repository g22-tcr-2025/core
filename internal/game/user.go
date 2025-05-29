package game

import (
	"bufio"
	"clash-royale/internal/network"
	"encoding/json"
	"log"
	"net"
	"os"
	"path/filepath"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserMetadata struct {
	Username string   `json:"username"`
	EXP      float64  `json:"exp"`
	Level    float64  `json:"level"`
	Troops   []*Troop `json:"troops"`
	Towers   []*Tower `json:"towers"`
}

type User struct {
	Conn      net.Conn
	Metadata  *UserMetadata
	Talk      chan network.Message
	Interrupt chan bool
}

func (u *User) ListenUser() error {
	defer u.Conn.Close()
	reader := bufio.NewReader(u.Conn)
	for {
		msg, err := network.ReceiveMessage(reader)
		if err != nil {
			log.Printf("[%s] disconnected\n", u.Metadata.Username)

			u.Interrupt <- true
			return err
		}

		u.Talk <- msg
	}
}

func (um *UserMetadata) SaveAll() error {
	usersDir := filepath.Join("assets", "metadata", um.Username)
	if err := saveJSON(filepath.Join(usersDir, "metadata.json"), um); err != nil {
		return err
	}
	return nil
}

func saveJSON(path string, v any) error {
	data, _ := json.MarshalIndent(v, "", "  ")
	return os.WriteFile(path, data, 0644)
}
