package models

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID             uuid.UUID
	Username       string
	Connection     net.Conn
	LastTimeActive time.Time
}
