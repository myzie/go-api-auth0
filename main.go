package main

import (
	"crypto/rsa"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/namsral/flag"
	"github.com/pavel-kiselyov/echo-logrusmiddleware"
	log "github.com/sirupsen/logrus"
)

func main() {

	var authDomain string
	flag.StringVar(&authDomain, "auth-domain", "", "Authentication domain")
	flag.Parse()

	keys, err := GetKeys("https://" + authDomain + "/.well-known/jwks.json")
	if err != nil {
		log.Fatal(err)
	}

	if len(keys.Keys) == 0 {
		log.Fatal("No public key found")
	}

	var publicKeys []*rsa.PublicKey

	for _, key := range keys.Keys {
		cert := key.GetCertificate()
		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		if err != nil {
			log.Fatal(err)
		}
		publicKeys = append(publicKeys, verifyKey)
	}

	e := echo.New()

	e.HideBanner = true
	e.HidePort = false
	e.Logger = logrusmiddleware.Logger{Logger: log.StandardLogger()}

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(logrusmiddleware.Hook())
	e.Use(middleware.Recover())

	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper: func(c echo.Context) bool {
			return false
		},
		ContextKey:    "user",
		SigningKey:    publicKeys[0],
		SigningMethod: "RS256",
		AuthScheme:    "Bearer",
		TokenLookup:   "header:" + echo.HeaderAuthorization,
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
