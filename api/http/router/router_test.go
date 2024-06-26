package router

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/mrinalwahal/boilerplate/model"
	"github.com/mrinalwahal/boilerplate/pkg/middleware"
	"github.com/mrinalwahal/boilerplate/records/db"
	v1 "github.com/mrinalwahal/boilerplate/records/handlers/http/v1"
	"github.com/mrinalwahal/boilerplate/records/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// testconfig contains all the configuration that is required by our tests.
type testconfig struct {

	// Logger instance.
	log *slog.Logger

	// Service layer.
	service service.Service
}

// configure configures a suitable and reliable environment for the tests.
func configure(t *testing.T) *testconfig {

	// Open an in-memory database connection with SQLite.
	conn, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open the database connection: %v", err)
	}

	// Migrate the schema.
	if err := conn.AutoMigrate(&model.Record{}); err != nil {
		t.Fatalf("failed to migrate the schema: %v", err)
	}

	// Cleanup the environment after the test is complete.
	t.Cleanup(func() {

		// Close the connection.
		sqlDB, err := conn.DB()
		if err != nil {
			t.Fatalf("failed to get the database connection: %v", err)
		}
		if err := sqlDB.Close(); err != nil {
			t.Fatalf("failed to close the database connection: %v", err)
		}
	})

	// Initialize the database layer.
	db := db.NewSQLDB(&db.SQLDBConfig{
		DB: conn,
	})

	// Initialize the service.
	service := service.NewService(&service.Config{
		DB:     db,
		Logger: slog.Default(),
	})

	return &testconfig{
		service: service,
		log:     slog.Default(),
	}
}

func Test_Router(t *testing.T) {

	// Configure the test environment.
	config := configure(t)

	t.Run("request to create record w/ valid body", func(t *testing.T) {

		// Prepare a body with invalid JSON.
		body, err := json.Marshal(v1.CreateOptions{
			Title: "test",
		})
		if err != nil {
			t.Fatalf("failed to marshal the dummy body for request: %v", err)
		}

		// Prepare the r and response recorder.
		r := httptest.NewRequest(http.MethodPost, "/v1", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		// Set random UserID in the request context.
		ctx := context.WithValue(r.Context(), middleware.XJWTClaims, middleware.JWTClaims{
			XUserID: uuid.New(),
		})
		r = r.WithContext(ctx)

		// Prepare the router.
		router := NewHTTPRouter(&HTTPRouterConfig{
			Service: config.service,
			Logger:  config.log,
		})

		// Serve the request.
		router.ServeHTTP(w, r)

		// Check the response status code.
		if w.Code != http.StatusCreated {
			t.Logf("got response body = %v", w.Body.String())
			t.Fatalf("expected status code %d, got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("request to get record w/ valid id", func(t *testing.T) {

		claims := middleware.JWTClaims{
			XUserID: uuid.New(),
		}

		// Create a record.
		record, err := config.service.Create(context.Background(), &service.CreateOptions{
			Title:  "test",
			UserID: claims.XUserID,
		})
		if err != nil {
			t.Fatalf("failed to create a record: %v", err)
		}

		// Prepare the r and response recorder.
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/%s", record.ID), nil)
		w := httptest.NewRecorder()

		r = r.WithContext(context.WithValue(r.Context(), middleware.XJWTClaims, claims))

		// Prepare the router.
		router := NewHTTPRouter(&HTTPRouterConfig{
			Service: config.service,
			Logger:  config.log,
		})

		// Serve the request.
		router.ServeHTTP(w, r)

		// Check the response status code.
		if w.Code != http.StatusOK {
			t.Logf("got response body = %v", w.Body.String())
			t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("request to list records", func(t *testing.T) {

		// Prepare the r and response recorder.
		r := httptest.NewRequest(http.MethodGet, "/v1", nil)
		w := httptest.NewRecorder()

		ctx := context.WithValue(r.Context(), middleware.XJWTClaims, middleware.JWTClaims{
			XUserID: uuid.New(),
		})
		r = r.WithContext(ctx)

		// Prepare the router.
		router := NewHTTPRouter(&HTTPRouterConfig{
			Service: config.service,
			Logger:  config.log,
		})

		// Serve the request.
		router.ServeHTTP(w, r)

		// Check the response status code.
		if w.Code != http.StatusOK {
			t.Logf("got response body = %v", w.Body.String())
			t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// Check that the returned data in the response is a JSON array.
		var response v1.Response
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal the response body: %v", err)
		}

		if response.Data == nil {
			t.Fatalf("expected response data to be a JSON array, got nil")
		}
	})

	t.Run("request to update record w/ valid id", func(t *testing.T) {

		claims := middleware.JWTClaims{
			XUserID: uuid.New(),
		}

		// Create a record.
		record, err := config.service.Create(context.WithValue(context.Background(), middleware.XJWTClaims, claims), &service.CreateOptions{
			Title:  "test",
			UserID: claims.XUserID,
		})
		if err != nil {
			t.Fatalf("failed to create a record: %v", err)
		}

		// Prepare the body.
		body, err := json.Marshal(v1.UpdateOptions{
			Title: "updated",
		})
		if err != nil {
			t.Fatalf("failed to marshal the dummy body for request: %v", err)
		}

		// Prepare the r and response recorder.
		r := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/%s", record.ID), bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r = r.WithContext(context.WithValue(r.Context(), middleware.XJWTClaims, claims))

		// Prepare the router.
		router := NewHTTPRouter(&HTTPRouterConfig{
			Service: config.service,
			Logger:  config.log,
		})

		// Serve the request.
		router.ServeHTTP(w, r)

		// Check the response status code.
		if w.Code != http.StatusOK {
			t.Logf("got response body = %v", w.Body.String())
			t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// Validate the title of the updated record.
		var response v1.Response
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal the response body: %v", err)
		}

		if response.Data == nil {
			t.Fatalf("expected response data to be a JSON object, got nil")
		}

		data, ok := response.Data.(map[string]interface{})
		if !ok {
			t.Fatalf("expected response data to be a JSON object, got %T", response.Data)
		}

		if data["title"] != "updated" {
			t.Fatalf("expected title to be 'updated', got %s", data["title"])
		}
	})

	t.Run("request to delete record w/ valid id", func(t *testing.T) {

		claims := middleware.JWTClaims{
			XUserID: uuid.New(),
		}

		// Create a record.
		record, err := config.service.Create(context.WithValue(context.Background(), middleware.XJWTClaims, claims), &service.CreateOptions{
			Title:  "test",
			UserID: claims.XUserID,
		})
		if err != nil {
			t.Fatalf("failed to create a record: %v", err)
		}

		// Prepare the r and response recorder.
		r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/%s", record.ID), nil)
		w := httptest.NewRecorder()

		// Set random UserID in the request context.
		r = r.WithContext(context.WithValue(r.Context(), middleware.XJWTClaims, claims))

		// Prepare the router.
		router := NewHTTPRouter(&HTTPRouterConfig{
			Service: config.service,
			Logger:  config.log,
		})

		// Serve the request.
		router.ServeHTTP(w, r)

		// Check the response status code.
		if w.Code != http.StatusOK {
			t.Logf("got response body = %v", w.Body.String())
			t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// Try to fetch the deleted record and ensure it doesn't exist.
		_, err = config.service.Get(context.Background(), record.ID)
		if err == nil {
			t.Fatal("expected to get an error, got nil")
		}
	})
}
