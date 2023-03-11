package main

import (
	"net/http"

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
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "X-CSRF-TOKEN"},
		AllowCredentials: true,
	}))
	// e.Use(middleware.Logger())
	e.Use(middleware.Secure())
	e.Use(middleware.Gzip())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		// TokenLookup: "Cookie:_csrf",
		TokenLookup:  "header:X-CSRF-TOKEN",
		CookieSecure: true,
		// CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteStrictMode,
	}))
	// e.Use(myMiddleware)
	e.Static("/", "frontend/dist")
	r := e.Group("/api")
	router.Router(r)

	e.Logger.Fatal(e.Start(":8080"))
}
