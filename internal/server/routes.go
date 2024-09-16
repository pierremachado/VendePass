package server

import (
	"encoding/json"
	"fmt"
	"net"
	"vendepass/internal/dao"
	"vendepass/internal/models"
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

	fmt.Println(routeRequest.Source, routeRequest.Dest)

	src := dao.GetAirportDAO().FindByName(routeRequest.Source)
	dest := dao.GetAirportDAO().FindByName(routeRequest.Dest)

	if src == nil || dest == nil {
		WriteNewResponse(models.Response{
			Error: "not valid city name",
		}, conn)
		return
	}

	route, err := dao.GetFlightDAO().FindBySourceAndDest(src.Id, dest.Id)
	if err != nil {
		response.Error = "no route"
	} else {
		response.Data = map[string]interface{}{
			"route": route,
		}
	}

	WriteNewResponse(response, conn)

}
