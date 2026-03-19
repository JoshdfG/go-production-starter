package repo

import (
	"context"
	"database/sql"
	"fmt"

	"todo-clean/internal/entity"
	"todo-clean/internal/repo/sqlcgen"

	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	queries *sqlcgen.Queries
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{
		queries: sqlcgen.New(db),
	}
}

func (r *PostgresRepo) GetAll() ([]entity.Todo, error) {
	rows, err := r.queries.GetAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("GetAll: %w", err)
	}
	todos := make([]entity.Todo, len(rows))
	for i, row := range rows {
		todos[i] = toEntity(row)
	}
	return todos, nil
}

func (r *PostgresRepo) GetByID(id string) (*entity.Todo, error) {
	row, err := r.queries.GetByID(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("GetByID %s: %w", id, err)
	}
	t := toEntity(row)
	return &t, nil
}

func (r *PostgresRepo) Save(todo *entity.Todo) error {
	err := r.queries.Save(context.Background(), sqlcgen.SaveParams{
		ID:        todo.ID,
		Title:     todo.Title,
		Done:      todo.Done,
		CreatedAt: todo.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("Save %s: %w", todo.ID, err)
	}
	return nil
}

func (r *PostgresRepo) Delete(id string) error {
	err := r.queries.Delete(context.Background(), id)
	if err != nil {
		return fmt.Errorf("Delete %s: %w", id, err)
	}
	return nil
}

// toEntity converts a sqlcgen.Todo (DB type) to entity.Todo (domain type)
func toEntity(row sqlcgen.Todo) entity.Todo {
	return entity.Todo{
		ID:        row.ID,
		Title:     row.Title,
		Done:      row.Done,
		CreatedAt: row.CreatedAt,
	}
}
