package api

import (
	"github.com/labstack/echo/v4"
	"github.com/mrinalwahal/boilerplate/api/todo"
)

// AddRoutes adds all routes to the echo instance.
func AddRoutes(e *echo.Echo) {

	todo.AddRoutes(e.Group("/todo"))
}
