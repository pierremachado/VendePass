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

// getClient retrieves a client from the database using their username.
// If the client is found, it returns a pointer to the client and nil as the error.
// If the client is not found, it returns nil for the client and an error indicating that the client was not found.
//
// Parameters:
// - username: A string representing the username of the client to be retrieved.
//
// Return:
//   - A pointer to a models.Client representing the found client.
//     If no client with the given username is found, nil is returned.
//   - An error indicating whether the client was found or not.
//     Returns nil if the client is found, otherwise returns an error indicating that the client was not found.
func getClient(username string) (*models.Client, error) {
	client := findClient(username)

	var isError error = nil

	if client == nil {
		isError = errors.New("client not found")
	}

	return client, isError
}

// findClient searches for a client in the database using their username.
//
// Parameters:
// - username: A string representing the username of the client to be found.
//
// Return:
//   - A pointer to a models.Client representing the found client.
//     If no client with the given username is found, nil is returned.
func findClient(username string) *models.Client {
	for _, client := range dao.GetClientDAO().FindAll() {
		if client.Username == username {
			return client
		}
	}
	return nil
}

// passwordMatches checks if the provided password matches the password stored in the client's record.
//
// Parameters:
// - client: A pointer to a models.Client representing the user for whom the password needs to be checked.
// - password: A string representing the password provided by the user.
//
// Return:
//   - A boolean value indicating whether the provided password matches the password stored in the client's record.
//     Returns true if the passwords match, false otherwise.
func passwordMatches(client *models.Client, password string) bool {
	return client.Password == password
}

// login handles the login process for a user.
// It receives a data interface and a connection to the client.
// It first unmarshals the data into a LoginCredentials struct, then retrieves the client from the database using the provided username.
// If the client is not found, it sends an error response to the client and returns.
// If the client is found, it checks if the provided password matches the client's password.
// If the passwords match, it checks if the client is already logged in by searching for an active session.
// If the client is already logged in, it sends an error response to the client and returns.
// If the client is not logged in, it creates a new session for the client, stores it in the database, and sends a success response with the session token to the client.
// If the passwords do not match, it sends an error response to the client.
func login(data interface{}, conn net.Conn) {
	var logCred models.LoginCredentials

	response := models.Response{Data: make(map[string]interface{})}

	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &logCred)

	login, err := getClient(logCred.Username)
	fmt.Println(login)
	if err != nil {
		WriteNewResponse(
			models.Response{
				Error: err.Error(),
			}, conn)
		return
	}

	var session *models.Session

	if passwordMatches(login, logCred.Password) {

		if s := findUser(login); s != nil {
			WriteNewResponse(
				models.Response{
					Error: "more than one user logged",
				}, conn)
			return
		} else {
			session = &models.Session{ClientID: login.Id, LastTimeActive: time.Now()}
			dao.GetSessionDAO().Insert(session)
		}

		token := fmt.Sprintf("%s", session.ID)

		response.Data["token"] = token

	} else {
		response.Error = "invalid credentials"
	}
	WriteNewResponse(response, conn)
}

// findUser searches for an active session associated with a given client.
// It iterates through all sessions stored in the database and checks if the client's ID matches the session's ClientID.
// If a matching session is found, it is returned. Otherwise, nil is returned.
//
// Parameters:
// - login: A pointer to a models.Client representing the client for which the session needs to be found.
//
// Return:
//   - A pointer to a models.Session representing the active session associated with the given client.
//     If no matching session is found, nil is returned.
func findUser(login *models.Client) *models.Session {
	for _, s := range dao.GetSessionDAO().FindAll() {
		if s.ClientID == login.Id {
			return s
		}
	}
	return nil
}

// logout handles the logout process for a user.
// It closes the connection, prepares a response object, and checks if a session exists for the given authentication token.
// If the session is found, it deletes the session from the database and calls the removeReservations function to release reserved seats.
// Finally, it sends a success message in the response and writes it to the provided connection.
//
// Parameters:
// - auth: A string representing the authentication token.
// - conn: A net.Conn representing the connection to the client.
//
// Return:
// - None. The function writes the response directly to the connection.
func logout(auth string, conn net.Conn) {
	defer conn.Close()
	response := models.Response{Data: make(map[string]interface{})}

	session, exists := SessionIfExists(auth)

	if !exists {
		response.Error = "session not found"
		WriteNewResponse(response, conn)
		return
	}

	dao.GetSessionDAO().Delete(session)
	removeReservations(session)

	response.Data["msg"] = "logout successfully made"
	WriteNewResponse(response, conn)
}

// removeReservations decreases the number of available seats for each flight associated with a given session.
// It iterates through the reservations associated with the session, retrieves the corresponding flight from the database,
// and increments the number of available seats by one.
//
// Parameters:
// - session: A pointer to a models.Session representing the session for which reservations need to be processed.
//
// Return:
// - None. The function modifies the flight records directly in the database.
func removeReservations(session *models.Session) {
	for _, res := range session.Reservations {
		flight, _ := dao.GetFlightDAO().FindById(res.FlightId)
		flight.Seats++
	}
}

// getUserBySessionToken retrieves the user associated with a given session token.
// It uses the provided authentication token to find the corresponding session in the database.
// If the session is found, it retrieves the user's ID from the session and uses it to fetch the user from the database.
// If the session or user is not found, appropriate error messages are set in the response.
// The function then writes the response to the provided connection.
//
// Parameters:
// - auth: A string representing the authentication token.
// - conn: A net.Conn representing the connection to the client.
//
// Return:
// - None. The function writes the response directly to the connection.
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
