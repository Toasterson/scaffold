package auth

import (
	"github.com/satori/go.uuid"
)

type Session struct {
	UUID  uuid.UUID `gorm:"'uuid' pk"`
	Token string    `gorm:"unique"`
}

func NewSession(ID uuid.UUID, encodedToken string) *Session {
	return &Session{
		UUID:  ID,
		Token: encodedToken,
	}
}
