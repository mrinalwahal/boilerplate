package todo

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mrinalwahal/boilerplate/services/todo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Create Handler.
func Create(c echo.Context) error {

	//	Unmarshal the incoming payload.
	var payload CreateOptions
	if err := c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, "Invalid payload.")
	}

	//	Prepare a database connection.
	dsn := "host=127.0.0.1 user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to connect to the database.")
	}

	//	Initialize a default context.
	ctx := context.Background()

	//	Get the todo service.
	service := todo.GetService(db)

	//	Call the service function to execute the business logic.
	todo, err := service.Create(ctx, payload.Title)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create the todo.")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"todo": todo,
	})
}

// Get Handler.
func Get(c echo.Context) error {

	//	Get the object ID from the URL.
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID.")
	}

	//	Prepare a database connection.
	dsn := "host=127.0.0.1 user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to connect to the database.")
	}

	//	Initialize a default context.
	ctx := context.Background()

	//	Get the todo service.
	service := todo.GetService(db)

	//	Call the service function to execute the business logic.
	todo, err := service.Get(ctx, uuid)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create the todo.")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"todo": todo,
	})
}

// List Handler.
func List(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload ListOptions
	if err := c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, "Invalid payload.")
	}

	//	Prepare a database connection.
	dsn := "host=127.0.0.1 user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to connect to the database.")
	}

	//	Initialize a default context.
	ctx := context.Background()

	//	Get the todo service.
	service := todo.GetService(db)

	//	Call the service function to execute the business logic.
	todo, err := service.List(ctx, &todo.ListOptions{
		Skip:           payload.Skip,
		Limit:          payload.Limit,
		Title:          payload.Title,
		OrderBy:        payload.OrderBy,
		OrderDirection: payload.OrderDirection,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create the todo.")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"todos": todo,
	})
}

// Update Handler.
func Update(c echo.Context) error {

	//	Get the object ID from the URL.
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID.")
	}

	//	Unmarshal the incoming payload.
	var payload UpdateOptions
	if err := c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, "Invalid payload.")
	}

	//	Prepare a database connection.
	dsn := "host=127.0.0.1 user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to connect to the database.")
	}

	//	Initialize a default context.
	ctx := context.Background()

	//	Get the todo service.
	service := todo.GetService(db)

	//	Call the service function to execute the business logic.
	todo, err := service.Update(ctx, uuid, &todo.UpdateOptions{
		Title: payload.Title,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create the todo.")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"todo": todo,
	})
}

// Delete Handler.
func Delete(c echo.Context) error {

	//	Get the object ID from the URL.
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID.")
	}

	//	Prepare a database connection.
	dsn := "host=127.0.0.1 user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to connect to the database.")
	}

	//	Initialize a default context.
	ctx := context.Background()

	//	Get the todo service.
	service := todo.GetService(db)

	//	Call the service function to execute the business logic.
	err = service.Delete(ctx, uuid)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create the todo.")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Todo deleted successfully.",
	})
}
