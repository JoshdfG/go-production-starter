package usecase

import "todo-clean/internal/entity"

type UserRepository interface {
	Create(user *entity.User) error
	GetByEmail(email string) (*entity.User, error)
}
