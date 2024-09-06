package models

import "github.com/google/uuid"

type LogoutCredentials struct {
	TokenId uuid.UUID
}
