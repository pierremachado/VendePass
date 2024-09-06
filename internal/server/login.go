package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"vendepass/internal/dao"
	"vendepass/internal/models"
	"vendepass/internal/utils"
)

func GetClient(username string) (*models.Client, error) {
	clientDao := dao.GetClientDAO()
	client := utils.Find[models.Client](clientDao.FindAll(), func(c models.Client) bool {
		return c.Username == username
	})

	var isError error = nil

	if client == nil {
		isError = errors.New("client not found")
	}

	return client, isError
}

func PasswordMatches(client *models.Client, password string) bool {
	return client.Password == password
}

func login(data interface{}, conn net.Conn) {
	var logCred models.LoginCredentials

	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &logCred)

	login, err := GetClient(logCred.Username)

	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(PasswordMatches(login, logCred.Password))
	if PasswordMatches(login, logCred.Password) {
		welcome := fmt.Sprintf("Bem vindo, %s", login.Name)
		_, err = conn.Write([]byte(welcome))

		if err != nil {
			fmt.Println("erro na comunicação com o cliente", err)
			return
		}

	}

}
