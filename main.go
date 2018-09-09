package main

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	e := echo.New()

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(`pubkey`))
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper: func(c echo.Context) bool {
			return false
		},
		ContextKey:    "user",
		SigningKey:    publicKey,
		SigningMethod: "RS256",
		TokenLookup:   "header:" + echo.HeaderAuthorization,
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
