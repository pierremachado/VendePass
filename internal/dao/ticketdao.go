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

type MemoryTicketDAO struct {
	data map[uuid.UUID]models.Ticket
}

func (dao *MemoryTicketDAO) New() {
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	jsonPath := filepath.Join(baseDir, "internal", "stubs", "tickets.json")

	b, _ := os.ReadFile(jsonPath)

	var tickets []models.Ticket

	json.Unmarshal(b, &tickets)

	for _, ticket := range tickets {
		dao.data[ticket.Id] = ticket
	}

}

func (dao *MemoryTicketDAO) FindAll() []models.Ticket {
	v := make([]models.Ticket, 0, len(dao.data))

	for _, value := range dao.data {
		v = append(v, value)
	}

	return v
}

func (dao *MemoryTicketDAO) Insert(t *models.Ticket) {
	id := uuid.New()

	t.Id = id

	dao.data[id] = *t
}

func (dao *MemoryTicketDAO) Update(t *models.Ticket) error {

	lastTicket, exists := dao.data[t.Id]

	if !exists {
		return errors.New("not found")
	}

	dao.data[t.Id] = lastTicket

	return nil
}

func (dao *MemoryTicketDAO) Delete(t models.Ticket) {
	delete(dao.data, t.Id)
}

func (dao *MemoryTicketDAO) FindById(id uuid.UUID) (*models.Ticket, error) {
	ticket, exists := dao.data[id]

	if !exists {
		return nil, errors.New("not found")
	}

	return &ticket, nil
}
