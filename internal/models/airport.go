package models

import "github.com/google/uuid"

type Airport struct {
	Id   uuid.UUID `json:"Id"`
	Name string    `json:"Name"`
	City City      `json:"City"`
}
