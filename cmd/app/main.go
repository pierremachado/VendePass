package main

import (
	"fmt"
	"net"
	"os"
	"vendepass/internal/models"
	"vendepass/internal/server"

	"github.com/google/uuid"
)

const (
	port string = ":8080"
)

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Listener não criado", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("servidor ouvindo na porta :8080")

	sessions := make(map[uuid.UUID]*models.Session)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("erro ao aceitar conexão", err)
			continue
		}

		go server.HandleConn(conn, sessions)
	}

}
