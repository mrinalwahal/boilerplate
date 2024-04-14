package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/mrinalwahal/boilerplate/organisation/model"
	"github.com/mrinalwahal/boilerplate/pkg/middleware"
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
	if err := conn.AutoMigrate(&model.Organisation{}); err != nil {
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

func Test_NewSQLDB(t *testing.T) {

	t.Run("create db with nil config", func(t *testing.T) {

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected NewSQLDB to panic, but it didn't")
			}
		}()

		NewSQLDB(nil)
	})

	t.Run("create db with valid config", func(t *testing.T) {

		// Setup the test environment.
		environment := configure(t)

		// Initialize the database.
		db := NewSQLDB(&SQLDBConfig{
			DB: environment.conn,
		})

		if db == nil {
			t.Fatalf("expected db to be initialized, got nil")
		}
	})
}

func Test_Database_Create(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the database.
	db := &sqldb{
		conn: config.conn,
	}

	t.Run("create organisation with nil options", func(t *testing.T) {

		_, err := db.Create(context.Background(), nil)
		if err == nil || err != ErrInvalidOptions {
			t.Errorf("service.Create() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("create organisation with invalid options", func(t *testing.T) {

		options := CreateOptions{
			Title:   "",
			OwnerID: uuid.Nil,
		}

		_, err := db.Create(context.Background(), &options)
		if err == nil {
			t.Errorf("service.Create() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("create organisation with valid options", func(t *testing.T) {

		options := CreateOptions{
			Title:   "Test Organisation",
			OwnerID: uuid.New(),
		}

		organisation, err := db.Create(context.Background(), &options)
		if err != nil {
			t.Fatalf("failed to create organisation: %v", err)
		}

		if organisation.Title != options.Title {
			t.Fatalf("expected organisation title to be '%s', got '%s'", options.Title, organisation.Title)
		}
	})
}

func Test_Database_List(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the database.
	db := &sqldb{
		conn: config.conn,
	}

	ctx := context.Background()

	// Seed the database with some organisations.
	for i := 0; i < 5; i++ {
		_, err := db.Create(ctx, &CreateOptions{
			Title:   fmt.Sprintf("Organisation %d", i),
			OwnerID: uuid.New(),
		})
		if err != nil {
			t.Fatalf("failed to seed the database: %v", err)
		}
	}

	t.Run("list organisations with nil options", func(t *testing.T) {

		organisations, err := db.List(ctx, nil)
		if err != nil {
			t.Fatalf("failed to list organisations: %v", err)
		}

		if len(organisations) < 1 {
			t.Fatalf("expected at least 1 organisation, got %d", len(organisations))
		}
	})

	t.Run("list organisations with invalid options", func(t *testing.T) {

		organisations, err := db.List(ctx, &ListOptions{
			Skip:  -1,
			Limit: -1,
		})
		if err == nil {
			t.Errorf("service.List() error = %v, wantErr %v", err, true)
		}

		if len(organisations) != 0 {
			t.Errorf("expected 0 organisations, got %d", len(organisations))
		}
	})

	t.Run("list organisations with valid options", func(t *testing.T) {

		organisations, err := db.List(ctx, &ListOptions{})
		if err != nil {
			t.Fatalf("failed to list organisations: %v", err)
		}

		if len(organisations) < 1 {
			t.Fatalf("expected at least 1 organisation, got %d", len(organisations))
		}
	})

	t.Run("list organisations as a different user than the one who created them", func(t *testing.T) {

		// Add JWT claims to the context.
		ctx := context.WithValue(context.Background(), middleware.XJWTClaims, middleware.JWTClaims{
			XUserID: uuid.New(),
		})

		organisations, err := db.List(ctx, &ListOptions{})
		if err != nil {
			t.Fatalf("failed to list organisations: %v", err)
		}

		if len(organisations) != 0 {
			t.Fatalf("expected 0 organisations, got %d", len(organisations))
		}
	})

	t.Run("list w/ title filter", func(t *testing.T) {

		organisations, err := db.List(ctx, &ListOptions{
			Title: "Organisation 1",
		})
		if err != nil {
			t.Fatalf("failed to list organisations: %v", err)
		}

		if len(organisations) < 1 {
			t.Fatalf("expected at least 1 organisation, got %d", len(organisations))
		}
	})

	t.Run("list w/ skip filter", func(t *testing.T) {

		organisations, err := db.List(ctx, &ListOptions{
			Skip: 2,
		})
		if err != nil {
			t.Fatalf("failed to list organisations: %v", err)
		}

		if len(organisations) != 3 {
			t.Fatalf("expected 3 organisations, got %d", len(organisations))
		}
	})

	t.Run("list w/ limit filter", func(t *testing.T) {

		organisations, err := db.List(ctx, &ListOptions{
			Limit: 2,
		})
		if err != nil {
			t.Fatalf("failed to list organisations: %v", err)
		}

		if len(organisations) != 2 {
			t.Fatalf("expected 2 organisations, got %d", len(organisations))
		}
	})

	t.Run("list w/ orderBy filter", func(t *testing.T) {

		organisations, err := db.List(ctx, &ListOptions{
			OrderBy: "title",
		})
		if err != nil {
			t.Fatalf("failed to list organisations: %v", err)
		}

		if organisations[3].Title != "Organisation 3" {
			t.Logf("received: %v", organisations[3])
			t.Fatalf("expected third organisation to be 'Organisation 4', got '%s'", organisations[3].Title)
		}
	})

	t.Run("list w/ orderBy and orderDirection filter", func(t *testing.T) {

		organisations, err := db.List(ctx, &ListOptions{
			OrderBy:        "title",
			OrderDirection: "desc",
		})
		if err != nil {
			t.Fatalf("failed to list organisations: %v", err)
		}

		if organisations[0].Title != "Organisation 4" {
			t.Fatalf("expected first organisation to be 'Organisation 4', got '%s'", organisations[0].Title)
		}
	})
}

func Test_Database_Get(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the database.
	db := &sqldb{
		conn: config.conn,
	}

	// Seed the database with sample organisations.
	options := CreateOptions{
		Title:   "Test Organisation",
		OwnerID: uuid.New(),
	}

	ctx := context.Background()

	seed, err := db.Create(ctx, &options)
	if err != nil {
		t.Fatalf("failed to seed the database: %v", err)
	}

	t.Run("get organisation with nil ID", func(t *testing.T) {

		_, err := db.Get(ctx, uuid.Nil)
		if err == nil {
			t.Errorf("service.Get() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("get organisation with valid ID", func(t *testing.T) {

		organisation, err := db.Get(ctx, seed.ID)
		if err != nil {
			t.Fatalf("failed to get organisation: %v", err)
		}

		if organisation.ID != seed.ID {
			t.Fatalf("expected retrieved organisation to equal seed, got = %v", organisation)
		}
	})

	t.Run("get organisation as a different user than the one who created it", func(t *testing.T) {

		// Add JWT claims to the context.
		ctx := context.WithValue(context.Background(), middleware.XJWTClaims, middleware.JWTClaims{
			XUserID: uuid.New(),
		})

		_, err := db.Get(ctx, seed.ID)
		if err == nil {
			t.Errorf("service.Get() error = %v, wantErr %v", err, true)
		}
	})
}

func Test_Database_Update(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the database.
	db := &sqldb{
		conn: config.conn,
	}

	// Seed the database with sample organisations.
	options := CreateOptions{
		Title:   "Test Organisation",
		OwnerID: uuid.New(),
	}

	ctx := context.Background()

	seed, err := db.Create(ctx, &options)
	if err != nil {
		t.Fatalf("failed to seed the database: %v", err)
	}

	t.Run("update organisation with nil ID", func(t *testing.T) {

		_, err := db.Update(ctx, uuid.Nil, &UpdateOptions{
			Title: "Updated Organisation",
		})
		if err == nil {
			t.Errorf("service.Update() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("update organisation with nil options", func(t *testing.T) {

		_, err := db.Update(ctx, seed.ID, nil)
		if err == nil {
			t.Errorf("service.Update() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("update organisation with invalid options", func(t *testing.T) {

		_, err := db.Update(ctx, seed.ID, &UpdateOptions{
			Title: "",
		})
		if err == nil {
			t.Errorf("service.Update() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("update organisation with valid options", func(t *testing.T) {

		updatedTitle := "Updated Organisation"
		organisation, err := db.Update(ctx, seed.ID, &UpdateOptions{
			Title: updatedTitle,
		})
		if err != nil {
			t.Fatalf("failed to update organisation: %v", err)
		}

		if organisation.Title != updatedTitle {
			t.Fatalf("expected organisation title to be 'Updated Organisation', got '%s'", organisation.Title)
		}
	})

	t.Run("update organisation as a different user than the one who created it", func(t *testing.T) {

		// Add JWT claims to the context.
		ctx := context.WithValue(context.Background(), middleware.XJWTClaims, middleware.JWTClaims{
			XUserID: uuid.New(),
		})

		_, err := db.Update(ctx, seed.ID, &UpdateOptions{
			Title: "Updated Organisation",
		})
		if err == nil {
			t.Errorf("service.Update() error = %v, wantErr %v", err, true)
		}
	})
}

func Test_Database_Delete(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the database.
	db := &sqldb{
		conn: config.conn,
	}

	ctx := context.Background()

	t.Run("delete organisation with nil ID", func(t *testing.T) {

		err := db.Delete(ctx, uuid.Nil)
		if err == nil {
			t.Errorf("service.Delete() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("delete organisation with valid ID", func(t *testing.T) {

		seed, err := db.Create(ctx, &CreateOptions{
			Title:   "Test Organisation",
			OwnerID: uuid.New(),
		})
		if err != nil {
			t.Fatalf("failed to seed the database: %v", err)
		}

		if err := db.Delete(ctx, seed.ID); err != nil {
			t.Fatalf("failed to delete organisation: %v", err)
		}
	})

	t.Run("delete organisation as a different user than the one who created it", func(t *testing.T) {

		seed, err := db.Create(ctx, &CreateOptions{
			Title:   "Test Organisation",
			OwnerID: uuid.New(),
		})
		if err != nil {
			t.Fatalf("failed to seed the database: %v", err)
		}

		// Add JWT claims to the context.
		ctx := context.WithValue(context.Background(), middleware.XJWTClaims, middleware.JWTClaims{
			XUserID: uuid.New(),
		})

		err = db.Delete(ctx, seed.ID)
		if err == nil {
			t.Errorf("service.Delete() error = %v, wantErr %v", err, true)
		}
	})
}
