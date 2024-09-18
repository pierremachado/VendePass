package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"
	"vendepass/internal/dao"
	"vendepass/internal/models"
)

func getClient(username string) (*models.Client, error) {
	client := findClient(username)

	var isError error = nil

	if client == nil {
		isError = errors.New("client not found")
	}

	return client, isError
}

func findClient(username string) *models.Client {
	for _, client := range dao.GetClientDAO().FindAll() {
		if client.Username == username {
			return client
		}
	}
	return nil
}

func passwordMatches(client *models.Client, password string) bool {
	return client.Password == password
}

func login(data interface{}, conn net.Conn) {
	var logCred models.LoginCredentials

	response := models.Response{Data: make(map[string]interface{})}

	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &logCred)

	login, err := getClient(logCred.Username)

	if err != nil {
		WriteNewResponse(
			models.Response{
				Error: err.Error(),
			}, conn)
		return
	}

	if passwordMatches(login, logCred.Password) {

		session := &models.Session{ClientID: login.Id, LastTimeActive: time.Now()}
		dao.GetSessionDAO().Insert(session)

		token := fmt.Sprintf("%s", session.ID)

		response.Data["token"] = token

	} else {
		response.Error = "invalid credentials"
	}
	WriteNewResponse(response, conn)

}

func logout(auth string, conn net.Conn) {
	defer conn.Close()
	response := models.Response{Data: make(map[string]interface{})}

	session, exists := SessionIfExists(auth)

	if !exists {
		response.Error = "session not found"
		WriteNewResponse(response, conn)
		return
	}

	dao.GetSessionDAO().Delete(*session)
	removeReservations(session)

	response.Data["msg"] = "logout succesfully made"
	WriteNewResponse(response, conn)
}

func removeReservations(session *models.Session) {
	for _, res := range session.Reservations {
		flight, _ := dao.GetFlightDAO().FindById(res.FlightId)
		flight.Seats++
	}
}

func getUserBySessionToken(auth string, conn net.Conn) {
	defer conn.Close()
	response := models.Response{Data: make(map[string]interface{})}

	session, exists := SessionIfExists(auth)

	if !exists {
		response.Error = "session not found"
		WriteNewResponse(response, conn)
		return
	}

	id := session.ClientID

	client, err := dao.GetClientDAO().FindById(id)

	if err != nil {
		response.Error = "client not found"
	}

	response.Data["user"] = client
	WriteNewResponse(response, conn)
}
