package models

import (
	"github.com/google/uuid"
)

type Client struct {
	Id             uuid.UUID
	Name           string
	Client_flights []Ticket
}
