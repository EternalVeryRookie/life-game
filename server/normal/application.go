package main

import (
	"life-game/httpapi"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "hello world") })
	e.GET("/lifegame", httpapi.Simulate)
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Start(":5000")
}
