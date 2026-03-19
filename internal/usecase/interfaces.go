package usecase

import "todo-clean/internal/entity"

// TodoRepository is defined HERE, inside the usecase package.
// Not in the repository package. This is the inversion.
type TodoRepository interface {
	GetAll() ([]entity.Todo, error)
	GetByID(id string) (*entity.Todo, error)
	Save(todo *entity.Todo) error
	Delete(id string) error
}
