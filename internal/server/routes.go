package server

import (
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

	WriteNewResponse(models.Response{
		Data: map[string]interface{}{
			"all-routes": dao.FindAll(),
		},
	}, conn)

}
