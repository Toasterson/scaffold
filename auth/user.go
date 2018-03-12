package auth

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func CryptPass(pass string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(bytes)
}

func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

type User struct {
	UUID          uuid.UUID `json:"uuid" gorm:"'uuid' pk"`
	Username      string    `json:"username" gorm:"unique"`
	Vorname       string    `json:"vorname"`
	Nachname      string    `json:"nachname"`
	Password      string    `json:"-"`
	ClearPassword string    `json:"password" gorm:"-"` //Only used to receive the cleartext password from the client
	Email         string    `json:"email"`
	CreatedAt     time.Time `gorm:"created" json:"created_at"`
	UpdatedAt     time.Time `gorm:"updated" json:"updated_at"`
}

func NewUser(username, email, password string) (u *User) {
	u = new(User)
	u.Username = username
	u.Email = email
	u.Password = CryptPass(password)
	u.UUID, _ = uuid.NewV4()
	return
}

func (u *User) UpdatePassword(oldpass, newpass string) error {
	if err := ComparePassword(u.Password, oldpass); err != nil {
		return fmt.Errorf("password does not match")
	}
	u.Password = CryptPass(newpass)
	return nil
}
