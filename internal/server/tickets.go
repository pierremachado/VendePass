package server

import (
	"net"
	"vendepass/internal/dao"
	"vendepass/internal/models"
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
