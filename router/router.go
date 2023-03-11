package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/usagrada/login-system/db"
)

func Router(e *echo.Group) {
	e.GET("", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})
	e.POST("/initialize", func(c echo.Context) error {
		db.Initialize()
		return c.String(200, "Initialized database!")
	})
	e.GET("/login", func(c echo.Context) error {
		return c.String(200, "Login")
	})
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
			Name:  "token",
			Value: tokenString,
		})
		return c.String(200, "JWT token: "+tokenString)
	})
	e.POST("/signup", func(c echo.Context) error {
		ref := c.Request().Referer()
		fmt.Println("Request: Signup", ref)
		return c.String(200, "Signup")
	})
	e.GET("/users", func(c echo.Context) error {
		type User struct {
			Id       int    `json:"id"`
			Username string `json:"username"`
			Password string `json:"password"`
		}
		var users []User
		rows, err := db.DB.Query("SELECT * FROM users")
		if err != nil {
			panic(err)
		}

		for rows.Next() {
			u := &User{}
			err := rows.Scan(&u.Id, &u.Username, &u.Password)
			if err != nil {
				panic(err)
			}
			fmt.Println(u)
			users = append(users, *u)
		}

		return c.JSON(200, users)
	})
}
