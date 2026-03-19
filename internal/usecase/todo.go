package usecase

import (
	"fmt"
	"time"
	"todo-clean/internal/entity"

	"github.com/google/uuid"
)

type TodoUseCase struct {
	repo TodoRepository // depends on the interface, not a concrete type
}

// Constructor — this is your DI. No magic, just a function.
func NewTodoUseCase(repo TodoRepository) *TodoUseCase {
	return &TodoUseCase{repo: repo}
}

func (uc *TodoUseCase) CreateTodo(title string) (*entity.Todo, error) {
	if title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}
	todo := &entity.Todo{
		ID:        uuid.NewString(),
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}
	if err := uc.repo.Save(todo); err != nil {
		return nil, fmt.Errorf("CreateTodo: %w", err)
	}
	return todo, nil
}

func (uc *TodoUseCase) ListTodos() ([]entity.Todo, error) {
	todos, err := uc.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("ListTodos: %w", err)
	}
	return todos, nil
}

func (uc *TodoUseCase) CompleteTodo(id string) error {
	todo, err := uc.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("CompleteTodo: %w", err)
	}
	todo.Done = true
	return uc.repo.Save(todo)
}

func (uc *TodoUseCase) DeleteTodo(id string) error {
	return uc.repo.Delete(id)
}
