package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"vendepass/internal/dao"
	"vendepass/internal/server"
)

const (
	port      = ":8080"
	timeLimit = 30 * time.Minute
)

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Listener não criado", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("servidor ouvindo na porta :8080")

	go server.CleanupSessions(timeLimit)

	for _, flight := range dao.GetFlightDAO().FindAll() {
		go flight.ProcessReservations()
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("erro ao aceitar conexão", err)
			continue
		}

		go server.HandleConn(conn)
	}

}
