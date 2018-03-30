package auth

import (
	"github.com/satori/go.uuid"
)

type Session struct {
	ID  uuid.UUID `gorm:"primary_key"`
	Token string    `gorm:"unique"`
}

func NewSession(ID uuid.UUID, encodedToken string) *Session {
	return &Session{
		ID:  ID,
		Token: encodedToken,
	}
}
