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

type MemoryAirportDAO struct {
	data map[uuid.UUID]models.Airport
}

func (dao *MemoryAirportDAO) New() {

	var airports []models.Airport

	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	jsonPath := filepath.Join(baseDir, "internal", "stubs", "airports.json")

	b, _ := os.ReadFile(jsonPath)

	json.Unmarshal(b, &airports)

	for _, airport := range airports {
		dao.data[airport.Id] = airport
	}

}

func (dao *MemoryAirportDAO) FindAll() []models.Airport {
	v := make([]models.Airport, 0, len(dao.data))

	for _, value := range dao.data {
		v = append(v, value)
	}

	return v
}

func (dao *MemoryAirportDAO) Insert(t *models.Airport) {
	id := uuid.New()

	t.Id = id

	dao.data[id] = *t
}

func (dao *MemoryAirportDAO) Update(t *models.Airport) error {

	lastAirport, exists := dao.data[t.Id]

	if !exists {
		return errors.New("not found")
	}

	dao.data[t.Id] = lastAirport

	return nil
}

func (dao *MemoryAirportDAO) Delete(t models.Airport) {
	delete(dao.data, t.Id)
}

func (dao *MemoryAirportDAO) FindById(id uuid.UUID) (*models.Airport, error) {
	airport, exists := dao.data[id]

	if !exists {
		return nil, errors.New("not found")
	}

	return &airport, nil
}
