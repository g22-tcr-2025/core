package game

import (
	"encoding/json"
	"math/rand"
	"net"
	"os"
	"path/filepath"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PlayerMetadata struct {
	Username string  `json:"username"`
	EXP      float64 `json:"exp"`
	Level    float64 `json:"level"`
}

type PlayerDataManagement struct {
	Metadata PlayerMetadata `json:"metadata"`
	Troops   []Troop        `json:"troops"`
	Towers   []Tower        `json:"towers"`
}

// Just use for server logic
type Player struct {
	Conn net.Conn
	Data PlayerDataManagement
}

func (p *Player) HasCrit() bool {
	// 10% lucky
	return rand.Float64() < 0.1
}

func (pdm *PlayerDataManagement) SaveAll() error {
	userDir := filepath.Join("assets/users", pdm.Metadata.Username)
	if err := saveJSON(filepath.Join(userDir, "player.json"), pdm.Metadata); err != nil {
		return err
	}
	if err := saveJSON(filepath.Join(userDir, "troops.json"), pdm.Troops); err != nil {
		return err
	}
	if err := saveJSON(filepath.Join(userDir, "towers.json"), pdm.Towers); err != nil {
		return err
	}
	return nil
}

func saveJSON(path string, v any) error {
	data, _ := json.MarshalIndent(v, "", "  ")
	return os.WriteFile(path, data, 0644)
}
