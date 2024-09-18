package server

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
	"vendepass/internal/dao"
	"vendepass/internal/models"

	"github.com/google/uuid"
)

func HandleConn(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)

	if err != nil {
		WriteNewResponse(models.Response{
			Error: "error when reading the buffer",
		}, conn)
		return
	}

	if n > 0 {
		var request models.Request
		err = json.Unmarshal(buffer[:n], &request)
		if err != nil {
			WriteNewResponse(models.Response{
				Error: "error on request format",
			}, conn)
			return
		}

		handleRequest(request, conn)
	}

}

func CleanupSessions(timeout time.Duration) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		for _, session := range dao.GetSessionDAO().FindAll() {
			if time.Since(session.LastTimeActive) > timeout {
				fmt.Printf("Encerrando sessão %s por inatividade\n", session.ID)
				dao.GetSessionDAO().Delete(*session)
			}
			for key, reservation := range session.Reservations {
				if time.Since(reservation.CreatedAt) > timeout {
					fmt.Printf("Encerrando reserva %s por inatividade\n", reservation.Id)
					flight, _ := dao.GetFlightDAO().FindById(reservation.FlightId)
					flight.Seats++
					delete(session.Reservations, key)
				}
			}
		}
	}
}

func handleRequest(request models.Request, conn net.Conn) {
	switch request.Action {
	case "login":
		login(request.Data, conn)
	case "get-user":
		getUserBySessionToken(request.Auth, conn)
	case "logout":
		logout(request.Auth, conn)
	case "all-routes":
		AllRoutes(request.Auth, conn)
	case "route":
		Route(request.Auth, request.Data, conn)
	case "flights":
		Flights(request.Auth, request.Data, conn)
	case "reservation":
		Reservation(request.Auth, request.Data, conn)
	case "cart":
		GetCart(request.Auth, conn)
	case "buy":
		BuyTicket(request.Auth, request.Data, conn)
	case "tickets":
		GetTickets(request.Auth, conn)
	}

}

func WriteNewResponse(response models.Response, conn net.Conn) {
	if response.Data == nil {
		response.Data = make(map[string]interface{})
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshalling response:", err)
		return
	}
	_, err = conn.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing response:", err)
	}
}

func SessionIfExists(token string) (*models.Session, bool) {

	uuid, err := uuid.Parse(token)
	if err != nil {
		return nil, false
	}

	session, err := dao.GetSessionDAO().FindById(uuid)
	if err != nil {
		return nil, false
	}
	session.LastTimeActive = time.Now()
	dao.GetSessionDAO().Update(session)
	return session, true
}
