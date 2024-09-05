package models

import (
	"github.com/google/uuid"
)

type Flight struct {
	Id         uuid.UUID
	Source     Airport
	Dest       Airport
	Passengers []Ticket
	Seats      uint
}
