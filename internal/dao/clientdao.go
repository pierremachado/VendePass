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

type MemoryClientDAO struct {
	data map[uuid.UUID]models.Client
}

func (dao *MemoryClientDAO) New() {

	var clients []models.Client

	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	jsonPath := filepath.Join(baseDir, "internal", "stubs", "clients.json")

	b, _ := os.ReadFile(jsonPath)

	json.Unmarshal(b, &clients)

	for _, client := range clients {
		dao.data[client.Id] = client
	}

}

func (dao *MemoryClientDAO) FindAll() []models.Client {
	v := make([]models.Client, 0, len(dao.data))

	for _, value := range dao.data {
		v = append(v, value)
	}

	return v
}

func (dao *MemoryClientDAO) Insert(t *models.Client) {
	id := uuid.New()

	t.Id = id

	dao.data[id] = *t
}

func (dao *MemoryClientDAO) Update(t *models.Client) error {

	lastClient, exists := dao.data[t.Id]

	if !exists {
		return errors.New("not found")
	}

	dao.data[t.Id] = lastClient

	return nil
}

func (dao *MemoryClientDAO) Delete(t models.Client) {
	delete(dao.data, t.Id)
}

func (dao *MemoryClientDAO) FindById(id uuid.UUID) (*models.Client, error) {
	client, exists := dao.data[id]

	if !exists {
		return nil, errors.New("not found")
	}

	return &client, nil
}
