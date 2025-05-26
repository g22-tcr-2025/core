package data

import (
	"clash-royale/internal/game"
	"encoding/json"
	"log"
	"os"
)

type UserStore struct {
	Users map[string]string
}

func LoadUsers(path string) *UserStore {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read users data: %v", err)
	}

	var userList []game.LoginData
	if err := json.Unmarshal(data, &userList); err != nil {
		log.Fatalf("Failed to parse user data: %v", err)
	}

	users := make(map[string]string)
	for _, u := range userList {
		users[u.Username] = u.Password
	}

	log.Printf("✅ Loaded %d users\n", len(users))
	return &UserStore{Users: users}
}

func (us *UserStore) Validate(loginData game.LoginData) bool {
	pass, ok := us.Users[loginData.Username]
	return ok && pass == loginData.Password
}
