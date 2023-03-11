package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/usagrada/login-system/db"

	"github.com/usagrada/login-system/router"
)

func main() {
	e := echo.New()
	db.Setup()
	e.Use(middleware.Recover())
	// e.Use(middleware.CORS())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:*"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "X-CSRF-Token"},
		AllowCredentials: true,
	}))
	// e.Use(middleware.Logger())
	e.Use(middleware.Secure())
	e.Use(middleware.Gzip())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "Cookie:_csrf",
		// TokenLookup:  "header:X-CSRF-Token",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		// CookieSameSite: http.SameSiteStrictMode,
	}))
	config := echojwt.Config{
		SigningKey: []byte("SECRET_KEY"),
		ParseTokenFunc: func(c echo.Context, tokenString string) (interface{}, error) {
			keyFunc := func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte("SECRET_KEY"), nil
			}

			token, err := jwt.Parse(tokenString, keyFunc)
			if err != nil {
				return nil, err
			}
			if !token.Valid {
				return nil, errors.New("invalid token")
			}
			return token, nil
		},
	}
	e.GET("/jwt", func(c echo.Context) error {
		claims := jwt.MapClaims{
			"user_id": 12345678,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		}

		// ヘッダーとペイロードの生成
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("SECRET_KEY"))
		if err != nil {
			return err
		}
		// c.SetResponse(c.Response().Header().Set("Cookie", "token="+tokenString))
		c.SetCookie(&http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Secure:   true,
			HttpOnly: true,
		})
		return c.JSON(200, map[string]string{
			"token": tokenString,
		})
	})
	e.Static("/", "frontend/dist")
	r := e.Group("/api")
	r.Use(echojwt.WithConfig(config))
	router.Router(r)

	e.Logger.Fatal(e.Start(":8080"))
}
