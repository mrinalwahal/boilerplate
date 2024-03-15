package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mrinalwahal/boilerplate/api"
)

func main() {

	e := echo.New()

	//
	// Middlewares
	//

	//	Rate-limit requests to 10 per second.
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	//	Register the API routes.
	api.AddRoutes(e)

	e.Logger.Fatal(e.Start(":1323"))
}
