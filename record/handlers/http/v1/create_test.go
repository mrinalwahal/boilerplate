package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/mrinalwahal/boilerplate/pkg/middleware"
	"github.com/mrinalwahal/boilerplate/record/model"
	"github.com/mrinalwahal/boilerplate/record/service"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Contains all the configuration required by our tests.
type testconfig struct {

	// Database connection.
	service *service.MockService

	// Test log.
	log *slog.Logger
}

// Setup the test environment.
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

	return &testsqldbconfig{
		conn: conn,
	}
}

func TestCreateHandler_ServeHTTP(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	t.Run("create w/ invalid options", func(t *testing.T) {

		// Create the handler.
		handler := NewCreateHandler(&CreateHandlerConfig{
			Service: config.service,
			Logger:  config.log,
		})

		// Initialize test request and response recorder.
		r := httptest.NewRequest(http.MethodPost, "/v1/records", nil)
		w := httptest.NewRecorder()

		// The service layer should ideally not be expecting any calls to reach it.
		config.service.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

		// Serve the request.
		handler.ServeHTTP(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status code %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("create w/ valid options but w/o jwt claims", func(t *testing.T) {

		// Create the handler.
		handler := NewCreateHandler(&CreateHandlerConfig{
			Service: config.service,
			Logger:  config.log,
		})

		body, err := json.Marshal(CreateOptions{
			Title: "Test Record",
		})
		if err != nil {
			t.Fatalf("failed to marshal the dummy body for request: %v", err)
		}

		// Initialize test request and response recorder.
		r := httptest.NewRequest(http.MethodPost, "/v1/records", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		// The service layer should ideally return an error because the JWT claims are missing.
		config.service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, service.ErrInvalidOptions).Times(0)

		// Serve the request.
		handler.ServeHTTP(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status code %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("create w/ valid options and jwt claims", func(t *testing.T) {

		// Create the handler.
		handler := NewCreateHandler(&CreateHandlerConfig{
			Service: config.service,
			Logger:  config.log,
		})

		options := CreateOptions{
			Title: "Test Record",
		}
		body, err := json.Marshal(options)
		if err != nil {
			t.Fatalf("failed to marshal the dummy body for request: %v", err)
		}

		// Initialize test request and response recorder.
		r := httptest.NewRequest(http.MethodPost, "/v1/records", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		// Set the JWT claims in the request context.
		user_id := uuid.New()
		r = r.WithContext(context.WithValue(r.Context(), middleware.XJWTClaims, middleware.JWTClaims{
			XUserID: user_id,
		}))

		// The service layer is expected to return a record.
		config.service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&model.Record{
			Base: model.Base{
				ID: uuid.New(),
			},
			Title:  options.Title,
			UserID: user_id,
		}, nil).Times(1)

		// Serve the request.
		handler.ServeHTTP(w, r)

		if w.Code != http.StatusCreated {
			t.Logf("response: %s", w.Body.String())
			t.Fatalf("expected status code %d, got %d", http.StatusCreated, w.Code)
		}
	})
}
