package ui

import (
	"bufio"
	"clash-royale/internal/config"
	"clash-royale/internal/game"
	"clash-royale/internal/network"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func ListenServer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := network.ReceiveMessage(reader)
		if err != nil {
			log.Println("Server stopped.")
			return
		}
		switch msg.Type {
		case config.MsgMatchStart:
		case config.MsgStateUpdate:
			var mana float64
			json.Unmarshal(msg.Data.(json.RawMessage), &mana)
			RenderMana(mana)
			// ClearScreen()
		case config.MsgMatchEnd:
		}
	}
}

func ListenPlayer(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input:", err)
		}
		line = strings.TrimSpace(line)
		network.SendMessage(conn, network.Message{Type: "demo", Data: line})
		ClearInput()
	}
}

func LoginStep(conn net.Conn) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter USERNAME: ")
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	username := strings.TrimSpace(line)

	fmt.Printf("Enter PASSWORD: ")
	line, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	password := strings.TrimSpace(line)

	data := game.LoginData{Username: username, Password: password}
	msg := network.Message{Type: config.MsgLogin, Data: data}

	err = network.SendMessage(conn, msg)
	if err != nil {
		return err
	}
	// Response from server
	readerServer := bufio.NewReader(conn)
	msg, err = network.ReceiveMessage(readerServer)
	if err != nil {
		return err
	}
	if msg.Type != config.MsgLoginResponse {
		return fmt.Errorf("INVALID RESPONSE")
	}
	var ok bool
	json.Unmarshal(msg.Data.(json.RawMessage), &ok)

	if !ok {
		return fmt.Errorf("INVALID CREDENTIALS")
	}

	return nil
}

func ClearScreen() {
	fmt.Print("\033[2J\033[H")
	fmt.Println("Mana: 0")
	fmt.Print(">> ")
}

func ClearInput() {
	fmt.Print("\033[2;1H")
	fmt.Print("\033[K")
	fmt.Print(">> ")
}

func RenderMana(mana float64) {
	fmt.Print("\033[s")            // Save pointer
	fmt.Print("\033[1;1H")         // Move to line 1 col 1
	fmt.Print("\033[K")            // Clear line
	fmt.Printf("Mana: %.1f", mana) // Print mana
	fmt.Print("\033[u")            // Back to previous
}
