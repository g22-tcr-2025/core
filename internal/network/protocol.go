package network

import (
	"bufio"
	"encoding/json"
	"net"
)

type Message struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func SendMessage(conn net.Conn, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// Adding the delim
	data = append(data, byte('\n'))

	// Sending message
	_, err = conn.Write(data)
	return err
}

func ReceiveMessage(reader *bufio.Reader) (Message, error) {
	line, err := reader.ReadBytes(byte('\n'))
	if err != nil {
		return Message{}, err
	}

	var rawMsg map[string]json.RawMessage
	err = json.Unmarshal(line, &rawMsg)
	if err != nil {
		return Message{}, err
	}

	var msg Message
	err = json.Unmarshal(rawMsg["type"], &msg.Type)
	if err != nil {
		return Message{}, err
	}

	msg.Data = rawMsg["data"]
	return msg, nil
}
