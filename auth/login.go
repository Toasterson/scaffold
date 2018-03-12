package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

type LoginManager struct {
	sessions   map[string]jwt.Token
	db         *gorm.DB
	signingKey string
}

func NewLoginManager(signingKey string, db *gorm.DB) *LoginManager {
	if signingKey == "" {
		signingKey = viper.GetString("signingKey")
	}
	if signingKey == "" {
		if keyUUID, err := uuid.NewV4(); err != nil {
			panic(err)
		} else {
			signingKey = keyUUID.String()
		}
	}
	return &LoginManager{
		db:         db,
		signingKey: signingKey,
		sessions:   make(map[string]jwt.Token),
	}
}

func (manager *LoginManager) Login(c echo.Context) error {
	dbUser := new(User)
	u := new(User)
	if err := c.Bind(u); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error while trying to bind to user",
			Inner:   err,
		}
	}
	manager.db.First(&dbUser, &User{Username: u.Username})
	if u.Username == dbUser.Username {
		if err := ComparePassword(dbUser.Password, u.ClearPassword); err == nil {
			tokID, err := uuid.NewV4()
			if err != nil {
				return &echo.HTTPError{
					Code:    http.StatusInternalServerError,
					Message: "Error could not generate UUID for token",
					Inner:   err,
				}
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"aud":       tokID.String(),
				"exp":       time.Now().Add(time.Hour * 72).Unix(),
				"username":  u.Username,
				"firstname": u.Vorname,
				"lastname":  u.Nachname,
			})

			encoded, err := token.SignedString([]byte(manager.signingKey))
			if err != nil {
				return &echo.HTTPError{
					Code:    http.StatusInternalServerError,
					Message: "Could not Sign token",
					Inner:   err,
				}
			}
			manager.sessions[tokID.String()] = *token
			return c.JSON(http.StatusOK, map[string]string{
				"token_type": "Bearer",
				"token":      encoded,
				"expires":    "72h",
			})
		}
	}
	return echo.ErrUnauthorized
}

func (manager *LoginManager) Logout(c echo.Context) error {
	tokenRaw := c.Get("user")
	if token, ok := tokenRaw.(*jwt.Token); ok {
		claims, convertOK := token.Claims.(jwt.MapClaims)
		aud, retrieveOK := claims["aud"].(string)
		if !convertOK && !retrieveOK {
			return &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "cannot retrieve aud claim",
			}
		}
		if _, ok := manager.sessions[aud]; ok {
			delete(manager.sessions, aud)
			return c.JSON(http.StatusOK, nil)
		}
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "could not find session for Audience: Somebody is tampering with his token",
		}
	} else {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "could not retrieve token from Context",
		}
	}
}
