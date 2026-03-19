package repo

import (
	"context"
	"database/sql"
	"fmt"

	"todo-clean/internal/entity"
	"todo-clean/internal/repo/sqlcgen"
)

type UserPostgresRepo struct {
	queries *sqlcgen.Queries
}

func NewUserPostgresRepo(db *sql.DB) *UserPostgresRepo {
	return &UserPostgresRepo{queries: sqlcgen.New(db)}
}

func (r *UserPostgresRepo) Create(user *entity.User) error {
	err := r.queries.CreateUser(context.Background(), sqlcgen.CreateUserParams{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("UserRepo.Create: %w", err)
	}
	return nil
}

func (r *UserPostgresRepo) GetByEmail(email string) (*entity.User, error) {
	row, err := r.queries.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, fmt.Errorf("UserRepo.GetByEmail: %w", err)
	}
	return &entity.User{
		ID:        row.ID,
		Email:     row.Email,
		Password:  row.Password,
		CreatedAt: row.CreatedAt.UTC(),
	}, nil
}

// satisfy interface — placeholder for CreatedAt type mismatch
var _ interface{ Create(*entity.User) error } = (*UserPostgresRepo)(nil)
