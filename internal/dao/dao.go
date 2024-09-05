package dao

import (
	"vendepass/internal/dao/interfaces"
	"vendepass/internal/models"

	"github.com/google/uuid"
)

var ticketDao interfaces.TicketDAO
var flightDao interfaces.FlightDAO

func GetTicketDAO() interfaces.TicketDAO {

	if ticketDao == nil {
		ticketDao = &MemoryTicketDAO{make(map[uuid.UUID]models.Ticket)}
	}

	return ticketDao
}

func GetFlightDAO() interfaces.FlightDAO {
	if flightDao == nil {
		flightDao = &MemoryFlightDAO{make(map[uuid.UUID]models.Flight)}
	}

	return flightDao
}

var clientDao interfaces.ClientDAO

func GetClientDAO() interfaces.ClientDAO {
	if clientDao == nil {
		clientDao = &MemoryClientDAO{make(map[uuid.UUID]models.Client)}
	}

	return clientDao
}
