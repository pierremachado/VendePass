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
	data map[uuid.UUID]map[uuid.UUID]models.Flight
}

func (dao *MemoryFlightDAO) New() {

	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	jsonPath := filepath.Join(baseDir, "internal", "stubs", "flights.json")

	b, _ := os.ReadFile(jsonPath)

	var flights []models.Flight

	json.Unmarshal(b, &flights)

	for _, flight := range flights {
		dao.data[flight.Source.Id][flight.Dest.Id] = flight
	}

}

func (dao *MemoryFlightDAO) FindAll() []models.Flight {
	v := make([]models.Flight, 0, len(dao.data))

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

	dao.data[t.Source.Id][t.Dest.Id] = *t
}

func (dao *MemoryFlightDAO) Update(t *models.Flight) error {

	lastFlight, exists := dao.data[t.Source.Id][t.Dest.Id]

	if !exists {
		return errors.New("not found")
	}

	dao.data[t.Source.Id][t.Dest.Id] = lastFlight

	return nil
}

func (dao *MemoryFlightDAO) Delete(t models.Flight) error {
	delete(dao.data[t.Source.Id], t.Dest.Id)

	_, exists := dao.data[t.Source.Id][t.Dest.Id]

	if exists {
		return errors.New("Delete was unsuccessful")
	}

	return nil
}

func (dao *MemoryFlightDAO) FindById(id uuid.UUID) (*models.Flight, error) {
	for _, array := range dao.data {
		for _, flight := range array {
			if flight.Id == id {
				return &flight, nil
			}
		}
	}

	return nil, errors.New("Flight not found")
}

func (dao *MemoryFlightDAO) FindBySource(id uuid.UUID) ([]*models.Flight, error) {
	t, exists := dao.data[id]

	if !exists {
		return nil, errors.New("Airport not found")
	}

	flights := make([]*models.Flight, 0, len(t))

	for _, flight := range t {
		flights = append(flights, &flight)
	}

	return flights, nil
}

func (dao *MemoryFlightDAO) FindBySourceAndDest(source uuid.UUID, dest uuid.UUID) (*models.Flight, error) {
	t, exists := dao.data[source][dest]

	if !exists {
		return nil, errors.New("Flight not found")
	}

	return &t, nil
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

	for current := dest; parent[current] != source; current = parent[current] {
		flight := dao.data[parent[current]][current]
		path = append([]*models.Flight{&flight}, path...)
	}

	if len(path) != 0 {
		return path, nil
	}

	return nil, errors.New("No route available")
}
