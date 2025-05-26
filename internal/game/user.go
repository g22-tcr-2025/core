package game

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserMetadata struct {
	Username string  `json:"username"`
	EXP      float64 `json:"exp"`
	Level    float64 `json:"level"`
	Troops   []Troop `json:"troops"`
	Towers   []Tower `json:"towers"`
}

type User struct {
	Conn     net.Conn
	Metadata UserMetadata
}

func (um *UserMetadata) SaveAll() error {
	usersDir := filepath.Join("assets/metadata")
	if err := saveJSON(filepath.Join(usersDir, um.Username+".json"), um); err != nil {
		return err
	}
	return nil
}

func saveJSON(path string, v any) error {
	data, _ := json.MarshalIndent(v, "", "  ")
	return os.WriteFile(path, data, 0644)
}
