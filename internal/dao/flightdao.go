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
	data map[uuid.UUID]models.Flight
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

	for _, flights := range flights {
		dao.data[flights.Id] = flights
	}

}

func (dao *MemoryFlightDAO) FindAll() []models.Flight {
	v := make([]models.Flight, 0, len(dao.data))

	for _, value := range dao.data {
		v = append(v, value)
	}

	return v
}

func (dao *MemoryFlightDAO) Insert(t *models.Flight) {
	id := uuid.New()

	t.Id = id

	dao.data[id] = *t
}

func (dao *MemoryFlightDAO) Update(t *models.Flight) error {

	lastFlight, exists := dao.data[t.Id]

	if !exists {
		return errors.New("not found")
	}

	dao.data[t.Id] = lastFlight

	return nil
}

func (dao *MemoryFlightDAO) Delete(t models.Flight) {
	delete(dao.data, t.Id)
}

func (dao *MemoryFlightDAO) FindById(id uuid.UUID) (*models.Flight, error) {
	flight, exists := dao.data[id]

	if !exists {
		return nil, errors.New("not found")
	}

	return &flight, nil
}
