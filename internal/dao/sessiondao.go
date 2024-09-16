package dao

import (
	"errors"
	"vendepass/internal/models"

	"github.com/google/uuid"
)

type MemorySessionDAO struct {
	data map[uuid.UUID]models.Session
}

func (dao *MemorySessionDAO) New() {
	dao.data = make(map[uuid.UUID]models.Session)
}

func (dao *MemorySessionDAO) FindAll() []*models.Session {
	v := make([]*models.Session, 0, len(dao.data))

	for _, value := range dao.data {
		v = append(v, &value)
	}

	return v
}

func (dao *MemorySessionDAO) Insert(t *models.Session) {
	id := uuid.New()

	t.ID = id

	dao.data[id] = *t
}

func (dao *MemorySessionDAO) Update(t *models.Session) error {

	lastSession, exists := dao.data[t.ID]

	if !exists {
		return errors.New("not found")
	}

	dao.data[t.ID] = lastSession

	return nil
}

func (dao *MemorySessionDAO) Delete(t models.Session) {
	delete(dao.data, t.ID)
}

func (dao *MemorySessionDAO) FindById(id uuid.UUID) (*models.Session, error) {
	session, exists := dao.data[id]

	if !exists {
		return nil, errors.New("not found")
	}

	return &session, nil
}

func (dao *MemorySessionDAO) DeleteAll() {
	dao.data = make(map[uuid.UUID]models.Session)
}
