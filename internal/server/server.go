package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Listener não criado", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("servidor ouvindo na porta :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("erro ao aceitar conexão", err)
			continue
		}

		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)

	if err != nil {
		fmt.Println("erro na leitura do buffer", err)
		return
	}

	fmt.Println("recebido: ", string(buffer[:n]))

	_, err = conn.Write([]byte("mensagem chegou no server\n"))

	if err != nil {
		fmt.Println("erro na comunicação com o cliente", err)
		return
	}

}
