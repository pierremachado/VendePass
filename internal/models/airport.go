package models

import "github.com/google/uuid"

type Airport struct {
	Id     uuid.UUID
	Name   string
	Routes []Airport
}
