package models

import "github.com/google/uuid"

type Airport struct {
	Id      uuid.UUID
	Name    string
	Flights []uuid.UUID
}
