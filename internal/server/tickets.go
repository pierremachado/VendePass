package server

import (
	"encoding/json"
	"net"
	"vendepass/internal/dao"
	"vendepass/internal/models"

	"github.com/google/uuid"
)

// GetTickets retrieves all tickets associated with the authenticated client.
// It sends a response containing a list of tickets with their respective source, destination, and ID.
//
// Parameters:
// - auth: A string representing the authentication token.
// - conn: A net.Conn object representing the connection to the client.
//
// Return:
// - No return value.
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

// BuyTicket handles the process of purchasing a ticket for an authenticated client.
// It checks if the client is authorized, validates the reservation, updates the flight and client data,
// and sends a response indicating success or failure.
//
// Parameters:
// - auth: A string representing the authentication token.
// - data: An interface containing the necessary data for purchasing a ticket.
// - conn: A net.Conn object representing the connection to the client.
//
// Return:
// - No return value.
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

// CancelBuy handles the cancellation of a ticket for an authenticated client.
// It checks if the client is authorized, finds the ticket to be canceled, updates the flight and client data,
// and sends a response indicating success or failure.
//
// Parameters:
// - auth: A string representing the authentication token. This is used to identify the client.
// - data: An interface containing the necessary data for canceling a ticket.
// - conn: A net.Conn object representing the connection to the client. This is used to send a response.
//
// Return:
// - No return value.
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

// findTicketById searches for a ticket with the given ID in a list of tickets.
//
// Parameters:
// - tickets: A slice of pointers to models.Ticket, representing the list of tickets to search.
// - id: A uuid.UUID representing the ID of the ticket to find.
//
// Return:
// - A pointer to models.Ticket if a ticket with the given ID is found in the list.
// - nil if no ticket with the given ID is found in the list.
func findTicketById(tickets []*models.Ticket, id uuid.UUID) *models.Ticket {
	for _, ticket := range tickets {
		if ticket.Id == id {
			return ticket
		}
	}
	return nil
}

// removeTicketByID searches for and removes a ticket with the given ID from a list of tickets.
//
// Parameters:
// - tickets: A slice of pointers to models.Ticket, representing the list of tickets to search and remove from.
// - id: A uuid.UUID representing the ID of the ticket to find and remove.
//
// Return:
// - A slice of pointers to models.Ticket representing the updated list of tickets after removing the ticket with the given ID.
//   If no ticket with the given ID is found, the original list is returned.
func removeTicketByID(tickets []*models.Ticket, id uuid.UUID) []*models.Ticket {
	for i, ticket := range tickets {
		if ticket.Id == id {
			return append(tickets[:i], tickets[i+1:]...)
		}
	}
	return tickets
}
