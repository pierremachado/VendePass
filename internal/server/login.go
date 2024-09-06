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

func login(data interface{}, conn net.Conn, session *models.Session) {
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
		token := fmt.Sprintf("%s", session.ID)
		_, err = conn.Write([]byte(token))

		if err != nil {
			fmt.Println("erro na comunicação com o cliente", err)
			return
		}
	}

}

func logout(data interface{}, conn net.Conn) {
	defer conn.Close()
	var logCred models.LogoutCredentials

	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &logCred)

	session, err := dao.GetSessionDAO().FindById(logCred.TokenId)

	if err != nil {
		conn.Write([]byte("erro na remoção de sessão"))
	}

	dao.GetSessionDAO().Delete(session)
	conn.Write([]byte("logout realizado com sucesso"))
}
