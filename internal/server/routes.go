package server

import (
	"encoding/json"
	"fmt"
	"net"
	"vendepass/internal/dao"
	"vendepass/internal/models"

	"github.com/google/uuid"
)

func AllRoutes(auth string, conn net.Conn) {

	_, exists := SessionIfExists(auth)

	if !exists {
		WriteNewResponse(models.Response{
			Error: "not authorized",
		}, conn)
		return
	}

	dao := dao.GetFlightDAO()
	dao.New()

	WriteNewResponse(models.Response{
		Data: map[string]interface{}{
			"all-routes": dao.FindAll(),
		},
	}, conn)

}

func Route(auth string, data interface{}, conn net.Conn) {
	_, exists := SessionIfExists(auth)

	if !exists {
		WriteNewResponse(models.Response{
			Error: "not authorized",
		}, conn)
		return
	}

	var routeRequest models.RouteRequest
	var response models.Response

	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &routeRequest)

	src := dao.GetAirportDAO().FindByName(routeRequest.Source)
	dest := dao.GetAirportDAO().FindByName(routeRequest.Dest)

	if src == nil || dest == nil {
		WriteNewResponse(models.Response{
			Error: "not valid city name",
		}, conn)
		return
	}

	// route, err := dao.GetFlightDAO().FindBySourceAndDest(src.Id, dest.Id)
	path, path_err := dao.GetFlightDAO().BreadthFirstSearch(src.Id, dest.Id)
	if path_err != nil {
		response.Error = "no route"
	} else {
		cities_path := make([]models.Route, len(path))
		for i, flight := range path {
			cities_path[i].Path = make([]models.City, 2)
			cities_path[i].FlightId = flight.Id
			srcById, _ := dao.GetAirportDAO().FindById(flight.SourceAirportId)
			cities_path[i].Path[0] = srcById.City
			destById, _ := dao.GetAirportDAO().FindById(flight.DestAirportId)
			cities_path[i].Path[1] = destById.City
		}

		response.Data = map[string]interface{}{
			"path": cities_path,
		}

	}

	WriteNewResponse(response, conn)

}

func Flights(auth string, data interface{}, conn net.Conn) {
	_, exists := SessionIfExists(auth)

	if !exists {
		WriteNewResponse(models.Response{
			Error: "not authorized",
		}, conn)
		return
	}

	var flightsRequest models.FlightsRequest

	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &flightsRequest)

	responseData, err := getRoute(flightsRequest.FlightIds)
	if err != nil {
		WriteNewResponse(models.Response{
			Error: err.Error(),
		}, conn)
	}

	WriteNewResponse(models.Response{
		Data: map[string]interface{}{
			"Flights": responseData,
		},
	}, conn)
}

func getRoute(flightIds []uuid.UUID) ([]map[string]interface{}, error) {
	responseData := make([]map[string]interface{}, len(flightIds))
	for i, id := range flightIds {
		flightresponse := make(map[string]interface{})
		flight, err := dao.GetFlightDAO().FindById(id)
		if err != nil {

			return nil, fmt.Errorf("some flight doesnt exists: %s", id)
		}

		src, _ := dao.GetAirportDAO().FindById(flight.SourceAirportId)
		dest, _ := dao.GetAirportDAO().FindById(flight.DestAirportId)

		flightresponse["Seats"] = flight.Seats
		flightresponse["Src"] = src.City.Name
		flightresponse["Dest"] = dest.City.Name
		responseData[i] = flightresponse
	}
	return responseData, nil
}
