package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"github.com/toasterson/scaffold/storage"
)

var (
	userStor   storage.Storer
	signingKey string
)

func InitModule(stor *storage.Storer) {
	userStor = *stor
	signingKey = viper.GetString("signingKey")
}

func Login(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	dbUsers := []interface{}{}
	if err := userStor.All(&dbUsers, User{}); err != nil {
		return err
	} else {
		for _, dbUser := range dbUsers {
			user := dbUser.(*User)
			if u.Username == user.Username {
				if err := ComparePassword(user.Password, u.Password); err == nil {
					return newLoginToken(c, user)
				}
			}
		}
		return echo.ErrUnauthorized
	}
}

func newLoginToken(c echo.Context, u *User) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
		"username":  u.Username,
		"firstname": u.Vorname,
		"lastname":  u.Nachname,
	})

	encoded, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token_type":    "Bearer",
		"access_token":  encoded,
		"expires":       "72h",
		"refresh_token": "",
	})
}
