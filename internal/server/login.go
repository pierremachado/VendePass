package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"
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

	response := models.Response{Data: make(map[string]interface{})}

	// defer WriteNewResponse(response, conn)

	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &logCred)

	login, err := GetClient(logCred.Username)

	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(PasswordMatches(login, logCred.Password))

	if PasswordMatches(login, logCred.Password) {
		session := &models.Session{Client: *login, LastTimeActive: time.Now()}
		dao.GetSessionDAO().Insert(session)

		token := fmt.Sprintf("%s", session.ID)

		response.Data["token"] = token

	} else {
		response.Error = "invalid credentials"
	}
	WriteNewResponse(response, conn)
}

func logout(data interface{}, conn net.Conn) {
	defer conn.Close()
	var logCred models.LogoutCredentials

	response := models.Response{Data: make(map[string]interface{})}

	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &logCred)

	session, err := dao.GetSessionDAO().FindById(logCred.TokenId)

	if err != nil {
		response.Error = "session not found"
		WriteNewResponse(response, conn)
		return
	}

	dao.GetSessionDAO().Delete(session)

	response.Data["msg"] = "logout realizado com sucesso"
	WriteNewResponse(response, conn)
}
