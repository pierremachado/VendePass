package dao

import (
	"sync"
	"vendepass/internal/dao/interfaces"
	"vendepass/internal/models"

	"github.com/google/uuid"
)

var airportDao interfaces.AirportDAO
var flightDao interfaces.FlightDAO
var clientDao interfaces.ClientDAO
var sessionDao interfaces.SessionDAO

// GetFlightDAO returns a singleton instance of FlightDAO.
// If the instance does not exist, it creates a new one and initializes it.
//
// The FlightDAO interface provides methods for managing flight data.
// This implementation uses a MemoryFlightDAO, which stores flight data in memory.
//
// Parameters:
// None
//
// Return:
// flightDao (interfaces.FlightDAO) - A singleton instance of FlightDAO.
func GetFlightDAO() interfaces.FlightDAO {
	if flightDao == nil {
		flightDao = &MemoryFlightDAO{make(map[uuid.UUID]map[uuid.UUID]*models.Flight)}
		flightDao.New()
	}

	return flightDao
}

func GetClientDAO() interfaces.ClientDAO {
	if clientDao == nil {
		clientDao = &MemoryClientDAO{make(map[uuid.UUID]*models.Client)}
		clientDao.New()
	}

	return clientDao
}

func GetSessionDAO() interfaces.SessionDAO {
	if sessionDao == nil {
		sessionDao = &MemorySessionDAO{
			data: make(map[uuid.UUID]*models.Session),
			mu:   sync.RWMutex{}}
		sessionDao.New()
	}

	return sessionDao
}

func GetAirportDAO() interfaces.AirportDAO {
	if airportDao == nil {
		airportDao = &MemoryAirportDAO{make(map[uuid.UUID]models.Airport)}
		airportDao.New()
	}

	return airportDao
}
