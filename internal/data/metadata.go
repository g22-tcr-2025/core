package data

import (
	"clash-royale/internal/game"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type templateStandard struct {
	Troops []*game.Troop `json:"troops"`
	Towers []*game.Tower `json:"towers"`
}

func EnsureMetadata(username string) error {
	metadataDir := filepath.Join("assets", "metadata", username)
	if _, err := os.Stat(metadataDir); os.IsNotExist(err) {
		os.MkdirAll(metadataDir, 0755)

		var tlStd templateStandard
		loadJSON(filepath.Join("assets", "standard.json"), &tlStd)

		metadata := game.UserMetadata{Username: username, EXP: 0.0, Level: 1.0, Troops: tlStd.Troops, Towers: tlStd.Towers}
		data, _ := json.MarshalIndent(metadata, "", "  ")
		if err := os.WriteFile(filepath.Join(metadataDir, "metadata.json"), data, 0644); err != nil {
			return err
		}
	}
	return nil
}

func LoadMetadata(username string) (*game.UserMetadata, error) {
	metadataDir := filepath.Join("assets", "metadata", username)
	if _, err := os.Stat(metadataDir); os.IsNotExist(err) {
		return nil, errors.New("user metadata not found")
	}

	metadata := &game.UserMetadata{}
	if err := loadJSON(filepath.Join(metadataDir, "metadata.json"), metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

func loadJSON(path string, v any) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, v)
}
