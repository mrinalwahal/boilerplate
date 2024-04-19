package service

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/mrinalwahal/boilerplate/organisation/db"
	"github.com/mrinalwahal/boilerplate/organisation/model"
	"go.uber.org/mock/gomock"
)

// Contains all the configuration required by our tests.
type testconfig struct {

	// Mock database layer.
	db *db.MockDB

	// Test log.
	log *slog.Logger
}

// Setup the test environment.
func configure(t *testing.T) *testconfig {

	// Get the mock database layer.
	db := db.NewMockDB(gomock.NewController(t))
	return &testconfig{
		db:  db,
		log: slog.Default(),
	}
}

func Test_NewService(t *testing.T) {

	t.Run("nil config", func(t *testing.T) {

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("NewService() did not panic")
			}
		}()

		// Initialize the service.
		NewService(nil)
	})

	t.Run("valid config w/ db", func(t *testing.T) {

		// Get the mock database layer.
		db := db.NewMockDB(gomock.NewController(t))

		// Initialize the service.
		s := NewService(&Config{
			DB: db,
		})

		if s == nil {
			t.Errorf("NewService() = %v, want a valid service", s)
		}
	})

	t.Run("valid config w/ db and logger", func(t *testing.T) {

		// Get the mock database layer.
		db := db.NewMockDB(gomock.NewController(t))

		// Initialize the service.
		s := NewService(&Config{
			DB:     db,
			Logger: slog.Default(),
		})

		if s == nil {
			t.Errorf("NewService() = %v, want a valid service", s)
		}
	})
}

func Test_Service_Create(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the service.
	s := &service{
		db:     config.db,
		logger: config.log,
	}

	t.Run("create organisation with nil options", func(t *testing.T) {

		// Make sure the database layer is not expecting a call.
		config.db.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

		_, err := s.Create(context.Background(), nil)
		if err == nil || err != ErrInvalidOptions {
			t.Errorf("service.Create() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("create organisation with invalid options", func(t *testing.T) {

		// Make sure the database layer is not expecting a call.
		config.db.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

		_, err := s.Create(context.Background(), &CreateOptions{
			Title: "",
		})
		if err == nil {
			t.Errorf("service.Create() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("create organisation with valid options", func(t *testing.T) {

		organisation := model.Organisation{
			Title: "Test Organisation",
		}

		// Set the expectation at the database layer.
		config.db.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&model.Organisation{
			Base: model.Base{
				ID: uuid.New(),
			},
			Title: organisation.Title,
		}, nil).Times(1)

		got, err := s.Create(context.Background(), &CreateOptions{
			Title:   organisation.Title,
			OwnerID: uuid.New(),
		})
		if err != nil {
			t.Errorf("service.Create() error = %v, wantErr %v", err, false)
		}
		if got.ID == uuid.Nil {
			t.Errorf("service.Create() = %v, want a valid UUID", got.ID)
		}
		if got.Title != organisation.Title {
			t.Errorf("service.Create() = %v, want %v", got.Title, organisation.Title)
		}
	})
}

func Test_Service_List(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the service.
	s := &service{
		db:     config.db,
		logger: config.log,
	}

	t.Run("list organisations with nil options", func(t *testing.T) {

		// Make sure the database layer is not expecting a call.
		config.db.EXPECT().List(gomock.Any(), gomock.Any()).Times(0)

		_, err := s.List(context.Background(), nil)
		if err == nil || err != ErrInvalidOptions {
			t.Errorf("service.List() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("list organisations with invalid options", func(t *testing.T) {

		// Make sure the database layer is not expecting a call.
		config.db.EXPECT().List(gomock.Any(), gomock.Any()).Times(0)

		_, err := s.List(context.Background(), &ListOptions{
			Skip:  -1,
			Limit: -1,
		})
		if err == nil {
			t.Errorf("service.List() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("list organisations with valid options", func(t *testing.T) {

		organisations := []*model.Organisation{
			{
				Base: model.Base{
					ID: uuid.New(),
				},
				Title: "Test Organisation",
			},
		}

		// Set the expectation at the database layer.
		config.db.EXPECT().List(gomock.Any(), gomock.Any()).Return(organisations, nil).Times(1)

		got, err := s.List(context.Background(), &ListOptions{
			Skip:  0,
			Limit: 10,
		})
		if err != nil {
			t.Errorf("service.List() error = %v, wantErr %v", err, false)
		}
		if len(got) != len(organisations) {
			t.Errorf("service.List() = %v, want %v", len(got), len(organisations))
		}
	})
}

func Test_Service_Get(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the service.
	s := &service{
		db:     config.db,
		logger: config.log,
	}

	// Sample organisation UUID.
	id := uuid.New()

	t.Run("get organisation with invalid ID", func(t *testing.T) {

		// Make sure the database layer is not expecting a call.
		config.db.EXPECT().Get(gomock.Any(), gomock.Any()).Times(0)

		_, err := s.Get(context.Background(), uuid.Nil)
		if err == nil || err != ErrInvalidOptions {
			t.Errorf("service.Get() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("get organisation with valid ID", func(t *testing.T) {

		organisation := model.Organisation{
			Base: model.Base{
				ID: id,
			},
			Title: "Test Organisation",
		}

		// Set the expectation at the database layer.
		config.db.EXPECT().Get(gomock.Any(), id).Return(&organisation, nil).Times(1)

		got, err := s.Get(context.Background(), id)
		if err != nil {
			t.Errorf("service.Get() error = %v, wantErr %v", err, false)
		}
		if got.ID != id {
			t.Errorf("service.Get() = %v, want %v", got.ID, id)
		}
		if got.Title != organisation.Title {
			t.Errorf("service.Get() = %v, want %v", got.Title, organisation.Title)
		}
	})
}

func Test_Service_Update(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the service.
	s := &service{
		db:     config.db,
		logger: config.log,
	}

	// Sample organisation UUID.
	id := uuid.New()

	t.Run("update organisation with invalid ID", func(t *testing.T) {

		// Make sure the database layer is not expecting a call.
		config.db.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		_, err := s.Update(context.Background(), uuid.Nil, &UpdateOptions{
			Title: "Test Organisation",
		})
		if err == nil || err != ErrInvalidorganisationID {
			t.Errorf("service.Update() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("update organisation with nil options", func(t *testing.T) {

		// Make sure the database layer is not expecting a call.
		config.db.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		_, err := s.Update(context.Background(), id, nil)
		if err == nil || err != ErrInvalidOptions {
			t.Errorf("service.Update() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("update organisation with invalid options", func(t *testing.T) {

		// Make sure the database layer is not expecting a call.
		config.db.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		_, err := s.Update(context.Background(), id, &UpdateOptions{
			Title: "",
		})
		if err == nil {
			t.Errorf("service.Update() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("update organisation with valid options", func(t *testing.T) {

		organisation := model.Organisation{
			Base: model.Base{
				ID: id,
			},
			Title: "Test Organisation",
		}

		// Set the expectation at the database layer.
		config.db.EXPECT().Update(gomock.Any(), id, gomock.Any()).Return(&organisation, nil).Times(1)

		got, err := s.Update(context.Background(), id, &UpdateOptions{
			Title: "Updated Organisation",
		})
		if err != nil {
			t.Errorf("service.Update() error = %v, wantErr %v", err, false)
		}
		if got.ID != id {
			t.Errorf("service.Update() = %v, want %v", got.ID, id)
		}
		if got.Title != organisation.Title {
			t.Errorf("service.Update() = %v, want %v", got.Title, organisation.Title)
		}
	})
}

func Test_Service_Delete(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the service.
	s := &service{
		db:     config.db,
		logger: config.log,
	}

	// Sample organisation UUID.
	id := uuid.New()

	t.Run("delete organisation with invalid ID", func(t *testing.T) {

		// Make sure the database layer is not expecting a call.
		config.db.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(0)

		err := s.Delete(context.Background(), uuid.Nil)
		if err == nil || err != ErrInvalidorganisationID {
			t.Errorf("service.Delete() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("delete organisation with valid ID", func(t *testing.T) {

		// Set the expectation at the database layer.
		config.db.EXPECT().Delete(gomock.Any(), id).Return(nil).Times(1)

		err := s.Delete(context.Background(), id)
		if err != nil {
			t.Errorf("service.Delete() error = %v, wantErr %v", err, false)
		}
	})
}
