package service

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/mrinalwahal/boilerplate/pkg/middleware"
	"github.com/mrinalwahal/boilerplate/record/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Temporary testconfig that contains all the configuration required by our tests.
type testconfig struct {

	// Test service connection.
	conn *gorm.DB

	// Logger instance.
	log *slog.Logger
}

// Setup the test environment.
func configure(t *testing.T) *testconfig {

	// Open an in-memory service connection with SQLite.
	conn, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open the service connection: %v", err)
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
			t.Fatalf("failed to get the service connection: %v", err)
		}
		if err := sqlDB.Close(); err != nil {
			t.Fatalf("failed to close the service connection: %v", err)
		}
	})

	return &testconfig{
		conn: conn,
		log:  slog.Default(),
	}
}

func Test_NewService(t *testing.T) {

	t.Run("create service with nil config", func(t *testing.T) {

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected NewService to panic, but it didn't")
			}
		}()

		NewService(nil)
	})

	t.Run("create service with valid config", func(t *testing.T) {

		// Setup the test config.
		config := configure(t)

		// Initialize the service.
		service := NewService(&Config{
			DB:     config.conn,
			Logger: config.log,
		})

		if service == nil {
			t.Fatalf("expected service to be initialized, got nil")
		}
	})
}

func Test_Service_Create(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the service.
	service := &service{
		db:  config.conn,
		log: config.log,
	}

	t.Run("create record with nil options", func(t *testing.T) {

		_, err := service.Create(context.Background(), nil)
		if err == nil || err != ErrInvalidOptions {
			t.Errorf("service.Create() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("create record with invalid options", func(t *testing.T) {

		options := CreateOptions{
			Title:  "",
			UserID: uuid.Nil,
		}

		_, err := service.Create(context.Background(), &options)
		if err == nil {
			t.Errorf("service.Create() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("create record with valid options", func(t *testing.T) {

		options := CreateOptions{
			Title:  "Test Record",
			UserID: uuid.New(),
		}

		record, err := service.Create(context.Background(), &options)
		if err != nil {
			t.Fatalf("failed to create record: %v", err)
		}

		if record.Title != options.Title {
			t.Fatalf("expected record title to be '%s', got '%s'", options.Title, record.Title)
		}
	})
}

func Test_Service_List(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the service.
	service := &service{
		db:  config.conn,
		log: config.log,
	}

	ctx := context.Background()

	// Seed the service with some records.
	for i := 0; i < 5; i++ {
		_, err := service.Create(ctx, &CreateOptions{
			Title:  fmt.Sprintf("Record %d", i),
			UserID: uuid.New(),
		})
		if err != nil {
			t.Fatalf("failed to seed the service: %v", err)
		}
	}

	t.Run("list records with nil options", func(t *testing.T) {

		records, err := service.List(ctx, nil)
		if err != nil {
			t.Fatalf("failed to list records: %v", err)
		}

		if len(records) < 1 {
			t.Fatalf("expected at least 1 record, got %d", len(records))
		}
	})

	t.Run("list records with invalid options", func(t *testing.T) {

		records, err := service.List(ctx, &ListOptions{
			Skip:  -1,
			Limit: -1,
		})
		if err == nil {
			t.Errorf("service.List() error = %v, wantErr %v", err, true)
		}

		if len(records) != 0 {
			t.Errorf("expected 0 records, got %d", len(records))
		}
	})

	t.Run("list records with valid options", func(t *testing.T) {

		records, err := service.List(ctx, &ListOptions{})
		if err != nil {
			t.Fatalf("failed to list records: %v", err)
		}

		if len(records) < 1 {
			t.Fatalf("expected at least 1 record, got %d", len(records))
		}
	})

	t.Run("list records as a different user than the one who created them", func(t *testing.T) {

		// Add JWT claims to the context.
		ctx := context.WithValue(context.Background(), middleware.XJWTClaims, middleware.JWTClaims{
			XUserID: uuid.New(),
		})

		records, err := service.List(ctx, &ListOptions{})
		if err != nil {
			t.Fatalf("failed to list records: %v", err)
		}

		if len(records) != 0 {
			t.Fatalf("expected 0 records, got %d", len(records))
		}
	})

	t.Run("list w/ title filter", func(t *testing.T) {

		records, err := service.List(ctx, &ListOptions{
			Title: "Record 1",
		})
		if err != nil {
			t.Fatalf("failed to list records: %v", err)
		}

		if len(records) < 1 {
			t.Fatalf("expected at least 1 record, got %d", len(records))
		}
	})

	t.Run("list w/ skip filter", func(t *testing.T) {

		records, err := service.List(ctx, &ListOptions{
			Skip: 2,
		})
		if err != nil {
			t.Fatalf("failed to list records: %v", err)
		}

		if len(records) != 3 {
			t.Fatalf("expected 3 records, got %d", len(records))
		}
	})

	t.Run("list w/ limit filter", func(t *testing.T) {

		records, err := service.List(ctx, &ListOptions{
			Limit: 2,
		})
		if err != nil {
			t.Fatalf("failed to list records: %v", err)
		}

		if len(records) != 2 {
			t.Fatalf("expected 2 records, got %d", len(records))
		}
	})

	t.Run("list w/ orderBy filter", func(t *testing.T) {

		records, err := service.List(ctx, &ListOptions{
			OrderBy: "title",
		})
		if err != nil {
			t.Fatalf("failed to list records: %v", err)
		}

		if records[3].Title != "Record 3" {
			t.Logf("received: %v", records[3])
			t.Fatalf("expected third record to be 'Record 4', got '%s'", records[3].Title)
		}
	})

	t.Run("list w/ orderBy and orderDirection filter", func(t *testing.T) {

		records, err := service.List(ctx, &ListOptions{
			OrderBy:        "title",
			OrderDirection: "desc",
		})
		if err != nil {
			t.Fatalf("failed to list records: %v", err)
		}

		if records[0].Title != "Record 4" {
			t.Fatalf("expected first record to be 'Record 4', got '%s'", records[0].Title)
		}
	})
}

func Test_Service_Get(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the service.
	service := &service{
		db:  config.conn,
		log: config.log,
	}

	// Seed the service with sample records.
	options := CreateOptions{
		Title:  "Test Record",
		UserID: uuid.New(),
	}

	ctx := context.Background()

	seed, err := service.Create(ctx, &options)
	if err != nil {
		t.Fatalf("failed to seed the service: %v", err)
	}

	t.Run("get record with nil ID", func(t *testing.T) {

		_, err := service.Get(ctx, uuid.Nil)
		if err == nil {
			t.Errorf("service.Get() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("get record with valid ID", func(t *testing.T) {

		record, err := service.Get(ctx, seed.ID)
		if err != nil {
			t.Fatalf("failed to get record: %v", err)
		}

		if record.ID != seed.ID {
			t.Fatalf("expected retrieved record to equal seed, got = %v", record)
		}
	})

	t.Run("get record as a different user than the one who created it", func(t *testing.T) {

		// Add JWT claims to the context.
		ctx := context.WithValue(context.Background(), middleware.XJWTClaims, middleware.JWTClaims{
			XUserID: uuid.New(),
		})

		_, err := service.Get(ctx, seed.ID)
		if err == nil {
			t.Errorf("service.Get() error = %v, wantErr %v", err, true)
		}
	})
}

func Test_Service_Update(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the service.
	service := &service{
		db:  config.conn,
		log: config.log,
	}

	// Seed the service with sample records.
	options := CreateOptions{
		Title:  "Test Record",
		UserID: uuid.New(),
	}

	ctx := context.Background()

	seed, err := service.Create(ctx, &options)
	if err != nil {
		t.Fatalf("failed to seed the service: %v", err)
	}

	t.Run("update record with nil ID", func(t *testing.T) {

		_, err := service.Update(ctx, uuid.Nil, &UpdateOptions{
			Title: "Updated Record",
		})
		if err == nil {
			t.Errorf("service.Update() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("update record with nil options", func(t *testing.T) {

		_, err := service.Update(ctx, seed.ID, nil)
		if err == nil {
			t.Errorf("service.Update() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("update record with invalid options", func(t *testing.T) {

		_, err := service.Update(ctx, seed.ID, &UpdateOptions{
			Title: "",
		})
		if err == nil {
			t.Errorf("service.Update() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("update record with valid options", func(t *testing.T) {

		updatedTitle := "Updated Record"
		record, err := service.Update(ctx, seed.ID, &UpdateOptions{
			Title: updatedTitle,
		})
		if err != nil {
			t.Fatalf("failed to update record: %v", err)
		}

		if record.Title != updatedTitle {
			t.Fatalf("expected record title to be 'Updated Record', got '%s'", record.Title)
		}
	})

	t.Run("update record as a different user than the one who created it", func(t *testing.T) {

		// Add JWT claims to the context.
		ctx := context.WithValue(context.Background(), middleware.XJWTClaims, middleware.JWTClaims{
			XUserID: uuid.New(),
		})

		_, err := service.Update(ctx, seed.ID, &UpdateOptions{
			Title: "Updated Record",
		})
		if err == nil {
			t.Errorf("service.Update() error = %v, wantErr %v", err, true)
		}
	})
}

func Test_Service_Delete(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the service.
	service := &service{
		db:  config.conn,
		log: config.log,
	}

	ctx := context.Background()

	t.Run("delete record with nil ID", func(t *testing.T) {

		err := service.Delete(ctx, uuid.Nil)
		if err == nil {
			t.Errorf("service.Delete() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("delete record with valid ID", func(t *testing.T) {

		seed, err := service.Create(ctx, &CreateOptions{
			Title:  "Test Record",
			UserID: uuid.New(),
		})
		if err != nil {
			t.Fatalf("failed to seed the service: %v", err)
		}

		if err := service.Delete(ctx, seed.ID); err != nil {
			t.Fatalf("failed to delete record: %v", err)
		}
	})

	t.Run("delete record as a different user than the one who created it", func(t *testing.T) {

		seed, err := service.Create(ctx, &CreateOptions{
			Title:  "Test Record",
			UserID: uuid.New(),
		})
		if err != nil {
			t.Fatalf("failed to seed the service: %v", err)
		}

		// Add JWT claims to the context.
		ctx := context.WithValue(context.Background(), middleware.XJWTClaims, middleware.JWTClaims{
			XUserID: uuid.New(),
		})

		err = service.Delete(ctx, seed.ID)
		if err == nil {
			t.Errorf("service.Delete() error = %v, wantErr %v", err, true)
		}
	})
}
