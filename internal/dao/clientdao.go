package dao

import (
	"errors"
	"vendepass/internal/models"

	"github.com/google/uuid"
)

type MemoryClientDAO struct {
	data map[uuid.UUID]models.Client
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
