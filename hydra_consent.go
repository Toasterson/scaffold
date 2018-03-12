package main

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/ory/hydra/sdk/go/hydra"
	"github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/spf13/viper"
	"github.com/toasterson/scaffold/auth"
	"github.com/toasterson/scaffold/template"
)

var (
	client hydra.SDK
)

const state = "demostate"

func init() {
	curdir, _ := filepath.Abs("./")
	viper.SetDefault("templatesDirectory", filepath.Join(curdir, "templates"))
	viper.SetDefault("secret", "veryopenSecret")
	viper.SetDefault("listen", ":4445")
	viper.SetDefault("clientID", "demo")
	viper.SetDefault("clientSecret", "demo")
	viper.SetDefault("enpointURL", "http://localhost:4444")
	viper.SetDefault("scopes", []string{"hydra.consent"})

	viper.AddConfigPath("/etc/scaffold.yaml")

	template.InitTemplates(viper.GetString("templatesDirectory"))
}

func handleMain(c echo.Context) error {
	oauthConfig := client.GetOAuth2Config()
	oauthConfig.RedirectURL = "http://localhost:4445/callback"
	oauthConfig.Scopes = []string{"offline", "openid"}

	authURL := client.GetOAuth2Config().AuthCodeURL(state)
	return c.Render(http.StatusOK, "home.html", map[string]interface{}{
		"authURL": authURL,
	})
}

func handleConsentPost(c echo.Context) error {
	consentRequestID := c.QueryParam("consent")
	consentRequest, response, err := client.GetOAuth2ConsentRequest(consentRequestID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", map[string]string{
			"err": fmt.Sprintf("Consent Enspoint did not respond: %s", err),
		})
	} else if response.StatusCode != http.StatusOK {
		return c.Render(http.StatusInternalServerError, "error.html", map[string]string{
			"err": fmt.Sprintf("Consent Endpoint responded with %d: %s", response.StatusCode, response.Message),
		})
	}
	user := auth.Authenticated(c)
	if user == "" {
		return c.Redirect(http.StatusNotFound, "/login?consent="+consentRequestID)
	}

	grantedScopes := []string{}

	if values, err := c.FormParams(); err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", map[string]string{
			"err": fmt.Sprintf("Failed to get the Scopes for consent: %s", err),
		})
	} else {
		for key := range values {
			if key != "consent" {
				grantedScopes = append(grantedScopes, key)
			}
		}
	}

	acceptresponse, accepterr := client.AcceptOAuth2ConsentRequest(consentRequestID, swagger.ConsentRequestAcceptance{
		Subject:          user,
		GrantScopes:      grantedScopes,
		AccessTokenExtra: map[string]interface{}{"foo": "bar"},
		IdTokenExtra:     map[string]interface{}{"foo": "baz"},
	})
	if accepterr != nil {
		return c.Render(http.StatusInternalServerError, "error.html", map[string]string{
			"err": fmt.Sprintf("Accept Consent Enspoint did not respond: %s", accepterr),
		})
	} else if acceptresponse.StatusCode != http.StatusNoContent {
		return c.Render(http.StatusInternalServerError, "error.html", map[string]string{
			"err": fmt.Sprintf("Accept Consent Endpoint responded with %d: %s", acceptresponse.StatusCode, acceptresponse.Message),
		})
	}
	return c.Redirect(http.StatusFound, consentRequest.RedirectUrl)
}

func handleConsentGet(c echo.Context) error {
	consentRequestID := c.QueryParam("consent")
	consentRequest, response, err := client.GetOAuth2ConsentRequest(consentRequestID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", map[string]string{
			"err": fmt.Sprintf("Consent Enspoint did not respond: %s", err),
		})
	} else if response.StatusCode != http.StatusOK {
		return c.Render(http.StatusInternalServerError, "error.html", map[string]string{
			"err": fmt.Sprintf("Consent Endpoint responded with %d: %s", response.StatusCode, response.Message),
		})
	}
	user := auth.Authenticated(c)
	if user == "" {
		return c.Redirect(http.StatusFound, "/login?consent="+consentRequestID)
	}

	return c.Render(http.StatusOK, "consent.html", map[string]interface{}{
		"consentRequest":   consentRequest,
		"consentRequestID": consentRequestID,
	})
}

func handleLoginGet(c echo.Context) error {
	consentRequestID := c.QueryParam("consent")
	return c.Render(http.StatusOK, "login.html", map[string]string{
		"constenRequestID": consentRequestID,
	})
}

func handleLoginPost(c echo.Context) error {
	consentRequestID := c.QueryParam("consent")
	username := c.FormValue("username")
	password := c.FormValue("password")
	if err := auth.AuthenticateUser(username, password); err != nil {
		return c.Render(http.StatusBadRequest, "error.html", map[string]string{
			"err": "Wrong Username, password",
		})
	} else {
		session, _ := sessionStore.Get(c.Request(), "")
		session.Values["user"] = "username"
		if err := sessionStore.Save(c.Request(), c.Response(), session); err != nil {
			return c.Render(http.StatusBadRequest, "error.html", map[string]string{
				"err": fmt.Sprintf("could not persist cookie: %s", err),
			})
		}
	}
	return c.Redirect(http.StatusFound, "/consent?consent="+consentRequestID)
}

func handleCallback(c echo.Context) error {

	token, err := client.GetOAuth2Config().Exchange(context.Background(), c.QueryParam("code"))
	if err != nil {
		return c.Render(http.StatusBadRequest, "error.html", map[string]string{
			"err": fmt.Sprintf("could not exchange token: %s", err),
		})
	}

	return c.Render(http.StatusOK, "callback.html", map[string]interface{}{
		"token":   token,
		"IDToken": token.Extra("id_token"),
	})
}

func main() {
	var err error

	sessionStore = sessions.NewCookieStore([]byte(viper.GetString("secret")))

	client, err = hydra.NewSDK(&hydra.Configuration{
		ClientID:     viper.GetString("clientID"),
		ClientSecret: viper.GetString("clientSecret"),
		EndpointURL:  viper.GetString("enpointURL"),
		Scopes:       viper.GetStringSlice("scopes"),
	})
	if err != nil {
		panic(fmt.Errorf("could not open hydra SDK: %s", err))
	}

	e := echo.New()
	e.Renderer = template.DefaultRenderer

	e.GET("/", handleMain)
	e.GET("/login", handleLoginGet)
	e.GET("/consent", handleConsentGet)
	e.GET("/callback", handleCallback)

	e.POST("/login", handleLoginPost)
	e.POST("/consent", handleConsentPost)
	e.Start(viper.GetString("listen"))
}
