package dao

import (
	"errors"
	"sync"
	"vendepass/internal/models"

	"github.com/google/uuid"
)

type MemorySessionDAO struct {
	data map[uuid.UUID]*models.Session
	mu   sync.RWMutex
}

func (dao *MemorySessionDAO) New() {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	dao.data = make(map[uuid.UUID]*models.Session)
}

func (dao *MemorySessionDAO) FindAll() []*models.Session {
	dao.mu.RLock()
	defer dao.mu.RUnlock()

	v := make([]*models.Session, 0, len(dao.data))

	for _, value := range dao.data {
		v = append(v, value)
	}

	return v
}

func (dao *MemorySessionDAO) Insert(t *models.Session) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	id := uuid.New()
	t.ID = id
	t.Reservations = make(map[uuid.UUID]models.Reservation)
	dao.data[id] = t
}

func (dao *MemorySessionDAO) Update(t *models.Session) error {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	_, exists := dao.data[t.ID]
	if !exists {
		return errors.New("not found")
	}

	dao.data[t.ID] = t
	return nil
}

func (dao *MemorySessionDAO) Delete(t models.Session) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	delete(dao.data, t.ID)
}

func (dao *MemorySessionDAO) FindById(id uuid.UUID) (*models.Session, error) {
	dao.mu.RLock()
	defer dao.mu.RUnlock()

	session, exists := dao.data[id]
	if !exists {
		return nil, errors.New("not found")
	}

	return session, nil
}

func (dao *MemorySessionDAO) DeleteAll() {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	dao.data = make(map[uuid.UUID]*models.Session)
}
