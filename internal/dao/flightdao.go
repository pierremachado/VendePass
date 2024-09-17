package dao

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"vendepass/internal/models"

	"github.com/google/uuid"
)

type MemoryFlightDAO struct {
	data map[uuid.UUID]map[uuid.UUID]*models.Flight
}

func (dao *MemoryFlightDAO) New() {

	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	jsonPath := filepath.Join(baseDir, "internal", "stubs", "flights.json")

	b, _ := os.ReadFile(jsonPath)

	var flights []models.FlightJSON

	json.Unmarshal(b, &flights)

	for _, f := range flights {
		flight := models.Flight{
			Id:              f.Id,
			SourceAirportId: f.SourceAirportId,
			DestAirportId:   f.DestAirportId,
			Passengers:      f.Passengers,
			Seats:           f.Seats,
			Queue:           make(chan *models.Session, 10),
		}
		dao.data[flight.SourceAirportId] = make(map[uuid.UUID]*models.Flight)
		dao.data[flight.SourceAirportId][flight.DestAirportId] = &flight
	}
}

func (dao *MemoryFlightDAO) FindAll() []*models.Flight {
	v := make([]*models.Flight, 0, len(dao.data))

	for _, array := range dao.data {
		for _, flight := range array {
			v = append(v, flight)
		}
	}

	return v
}

func (dao *MemoryFlightDAO) Insert(t *models.Flight) {
	id := uuid.New()

	t.Id = id

	if dao.data[t.SourceAirportId] == nil {
		dao.data[t.SourceAirportId] = make(map[uuid.UUID]*models.Flight)
	}

	dao.data[t.SourceAirportId][t.DestAirportId] = t
}

func (dao *MemoryFlightDAO) Update(t *models.Flight) error {

	_, exists := dao.data[t.SourceAirportId][t.DestAirportId]

	if !exists {
		return errors.New("not found")
	}

	dao.data[t.SourceAirportId][t.DestAirportId] = t

	return nil
}

func (dao *MemoryFlightDAO) Delete(t models.Flight) error {
	delete(dao.data[t.SourceAirportId], t.DestAirportId)

	_, exists := dao.data[t.SourceAirportId][t.DestAirportId]

	if exists {
		return errors.New("delete was unsuccessful")
	}

	return nil
}

func (dao *MemoryFlightDAO) FindById(id uuid.UUID) (*models.Flight, error) {
	for _, array := range dao.data {
		for _, flight := range array {
			if flight.Id == id {
				return flight, nil
			}
		}
	}

	return nil, errors.New("flight not found")
}

func (dao *MemoryFlightDAO) FindBySource(id uuid.UUID) ([]*models.Flight, error) {
	t, exists := dao.data[id]

	if !exists {
		return nil, errors.New("airport not found")
	}

	flights := make([]*models.Flight, 0, len(t))

	for _, flight := range t {
		flights = append(flights, flight)
	}

	return flights, nil
}

func (dao *MemoryFlightDAO) FindBySourceAndDest(source uuid.UUID, dest uuid.UUID) (*models.Flight, error) {
	t, exists := dao.data[source][dest]

	if !exists {
		return nil, errors.New("flight not found")
	}

	return t, nil
}

func (dao *MemoryFlightDAO) BreadthFirstSearch(source uuid.UUID, dest uuid.UUID) ([]*models.Flight, error) {
	visited := make(map[uuid.UUID]bool, len(dao.data))
	queue := []uuid.UUID{source}
	visited[source] = true
	parent := make(map[uuid.UUID]uuid.UUID, len(dao.data))
	parent[source] = source

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current == dest {
			break
		}

		for neighbor, flight := range dao.data[current] {
			if !visited[neighbor] && flight.Seats > 0 {
				visited[neighbor] = true
				queue = append(queue, neighbor)
				parent[neighbor] = current
			}
		}
	}

	path := []*models.Flight{}
	current := dest

	if !visited[dest] {
		return nil, errors.New("no route available")
	}

	for current != source {
		prev := parent[current]
		flight := dao.data[prev][current]
		path = append([]*models.Flight{flight}, path...)
		current = prev
	}

	return path, nil
}

func (dao *MemoryFlightDAO) DeleteAll() {
	dao.data = make(map[uuid.UUID]map[uuid.UUID]*models.Flight)
}
