package server

import (
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
