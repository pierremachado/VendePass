package server

import (
	"encoding/json"
	"net"
	"vendepass/internal/dao"
	"vendepass/internal/models"
)

func Reservation(auth string, data interface{}, conn net.Conn) {
	// Verifica se a sessão existe
	session, exists := SessionIfExists(auth)

	if !exists {
		WriteNewResponse(models.Response{
			Error: "not authorized",
		}, conn)
		return
	}

	var flightRequest models.FlightsRequest

	// Deserializa os dados da requisição
	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &flightRequest)

	flights := make([]*models.Flight, len(flightRequest.FlightIds))
	var notAvailableFlights []*models.Flight
	// Processa cada voo solicitado
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

	// Verifica se houve algum voo indisponível e responde com erro
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
		// Envia a sessão para a fila de reserva
		flight.Queue <- session
	}

	// Sucesso: Reservas criadas com sucesso
	WriteNewResponse(models.Response{
		Data: map[string]interface{}{
			"msg": "success",
		},
	}, conn)
}
