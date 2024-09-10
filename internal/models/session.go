package models

import (
	"time"
	"github.com/google/uuid"
)

type Session struct {
	ID             uuid.UUID
	Client         Client
	LastTimeActive time.Time
}
