package main

import (
	"encoding/json"
	"fmt"
	"net"
	"vendepass/internal/models"
)

const (
	address = "localhost:8080"
)

func main() {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("erro na conexão: ", err)
		return
	}
	defer conn.Close()

	login := models.Request{
		Action: "login",
		Data: models.LoginCredentials{
			Username: "pedrocosta",
			Password: "senhaSegura789",
		},
	}

	buffer, _ := json.Marshal(login)

	n, write_err := conn.Write(buffer)

	if write_err != nil {
		fmt.Println("erro na escrita: ", err)
	}

	fmt.Println(n)
}
