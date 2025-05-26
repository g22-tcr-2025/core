package main

import (
	"clash-royale/internal/data"
	"clash-royale/internal/logic"
	"clash-royale/internal/network"
	"log"
)

func main() {
	log.Println("Starting server...")

	userStore := data.LoadUsers("assets/users.json")
	network.StartServer(3000, &logic.MatchMaker{UserStore: userStore})
}
