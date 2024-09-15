package interfaces

import (
	"vendepass/internal/models"

	"github.com/google/uuid"
)

type TicketDAO interface {
	FindAll() []models.Ticket
	Insert(*models.Ticket)
	Update(*models.Ticket) error
	Delete(models.Ticket)
	FindById(uuid.UUID) (*models.Ticket, error)
	New()
}

type FlightDAO interface {
	FindAll() []models.Flight
	Insert(*models.Flight)
	Update(*models.Flight) error
	Delete(models.Flight) error
	FindById(uuid.UUID) (*models.Flight, error)
	FindBySource(uuid.UUID) ([]*models.Flight, error)
	FindBySourceAndDest(uuid.UUID, uuid.UUID) (*models.Flight, error)
	BreadthFirstSearch(source uuid.UUID, dest uuid.UUID) ([]*models.Flight, error)
	New()
}

type ClientDAO interface {
	FindAll() []models.Client
	Insert(*models.Client)
	Update(*models.Client) error
	Delete(models.Client)
	FindById(uuid.UUID) (*models.Client, error)
	New()
}

type SessionDAO interface {
	FindAll() []*models.Session
	Insert(*models.Session)
	Update(*models.Session) error
	Delete(*models.Session)
	FindById(uuid.UUID) (*models.Session, error)
	New()
}
