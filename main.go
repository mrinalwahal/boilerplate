package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mrinalwahal/boilerplate/api"
)

func main() {

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	//	Register the API routes.
	api.AddRoutes(e)

	e.Logger.Fatal(e.Start(":1323"))
}
