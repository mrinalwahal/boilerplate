package todo

import "github.com/labstack/echo/v4"

func AddRoutes(sg *echo.Group) {

	//	API v1 routes.
	v1Group := sg.Group("/v1")

	v1Group.POST("", Create)
	v1Group.GET("", List)
	v1Group.GET("/:id", Get)
	v1Group.PATCH("/:id", Update)
	v1Group.DELETE("/:id", Delete)
}
