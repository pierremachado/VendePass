package server

import (
	"encoding/json"
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

	source, errsource := uuid.Parse(routeRequest.Source)
	dest, errdest := uuid.Parse(routeRequest.Dest)

	if errsource != nil || errdest != nil {
		WriteNewResponse(models.Response{
			Error: "not uuid",
		}, conn)
		return
	}

	route, err := dao.GetFlightDAO().BreadthFirstSearch(source, dest)
	if err != nil {
		response.Error = "no route"
	} else {
		response.Data = map[string]interface{}{
			"route": route,
		}
	}

	WriteNewResponse(response, conn)

}
