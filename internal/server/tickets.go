package server

import (
	"encoding/json"
	"net"
	"vendepass/internal/dao"
	"vendepass/internal/models"

	"github.com/google/uuid"
)

func GetTickets(auth string, conn net.Conn) {
	session, exists := SessionIfExists(auth)

	if !exists {
		WriteNewResponse(models.Response{
			Error: "not authorized",
		}, conn)
		return
	}
	responseData := make([]map[string]interface{}, 0)

	client, _ := dao.GetClientDAO().FindById(session.ClientID)

	for _, ticket := range client.Client_flights {
		flight, _ := dao.GetFlightDAO().FindById(ticket.FlightId)

		src, _ := dao.GetAirportDAO().FindById(flight.SourceAirportId)
		dest, _ := dao.GetAirportDAO().FindById(flight.DestAirportId)

		flightresponse := make(map[string]interface{})

		flightresponse["Src"] = src.City
		flightresponse["Dest"] = dest.City
		flightresponse["Id"] = ticket.Id
		responseData = append(responseData, flightresponse)
	}

	WriteNewResponse(models.Response{
		Data: map[string]interface{}{
			"Tickets": responseData,
		},
	}, conn)
}

func BuyTicket(auth string, data interface{}, conn net.Conn) {
	session, exists := SessionIfExists(auth)

	if !exists {
		WriteNewResponse(models.Response{
			Error: "not authorized",
		}, conn)
		return
	}

	var buyTicket models.BuyTicket

	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &buyTicket)

	res, exists := session.Reservations[buyTicket.ReservationId]

	if !exists {
		WriteNewResponse(models.Response{
			Error: "reservation do not exists",
		}, conn)
	} else {
		client, _ := dao.GetClientDAO().FindById(res.ClientId)
		flight, _ := dao.GetFlightDAO().FindById(res.FlightId)

		flight.Mu.Lock()

		flight.Passengers = append(flight.Passengers, res.Ticket)

		flight.Mu.Unlock()

		client.Client_flights = append(client.Client_flights, res.Ticket)

		delete(session.Reservations, res.Id)
	}

	WriteNewResponse(models.Response{
		Data: map[string]interface{}{
			"msg": "success",
		},
	}, conn)
}

func CancelBuy(auth string, data interface{}, conn net.Conn) {
	session, exists := SessionIfExists(auth)

	if !exists {
		WriteNewResponse(models.Response{
			Error: "not authorized",
		}, conn)
		return
	}

	var cancelReservation models.CancelBuyRequest

	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &cancelReservation)

	client, _ := dao.GetClientDAO().FindById(session.ClientID)
	ticket := findTicketById(client.Client_flights, cancelReservation.TicketId)
	client.Client_flights = removeTicketByID(client.Client_flights, ticket.Id)

	flight, _ := dao.GetFlightDAO().FindById(ticket.FlightId)

	flight.Mu.Lock()
	flight.Seats++
	flight.Passengers = removeTicketByID(flight.Passengers, ticket.Id)
	flight.Mu.Unlock()

	WriteNewResponse(models.Response{
		Data: map[string]interface{}{
			"msg": "success",
		},
	}, conn)

}

func findTicketById(tickets []*models.Ticket, id uuid.UUID) *models.Ticket {
	for _, ticket := range tickets {
		if ticket.Id == id {
			return ticket
		}
	}
	return nil
}

func removeTicketByID(tickets []*models.Ticket, id uuid.UUID) []*models.Ticket {
	for i, ticket := range tickets {
		if ticket.Id == id {
			return append(tickets[:i], tickets[i+1:]...)
		}
	}
	return tickets
}
