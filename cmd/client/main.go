package main

import (
	"clash-royale/internal/network"
	"clash-royale/internal/ui"
	"log"
)

func main() {
	log.Println("Starting client...")
	conn := network.ConnectToServer("localhost:3000")
	defer conn.Close()

	err := ui.LoginStep(conn)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	ui.ClearScreen()

	go ui.ListenPlayer(conn)
	go ui.ListenServer(conn)

	select {}
}
