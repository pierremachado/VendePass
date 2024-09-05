package models

import (
	"github.com/google/uuid"
)

type Client struct {
	Id             uuid.UUID
	Name           string
	Username       string
	Password       string
	Client_flights []Ticket
}
