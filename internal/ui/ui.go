package ui

import (
	"bufio"
	"clash-royale/internal/config"
	"clash-royale/internal/game"
	"clash-royale/internal/network"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
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
			var template game.MatchData
			json.Unmarshal(msg.Data.(json.RawMessage), &template)
			RenderTemplate(template)
		case config.MsgUpdateMnana:
			var mana float64
			json.Unmarshal(msg.Data.(json.RawMessage), &mana)
			RenderMana(mana)
		case config.MsgAttackResult:
			var combatResult game.CombatResult
			json.Unmarshal(msg.Data.(json.RawMessage), &combatResult)

			RenderNotification(combatString(combatResult))
		case config.MsgError:
			var err string
			json.Unmarshal(msg.Data.(json.RawMessage), &err)
			RenderNotification(err)
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

		data, err := validCommand(line)
		if err != nil {
			RenderNotification(err.Error())
		} else {
			network.SendMessage(conn, network.Message{Type: config.MsgAttack, Data: data})
		}
		ClearInput()
	}
}

func validCommand(line string) (game.Command, error) {
	data := strings.Split(line, " ")
	if len(data) != 2 {
		return game.Command{}, errors.New("invalid command")
	}

	indexTroop, err := strconv.Atoi(data[0])
	if err != nil {
		return game.Command{}, err
	}

	indexTower, err := strconv.Atoi(data[1])
	if err != nil {
		return game.Command{}, err
	}

	return game.Command{TroopIndex: indexTroop, TowerIndex: indexTower}, nil
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
}

func ClearInput() {
	fmt.Print("\033[28;1H")
	fmt.Print("\033[K")
	fmt.Print(">> ")
}

func RenderTemplate(matchData game.MatchData) {
	ClearScreen()

	fmt.Println("============ You ============")
	fmt.Printf("============ %s - %d ============\n", matchData.PUsername, int(matchData.PLevel))
	fmt.Println(manaString(matchData.PMana))
	fmt.Println()
	for i, troop := range matchData.PTroops {
		fmt.Println(troopString(i, troop))
	}
	fmt.Println()
	for i, tower := range matchData.PTowers {
		fmt.Println(towerString(i, tower))
	}
	fmt.Println()

	fmt.Println("============ Opponent ============")
	fmt.Printf("============ %s - %d ============\n", matchData.OUsername, int(matchData.OLevel))
	fmt.Println(manaString(matchData.OMana))
	fmt.Println()
	for i, troop := range matchData.OTroops {
		fmt.Println(troopString(i, troop))
	}
	fmt.Println()
	for i, tower := range matchData.OTowers {
		fmt.Println(towerString(i, tower))
	}
	fmt.Println()

	fmt.Println(":: ")
	fmt.Println()

	fmt.Println("Command: <troop_index> <tower_index>")
	fmt.Print(">> ")
}

func manaString(mana float64) string {
	manaInt := int(mana)
	str := ""
	for range manaInt {
		str += "#"
	}
	return str
}

func troopString(index int, troop game.Troop) string {
	str := ""
	str += fmt.Sprintf("[%d]", index)
	if troop.HP <= 0 {
		str += fmt.Sprintf(" 🪦 %s", troop.Name)
	} else {
		str += fmt.Sprintf(" 🤖 %s", troop.Name)
	}
	str += fmt.Sprintf("\t\t❤️ %d", int(troop.HP))
	str += fmt.Sprintf("\t🛡️ %d", int(troop.DEF))
	str += fmt.Sprintf("\t⚔️ %d", int(troop.ATK))

	return str
}

func combatString(combatResult game.CombatResult) string {
	str := combatResult.Attacker
	str += fmt.Sprintf(" 🎯 %s", combatResult.Defender)
	str += fmt.Sprintf(" | 🤖 %s ⛏️  %s 🏰", combatResult.UsingTroop.Name, combatResult.TargetTower.Type)
	str += fmt.Sprintf(" | 🤖 (-%d🩸) ~ 🏰 (-%d🩸)", int(math.Ceil(combatResult.DamgeToTroop)), int(math.Ceil(combatResult.DamgeToTower)))

	return str
}

func towerString(index int, tower game.Tower) string {
	str := ""
	str += fmt.Sprintf("[%d]", index)
	if tower.HP <= 0 {
		str += fmt.Sprintf(" 🪨 %s", tower.Type)
	} else {
		str += fmt.Sprintf(" 🏰 %s", tower.Type)
	}
	str += fmt.Sprintf("\t❤️ %d", int(tower.HP))
	str += fmt.Sprintf("\t🛡️ %d", int(tower.DEF))
	str += fmt.Sprintf("\t⚔️ %d", int(tower.ATK))

	return str
}

func RenderMana(mana float64) {
	fmt.Print("\033[s")         // Save pointer
	fmt.Print("\033[3;1H")      // Move to line 3 col 1
	fmt.Print("\033[K")         // Clear line
	fmt.Print(manaString(mana)) // Print mana

	fmt.Print("\033[15;1H")     // Move to line 3 col 1
	fmt.Print("\033[K")         // Clear line
	fmt.Print(manaString(mana)) // Print mana
	fmt.Print("\033[u")         // Back to previous
}

func RenderNotification(content string) {
	fmt.Print("\033[s")       // Save pointer
	fmt.Print("\033[25;1H")   // Move to line 3 col 1
	fmt.Print("\033[K")       // Clear line
	fmt.Print(":: ", content) // Print mana
	fmt.Print("\033[u")       // Back to previous
}
