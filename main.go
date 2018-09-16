package main

import (
	"crypto/rsa"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"
)

func main() {

	var authDomain, corsOrigins string
	flag.StringVar(&authDomain, "auth-domain", "", "Authentication domain")
	flag.StringVar(&corsOrigins, "cors-origins", "", "Origins to allow")
	flag.Parse()

	origins := strings.Split(corsOrigins, ",")

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
	e.HidePort = true

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: origins,
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

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
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "hello"})
	})

	e.Logger.Info("CORS origins:", origins)
	e.Logger.Fatal(e.Start(":8080"))
}
