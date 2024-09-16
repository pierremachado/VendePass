package tests

import (
	"fmt"
	"testing"
	"vendepass/internal/dao"
	"vendepass/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInsertAndRetrieveFlight(t *testing.T) {
	flightDAO := dao.GetFlightDAO()
	defer flightDAO.DeleteAll()

	flight := &models.Flight{
		Id:              uuid.New(),
		SourceAirportId: uuid.New(),
		DestAirportId:   uuid.New(),
		Seats:           10,
	}

	flightDAO.Insert(flight)

	retrievedFlight, err := flightDAO.FindById(flight.Id)

	assert.NoError(t, err, "expected no error, got %v", err)
	assert.Equal(t, flight, retrievedFlight, "flights should be the same")
}

func TestFindAllFlights(t *testing.T) {
	flightDAO := dao.GetFlightDAO()
	defer flightDAO.DeleteAll()

	flight1 := &models.Flight{
		Id:              uuid.New(),
		SourceAirportId: uuid.New(),
		DestAirportId:   uuid.New(),
		Seats:           10,
	}

	flight2 := &models.Flight{
		Id:              uuid.New(),
		SourceAirportId: uuid.New(),
		DestAirportId:   uuid.New(),
		Seats:           5,
	}

	flightDAO.Insert(flight1)
	flightDAO.Insert(flight2)

	flights := flightDAO.FindAll()

	for _, flight := range flights {
		fmt.Println(flight.Seats)
	}

	assert.Equal(t, 2, len(flights), "expected 2 flights, got: %d", len(flights))
	assert.Contains(t, flights, flight1, "expected flight1 to be in the list")
	assert.Contains(t, flights, flight2, "expected flight2 to be in the list")
	assert.NotContains(t, flights, &models.Flight{Id: uuid.New()}, "expected flight3 to not be in the list")
}

func TestUpdateFlight(t *testing.T) {
	flightDAO := dao.GetFlightDAO()
	defer flightDAO.DeleteAll()

	flight := &models.Flight{
		Id:              uuid.New(),
		SourceAirportId: uuid.New(),
		DestAirportId:   uuid.New(),
		Seats:           10,
	}

	flightDAO.Insert(flight)

	flight.Seats = 5

	flightDAO.Update(flight)

	retrievedFlight, err := flightDAO.FindById(flight.Id)

	assert.NoError(t, err, "expected no error, got: %v", retrievedFlight)
	assert.Equal(t, uint(5), retrievedFlight.Seats, "seats were not updated")

	flight = &models.Flight{SourceAirportId: uuid.New(), DestAirportId: uuid.New(), Seats: 0}
	err = flightDAO.Update(flight)

	assert.Error(t, err, "expected error upon trying to update an nonexistent flight")
}

func TestDeleteFlight(t *testing.T) {
	flightDAO := dao.GetFlightDAO()
	defer flightDAO.DeleteAll()

	flight := &models.Flight{
		Id:              uuid.New(),
		SourceAirportId: uuid.New(),
		DestAirportId:   uuid.New(),
		Seats:           10,
	}

	flightDAO.Insert(flight)

	err := flightDAO.Delete(*flight)

	assert.NoError(t, err, "expected no error, got: %v", err)

	_, err = flightDAO.FindById(flight.Id)

	assert.Error(t, err, "expected error upon trying to find deleted flight, got %v", err)
}

func TestFindBySourceAndDest(t *testing.T) {
	flightDAO := dao.GetFlightDAO()
	defer flightDAO.DeleteAll()

	sourceID := uuid.New()
	destID := uuid.New()

	flightDAO.Insert(&models.Flight{SourceAirportId: sourceID, DestAirportId: destID, Seats: 10})

	flight, err := flightDAO.FindBySourceAndDest(sourceID, destID)

	assert.NoError(t, err, "expected no error, got %v", err)
	assert.Equal(t, sourceID, flight.SourceAirportId)
	assert.Equal(t, destID, flight.DestAirportId)
}

func TestFindBySource(t *testing.T) {
	flightDAO := dao.GetFlightDAO()
	defer flightDAO.DeleteAll()

	sourceID := uuid.New()

	flightDAO.Insert(&models.Flight{SourceAirportId: sourceID, DestAirportId: uuid.New(), Seats: 10})
	flightDAO.Insert(&models.Flight{SourceAirportId: sourceID, DestAirportId: uuid.New(), Seats: 5})

	flights, err := flightDAO.FindBySource(sourceID)

	assert.NoError(t, err, "expected no error, got %v", err)
	assert.Equal(t, 2, len(flights), "expected 2 flights, got %d", len(flights))

	for _, flight := range flights {
		assert.Equal(t, sourceID, flight.SourceAirportId)
	}
}

func TestFindBySourceAndDestWithInvalidIDs(t *testing.T) {
	flightDAO := dao.GetFlightDAO()
	defer flightDAO.DeleteAll()

	_, err := flightDAO.FindBySourceAndDest(uuid.Nil, uuid.Nil)
	assert.Error(t, err, "expected error, got %v", err)

	_, err = flightDAO.FindBySource(uuid.Nil)
	assert.Error(t, err, "expected error, got %v", err)
}

func TestBFS(t *testing.T) {
	flightDAO := dao.GetFlightDAO()
	defer flightDAO.DeleteAll()

	sourceID := uuid.New()
	destID := uuid.New()
	middleID := uuid.New()

	flightDAO.Insert(&models.Flight{SourceAirportId: sourceID, DestAirportId: middleID, Seats: 10})
	flightDAO.Insert(&models.Flight{SourceAirportId: middleID, DestAirportId: destID, Seats: 5})

	path, err := flightDAO.BreadthFirstSearch(sourceID, destID)

	assert.NoError(t, err, "expected no errors, got %v", err)
	assert.NotNil(t, path, "expected path, got %v", path)

	expectedPath := []*models.Flight{
		{SourceAirportId: sourceID, DestAirportId: middleID},
		{SourceAirportId: middleID, DestAirportId: destID},
	}
	assert.Equal(t, len(expectedPath), len(path), "seats should be equal")
	for i := range expectedPath {
		assert.Equal(t, expectedPath[i].SourceAirportId, path[i].SourceAirportId)
		assert.Equal(t, expectedPath[i].DestAirportId, path[i].DestAirportId)
	}

	unreachableID := uuid.New()
	path, err = flightDAO.BreadthFirstSearch(sourceID, unreachableID)

	assert.Error(t, err, "expected error, got %v", err)
	assert.Nil(t, path, "expected path to be nil, got %v")
}
