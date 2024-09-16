package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID             uuid.UUID
	ClientID       uuid.UUID
	LastTimeActive time.Time
	Reservations   []Reservation
}
