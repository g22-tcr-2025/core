package network

import (
	"fmt"
	"log"
	"net"
)

type MatchHandler interface {
	HandleConnection(conn net.Conn)
}

func StartServer(port int, handler MatchHandler) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Printf("Server listening on port %d\n", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}
		go handler.HandleConnection(conn)
	}
}
