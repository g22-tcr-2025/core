package data

import (
	"clash-royale/internal/game"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

func EnsurePlayerMetadata(username string) error {
	userDir := filepath.Join("assets/users", username)
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		if err := os.MkdirAll(userDir, 0755); err != nil {
			return err
		}
		if err := copyFile("assets/sd_troops.json", filepath.Join(userDir, "troops.json")); err != nil {
			return err
		}
		if err := copyFile("assets/sd_towers.json", filepath.Join(userDir, "towers.json")); err != nil {
			return err
		}

		meta := game.PlayerMetadata{Username: username, EXP: 0.0, Level: 1.0}
		data, _ := json.MarshalIndent(meta, "", "  ")
		if err := os.WriteFile(filepath.Join(userDir, "player.json"), data, 0644); err != nil {
			return err
		}
	}
	return nil
}

func LoadPlayerData(username string) (*game.PlayerDataManagement, error) {
	userDir := filepath.Join("assets/users", username)
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		return nil, errors.New("player data not found")
	}

	meta := &game.PlayerMetadata{}
	if err := loadJSON(filepath.Join(userDir, "player.json"), meta); err != nil {
		return nil, err
	}

	troops := []game.Troop{}
	towers := []game.Tower{}
	loadJSON(filepath.Join(userDir, "troops.json"), &troops)
	loadJSON(filepath.Join(userDir, "towers.json"), &towers)

	return &game.PlayerDataManagement{
		Metadata: *meta,
		Troops:   troops,
		Towers:   towers,
	}, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func loadJSON(path string, v any) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, v)
}
