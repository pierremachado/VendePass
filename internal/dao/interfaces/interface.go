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
}

type FlightDAO interface {
	FindAll() []models.Flight
	Insert(*models.Flight)
	Update(*models.Flight) error
	Delete(models.Flight)
	FindById(uuid.UUID) (*models.Flight, error)
}

type ClientDAO interface {
	FindAll() []models.Client
	Insert(*models.Client)
	Update(*models.Client) error
	Delete(models.Client)
	FindById(uuid.UUID) (*models.Client, error)
}
