package repo_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	"todo-clean/internal/entity"
	"todo-clean/internal/repo"
)

// PostgresTestSuite runs all repository tests against a real Postgres instance.
// It reads connection details from env vars so it works both locally and in CI.
type PostgresTestSuite struct {
	suite.Suite
	db       *sql.DB
	todoRepo *repo.PostgresRepo
}

func TestPostgresSuite(t *testing.T) {
	// skip if no database URL provided — unit tests still run fine
	if os.Getenv("POSTGRES_USER") == "" {
		t.Skip("skipping integration tests: POSTGRES_USER not set")
	}
	suite.Run(t, new(PostgresTestSuite))
}

func (s *PostgresTestSuite) SetupSuite() {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		getEnvOrDefault("POSTGRES_HOST", "localhost"),
		getEnvOrDefault("POSTGRES_PORT", "5432"),
		os.Getenv("POSTGRES_DB"),
		getEnvOrDefault("POSTGRES_SSLMODE", "disable"),
	)

	db, err := sql.Open("postgres", dsn)
	s.Require().NoError(err)

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(time.Minute)

	s.Require().NoError(db.Ping())
	s.db = db

	// run migrations so schema is ready
	s.runMigrations(dsn)

	s.todoRepo = repo.NewPostgresRepo(db)
}

func (s *PostgresTestSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *PostgresTestSuite) SetupTest() {
	// clean todos table before each test — fresh slate
	_, err := s.db.ExecContext(context.Background(), "TRUNCATE TABLE todos")
	s.Require().NoError(err)
}

func (s *PostgresTestSuite) runMigrations(dsn string) {
	m, err := migrate.New("file://../../migrations", dsn)
	s.Require().NoError(err)
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		s.Require().NoError(err)
	}
}

// ── tests ────────────────────────────────────────────────────────────────────

func (s *PostgresTestSuite) TestSave_andGetByID() {
	todo := &entity.Todo{
		ID:        uuid.NewString(),
		Title:     "integration test todo",
		Done:      false,
		CreatedAt: time.Now().UTC().Truncate(time.Millisecond),
	}

	err := s.todoRepo.Save(todo)
	s.Require().NoError(err)

	got, err := s.todoRepo.GetByID(todo.ID)
	s.Require().NoError(err)
	s.Equal(todo.ID, got.ID)
	s.Equal(todo.Title, got.Title)
	s.Equal(todo.Done, got.Done)
}

func (s *PostgresTestSuite) TestSave_Update() {
	todo := &entity.Todo{
		ID:        uuid.NewString(),
		Title:     "original title",
		Done:      false,
		CreatedAt: time.Now().UTC().Truncate(time.Millisecond),
	}
	s.Require().NoError(s.todoRepo.Save(todo))

	// update — same ID, different fields
	todo.Title = "updated title"
	todo.Done = true
	s.Require().NoError(s.todoRepo.Save(todo))

	got, err := s.todoRepo.GetByID(todo.ID)
	s.Require().NoError(err)
	s.Equal("updated title", got.Title)
	s.True(got.Done)
}

func (s *PostgresTestSuite) TestGetAll_ReturnsAllTodos() {
	todos := []*entity.Todo{
		{ID: uuid.NewString(), Title: "first", CreatedAt: time.Now().UTC()},
		{ID: uuid.NewString(), Title: "second", CreatedAt: time.Now().UTC()},
		{ID: uuid.NewString(), Title: "third", CreatedAt: time.Now().UTC()},
	}

	for _, t := range todos {
		s.Require().NoError(s.todoRepo.Save(t))
	}

	all, err := s.todoRepo.GetAll()
	s.Require().NoError(err)
	s.Len(all, 3)
}

func (s *PostgresTestSuite) TestGetAll_Empty() {
	all, err := s.todoRepo.GetAll()
	s.Require().NoError(err)
	s.Empty(all)
}

func (s *PostgresTestSuite) TestGetByID_NotFound() {
	_, err := s.todoRepo.GetByID("non-existent-id")
	s.Require().Error(err)
}

func (s *PostgresTestSuite) TestDelete() {
	todo := &entity.Todo{
		ID:        uuid.NewString(),
		Title:     "to be deleted",
		CreatedAt: time.Now().UTC(),
	}
	s.Require().NoError(s.todoRepo.Save(todo))

	s.Require().NoError(s.todoRepo.Delete(todo.ID))

	all, err := s.todoRepo.GetAll()
	s.Require().NoError(err)
	s.Empty(all)
}

func (s *PostgresTestSuite) TestDelete_NonExistent() {
	// deleting non-existent ID should not error
	err := s.todoRepo.Delete("non-existent-id")
	s.Require().NoError(err)
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
