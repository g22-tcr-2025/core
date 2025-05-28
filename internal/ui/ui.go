package ui

import (
	"bufio"
	"clash-royale/internal/config"
	"clash-royale/internal/game"
	"clash-royale/internal/network"
	"clash-royale/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	borderTop       = "╭───────────────────────────────────────────────────────╮"
	borderLeftRight = "│"
	borderMiddle    = "├───────────────────────────────────────────────────────┤"
	borderBottom    = "╰───────────────────────────────────────────────────────╯"
	tempContent     = "                                                     "
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
		case config.MsgTick:
			var current int
			json.Unmarshal(msg.Data.(json.RawMessage), &current)
			RenderDuration(current)
		case config.MsgUpdatePlayerMnana:
			var mana float64
			json.Unmarshal(msg.Data.(json.RawMessage), &mana)
			RenderPlayerMana(mana)
		case config.MsgUpdateOpponentMana:
			var mana float64
			json.Unmarshal(msg.Data.(json.RawMessage), &mana)
			RenderOpponentMana(mana)
		case config.MsgAttackResult:
			var combatResult game.CombatResult
			json.Unmarshal(msg.Data.(json.RawMessage), &combatResult)

			RenderNotification(combatString(combatResult)...)
		case config.MsgError:
			var err []string
			json.Unmarshal(msg.Data.(json.RawMessage), &err)
			RenderNotification(err...)
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
	fmt.Print("\033[38;1H")
	fmt.Print("\033[K")
	fmt.Print(">> ")
}

func RenderTemplate(matchData game.MatchData) {
	ClearScreen()

	// Timer
	fmt.Println(centerTitle("Timer", borderTop))
	fmt.Printf("%s\n", centerContent(durationString(0), borderTop))
	fmt.Println(borderBottom)

	// Player
	fmt.Println(centerTitle("You - Level", borderTop))
	fmt.Println(centerContent(fmt.Sprintf("%s - %d", matchData.PUsername, int(matchData.PLevel)), borderTop))
	fmt.Println(borderMiddle)
	fmt.Println(centerContent(manaString(matchData.PMana), borderTop))
	fmt.Println(borderMiddle)
	for i, troop := range matchData.PTroops {
		fmt.Println("│ " + troopString(i, troop) + "\t│")
	}
	fmt.Println(borderMiddle)
	for i, tower := range matchData.PTowers {
		fmt.Println("│ " + towerString(i, tower) + " \t\t│")
	}
	fmt.Println(borderBottom)

	// Opponent
	fmt.Println(centerTitle("Opponent - Level", borderTop))
	fmt.Println(centerContent(fmt.Sprintf("%s - %d", matchData.OUsername, int(matchData.OLevel)), borderTop))
	fmt.Println(borderMiddle)
	fmt.Println(centerContent(manaString(matchData.OMana), borderTop))
	fmt.Println(borderMiddle)
	for i, troop := range matchData.OTroops {
		fmt.Println("│ " + troopString(i, troop) + "\t│")
	}
	fmt.Println(borderMiddle)
	for i, tower := range matchData.OTowers {
		fmt.Println("│ " + towerString(i, tower) + " \t\t│")
	}
	fmt.Println(borderBottom)

	fmt.Println(centerTitle("Notification", borderTop))
	fmt.Println(centerContent(tempContent, borderTop))
	fmt.Println(centerContent(tempContent, borderTop))
	fmt.Println(centerContent(tempContent, borderTop))
	fmt.Println(borderBottom)

	fmt.Println(centerTitle("Command", borderTop))
	fmt.Println(centerContent("<troop_index> <tower_index>", borderTop))
	fmt.Println(borderBottom)
	fmt.Print(">> ")
}

func durationString(current int) string {
	elapsed := time.Duration(current) * time.Second

	remain := config.MatchDuration - elapsed

	minutes := int(remain.Minutes())
	seconds := int(remain.Seconds()) % 60

	return fmt.Sprintf("%d:%02d", minutes, seconds)
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
	str += fmt.Sprintf("\t🗡️ %d", int(troop.ATK))
	str += fmt.Sprintf("\t💧 %d", int(troop.Mana))

	return str
}

func combatString(combatResult game.CombatResult) []string {
	rs := []string{}

	str := combatResult.Attacker
	str += fmt.Sprintf(" ⚔️  %s", combatResult.Defender)
	rs = append(rs, str)

	str = fmt.Sprintf("🤖 %s ⚔️  %s 🏰", combatResult.UsingTroop.Name, combatResult.TargetTower.Type)
	rs = append(rs, str)

	str = fmt.Sprintf("🤖 (-%d🩸) ⚔️  🏰 (-%d🩸)", int(math.Ceil(combatResult.DamgeToTroop)), int(math.Ceil(combatResult.DamgeToTower)))
	rs = append(rs, str)

	return rs
}

func centerContent(content string, lineBase string) string {
	lineLength := utils.StringDisplayWidth(lineBase)

	innerWidth := lineLength - 4
	contentWidth := utils.StringDisplayWidth(content)
	if contentWidth >= innerWidth {
		return "│ " + content + " │"
	}

	padding := innerWidth - contentWidth

	leftPadding := padding / 2
	rightPadding := padding - leftPadding

	return "│ " + strings.Repeat(" ", leftPadding) + content + strings.Repeat(" ", rightPadding) + " │"
}

func centerTitle(content string, lineBase string) string {
	lineLength := utils.StringDisplayWidth(lineBase)

	innerWidth := lineLength - 4
	contentWidth := utils.StringDisplayWidth(content)
	if contentWidth >= innerWidth {
		return "╭ " + content + " ╮"
	}

	padding := innerWidth - contentWidth

	leftPadding := padding / 2
	rightPadding := padding - leftPadding

	return "╭" + strings.Repeat("─", leftPadding) + " " + content + " " + strings.Repeat("─", rightPadding) + "╮"
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
	str += fmt.Sprintf("\t🗡️ %d", int(tower.ATK))

	return str
}

func RenderDuration(current int) {
	fmt.Print("\033[s")    // Save pointer
	fmt.Print("\033[2;1H") // Move to line 2 col 1
	fmt.Print("\033[K")    // Clear line
	fmt.Printf("%s", centerContent(durationString(current), borderTop))
	fmt.Print("\033[u") // Back to previous
}

func RenderPlayerMana(mana float64) {
	fmt.Print("\033[s")                                   // Save pointer
	fmt.Print("\033[7;1H")                                // Move to line 7 col 1
	fmt.Print("\033[K")                                   // Clear line
	fmt.Print(centerContent(manaString(mana), borderTop)) // Print mana
	fmt.Print("\033[u")                                   // Back to previous
}

func RenderOpponentMana(mana float64) {
	fmt.Print("\033[s")                                   // Save pointer
	fmt.Print("\033[20;1H")                               // Move to line 20 col 1
	fmt.Print("\033[K")                                   // Clear line
	fmt.Print(centerContent(manaString(mana), borderTop)) // Print mana
	fmt.Print("\033[u")                                   // Back to previous
}

func manaString(mana float64) string {
	manaInt := int(mana)
	str := ""
	for range manaInt {
		str += "💧"
	}
	return str
}

func RenderNotification(content ...string) {
	fmt.Print("\033[s") // Save pointer

	switch len(content) {
	case 1:
		fmt.Print("\033[31;1H") // Move to line 30 col 1
		fmt.Print("\033[K")
		fmt.Print(centerContent(tempContent, borderTop)) // Print notification

		fmt.Print("\033[32;1H") // Move to line 30 col 1
		fmt.Print("\033[K")
		fmt.Print(centerContent(content[0], borderTop)) // Print notification

		fmt.Print("\033[33;1H") // Move to line 30 col 1
		fmt.Print("\033[K")
		fmt.Print(centerContent(tempContent, borderTop)) // Print notification
	case 2:
		fmt.Print("\033[31;1H") // Move to line 30 col 1
		fmt.Print("\033[K")
		fmt.Print(centerContent(content[0], borderTop)) // Print notification

		fmt.Print("\033[32;1H") // Move to line 30 col 1
		fmt.Print("\033[K")
		fmt.Print(centerContent(content[1], borderTop)) // Print notification

		fmt.Print("\033[33;1H") // Move to line 30 col 1
		fmt.Print("\033[K")
		fmt.Print(centerContent(tempContent, borderTop)) // Print notification
	case 3:
		fmt.Print("\033[31;1H") // Move to line 30 col 1
		fmt.Print("\033[K")
		fmt.Print(centerContent(content[0], borderTop)) // Print notification

		fmt.Print("\033[32;1H") // Move to line 30 col 1
		fmt.Print("\033[K")
		fmt.Print(centerContent(content[1], borderTop)) // Print notification

		fmt.Print("\033[33;1H") // Move to line 30 col 1
		fmt.Print("\033[K")
		fmt.Print(centerContent(content[2], borderTop)) // Print notification
	}
	fmt.Print("\033[u") // Back to previous
}
