package network

import (
	"log"
	"net"
)

func ConnectToServer(addr string) net.Conn {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Could not connect to server", err)
	}

	return conn
}
