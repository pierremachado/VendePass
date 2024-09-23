package server

import (
	"encoding/json"
	"fmt"
	"net"
	"vendepass/internal/dao"
	"vendepass/internal/models"
)

// Reservation handles the creation of reservations for a given set of flights.
// It verifies the session's existence, deserializes the request data, processes each requested flight,
// checks for availability, and sends the session to the respective flight's reservation queue.
// If any flight is not available, it responds with an error.
//
// Parameters:
// 	- auth: A string representing the session's authentication token.
// 	- data: An interface containing the request data. It should be of type models.FlightsRequest.
// 	- conn: A net.Conn representing the connection to the client.
func Reservation(auth string, data interface{}, conn net.Conn) {
	// Check if the session exists
	session, exists := SessionIfExists(auth)

	if !exists {
		WriteNewResponse(models.Response{
			Error: "not authorized",
		}, conn)
		return
	}

	var flightRequest models.FlightsRequest

	// Deserialize the request data
	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &flightRequest)

	flights := make([]*models.Flight, len(flightRequest.FlightIds))
	var notAvailableFlights []*models.Flight
	// Process each requested flight
	for i, id := range flightRequest.FlightIds {
		flight, _ := dao.GetFlightDAO().FindById(id)
		flight.Mu.Lock()
		if flight.Seats <= 0 {
			notAvailableFlights = append(notAvailableFlights, flight)
		} else {
			flights[i] = flight
		}
		flight.Mu.Unlock()
	}

	// Check if any flight is not available and respond with error
	if len(notAvailableFlights) > 0 {
		responseData, _ := getRoute(flightRequest.FlightIds)
		WriteNewResponse(models.Response{
			Error: "at least one flight is not available",
			Data: map[string]interface{}{
				"Flights": responseData,
			},
		}, conn)
		return
	}

	for _, flight := range flights {
		// Send the session to the flight's reservation queue
		flight.Queue <- session
	}

	// Success: Reservations created successfully
	WriteNewResponse(models.Response{
		Data: map[string]interface{}{
			"msg": "success",
		},
	}, conn)
}

// CancelReservation cancels a reservation for a specific flight.
// It verifies the session's existence, deserializes the request data, retrieves the reservation,
// releases the seat on the flight, and removes the reservation from the session.
//
// Parameters:
// 	- auth: A string representing the session's authentication token.
// 	- data: An interface containing the request data. It should be of type models.CancelReservationRequest.
// 	- conn: A net.Conn representing the connection to the client.
func CancelReservation(auth string, data interface{}, conn net.Conn) {
	// Verify if the session exists
	session, exists := SessionIfExists(auth)

	if !exists {
		WriteNewResponse(models.Response{
			Error: "not authorized",
		}, conn)
		return
	}

	var cancelReservation models.CancelReservationRequest

	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &cancelReservation)

	reservation := session.Reservations[cancelReservation.ReservationId]

	flight, _ := dao.GetFlightDAO().FindById(reservation.FlightId)

	flight.Mu.Lock()
	flight.Seats++
	flight.Mu.Unlock()

	delete(session.Reservations, cancelReservation.ReservationId)

	fmt.Printf("Session %s: flight %s canceled successfully\n", session.ID, flight.Id)

	WriteNewResponse(models.Response{
		Data: map[string]interface{}{
			"msg": "success",
		},
	}, conn)
}
