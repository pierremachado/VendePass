package main

import (
	"encoding/json"
	"fmt"
	"net"
	"vendepass/internal/models"

	"github.com/google/uuid"
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

	_, write_err := conn.Write(buffer)

	if write_err != nil {
		fmt.Println("erro na escrita: ", write_err)
		return
	}

	receive := make([]byte, 2048)

	n, read_err := conn.Read(receive)

	if read_err != nil {
		fmt.Println("erro na escrita: ", read_err)
		return
	}

	token, _ := uuid.Parse(string(receive[:n]))

	fmt.Println(token)

	logout := models.Request{
		Action: "logout",
		Data: models.LogoutCredentials{
			TokenId: token,
		},
	}

	buffer, _ = json.Marshal(logout)

	_, write_err = conn.Write(buffer)

	if write_err != nil {
		fmt.Println("erro na escrita: ", write_err)
		return
	}

	receive = make([]byte, 2048)

	n, read_err = conn.Read(receive)

	if read_err != nil {
		fmt.Println("erro na escrita: ", read_err)
		return
	}
	fmt.Println(string(receive[:n]))

}
