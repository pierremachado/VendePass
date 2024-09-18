package server

import (
	"encoding/json"
	"net"
	"vendepass/internal/dao"
	"vendepass/internal/models"
)

func GetCart(auth string, conn net.Conn) {
	session, exists := SessionIfExists(auth)

	if !exists {
		WriteNewResponse(models.Response{
			Error: "not authorized",
		}, conn)
		return
	}
	responseData := make([]map[string]interface{}, 0)

	for _, reservation := range session.Reservations {
		flight, _ := dao.GetFlightDAO().FindById(reservation.FlightId)

		src, _ := dao.GetAirportDAO().FindById(flight.SourceAirportId)
		dest, _ := dao.GetAirportDAO().FindById(flight.DestAirportId)

		flightresponse := make(map[string]interface{})

		flightresponse["Src"] = src.City
		flightresponse["Dest"] = dest.City
		flightresponse["Id"] = reservation.Id
		responseData = append(responseData, flightresponse)
	}

	WriteNewResponse(models.Response{
		Data: map[string]interface{}{
			"Reservations": responseData,
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
