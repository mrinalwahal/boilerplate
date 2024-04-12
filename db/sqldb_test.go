package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/mrinalwahal/service/model"
	"github.com/mrinalwahal/service/pkg/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Temporary testsqldbconfig that contains all the configuration required by our tests.
type testsqldbconfig struct {

	// Test database connection.
	conn *gorm.DB
}

// Setup the test environment.
func configure(t *testing.T) *testsqldbconfig {

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

func Test_Database_Create(t *testing.T) {

	// Setup the test environment.
	environment := configure(t)

	// Initialize the database.
	db := &sqldb{
		conn: environment.conn,
	}

	t.Run("create record with nil options", func(t *testing.T) {

		_, err := db.Create(context.Background(), nil)
		if err == nil || err != ErrInvalidOptions {
			t.Errorf("service.Create() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("create record with options w/o JWT claims", func(t *testing.T) {

		options := &CreateOptions{
			Title: "Test Record",
		}

		_, err := db.Create(context.Background(), options)
		if err == nil {
			t.Errorf("service.Create() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("create record with options w/ JWT claims", func(t *testing.T) {

		options := &CreateOptions{
			Title: "Test Record",
		}

		// Add JWT claims to the context.
		ctx := context.WithValue(context.Background(), middleware.XJWTClaims, jwt.MapClaims{
			XUserID: uuid.New(),
		})
		record, err := db.Create(ctx, options)
		if err != nil {
			t.Fatalf("failed to create a record: %v", err)
		}

		// Check if the response contains a valid UUID and correct title.
		if record.ID == uuid.Nil {
			t.Fatalf("expected record ID to be generated automatically, got empty UUID")
		}

		if record.Title != "Test Record" {
			t.Fatalf("expected record title to be 'Test Record', got '%s'", options.Title)
		}
	})
}

func Test_Database_List(t *testing.T) {

	// Setup the test environment.
	environment := configure(t)

	// Initialize the database.
	db := &sqldb{
		conn: environment.conn,
	}

	// Add JWT claims to the context.
	ctx := context.WithValue(context.Background(), middleware.XJWTClaims, jwt.MapClaims{
		XUserID: uuid.New(),
	})

	// Seed the database with some records.
	for i := 0; i < 5; i++ {
		_, err := db.Create(ctx, &CreateOptions{
			Title: fmt.Sprintf("Record %d", i),
		})
		if err != nil {
			t.Fatalf("failed to seed the database: %v", err)
		}
	}

	t.Run("list all records", func(t *testing.T) {

		records, err := db.List(ctx, &ListOptions{})
		if err != nil {
			t.Fatalf("failed to list records: %v", err)
		}

		if len(records) < 1 {
			t.Fatalf("expected at least 1 record, got %d", len(records))
		}
	})

	t.Run("list w/ title filter", func(t *testing.T) {

		records, err := db.List(ctx, &ListOptions{
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

		records, err := db.List(ctx, &ListOptions{
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

		records, err := db.List(ctx, &ListOptions{
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

		records, err := db.List(ctx, &ListOptions{
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

		records, err := db.List(ctx, &ListOptions{
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

func Test_Database_Get(t *testing.T) {

	// Setup the test environment.
	environment := configure(t)

	// Initialize the database.
	db := &sqldb{
		conn: environment.conn,
	}

	// Seed the database with sample records.
	options := CreateOptions{
		Title: "Test Record",
	}

	// Add JWT claims to the context.
	ctx := context.WithValue(context.Background(), middleware.XJWTClaims, jwt.MapClaims{
		XUserID: uuid.New(),
	})

	seed, err := db.Create(ctx, &options)
	if err != nil {
		t.Fatalf("failed to seed the database: %v", err)
	}

	t.Run("get seed record", func(t *testing.T) {

		record, err := db.Get(ctx, seed.ID)
		if err != nil {
			t.Fatalf("failed to get record: %v", err)
		}

		if record.ID != seed.ID {
			t.Fatalf("expected retrieved record to equal seed, got = %v", record)
		}
	})
}

func Test_Database_Update(t *testing.T) {

	// Setup the test environment.
	environment := configure(t)

	// Initialize the database.
	db := &sqldb{
		conn: environment.conn,
	}

	// Seed the database with sample records.
	options := CreateOptions{
		Title: "Test Record",
	}

	// Add JWT claims to the context.
	ctx := context.WithValue(context.Background(), middleware.XJWTClaims, jwt.MapClaims{
		XUserID: uuid.New(),
	})

	seed, err := db.Create(ctx, &options)
	if err != nil {
		t.Fatalf("failed to seed the database: %v", err)
	}

	t.Run("update seed record", func(t *testing.T) {

		options := UpdateOptions{
			Title: "Updated Title",
		}
		updated, err := db.Update(ctx, seed.ID, &options)
		if err != nil {
			t.Fatalf("failed to update record: %v", err)
		}

		if updated.Title != options.Title {
			t.Fatalf("expected updated record title to be '%s', got '%s'", options.Title, updated.Title)
		}
	})
}

func Test_Database_Delete(t *testing.T) {

	// Setup the test environment.
	environment := configure(t)

	// Initialize the database.
	db := &sqldb{
		conn: environment.conn,
	}

	// Seed the database with sample records.
	options := CreateOptions{
		Title: "Test Record",
	}

	// Add JWT claims to the context.
	ctx := context.WithValue(context.Background(), middleware.XJWTClaims, jwt.MapClaims{
		XUserID: uuid.New(),
	})

	seed, err := db.Create(ctx, &options)
	if err != nil {
		t.Fatalf("failed to seed the database: %v", err)
	}

	t.Run("delete seed record", func(t *testing.T) {

		err := db.Delete(context.Background(), seed.ID)
		if err != nil {
			t.Fatalf("failed to delete record: %v", err)
		}

		_, err = db.Get(context.Background(), seed.ID)
		if err == nil {
			t.Fatalf("expected record to be deleted, got nil")
		}
	})
}