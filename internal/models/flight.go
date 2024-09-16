package models

import (
	"github.com/google/uuid"
)

type Flight struct {
	Id              uuid.UUID
	SourceAirportId uuid.UUID
	DestAirportId   uuid.UUID
	Passengers      []Ticket
	Seats           uint
}
