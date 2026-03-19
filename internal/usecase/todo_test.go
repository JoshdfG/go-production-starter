package usecase_test

import (
	"testing"

	"todo-clean/internal/repo"
	"todo-clean/internal/usecase"
)

func newTestUseCase() *usecase.TodoUseCase {
	// InMemoryRepo satisfies TodoRepository — no DB needed
	r := repo.NewInMemoryRepo()
	return usecase.NewTodoUseCase(r)
}

func TestCreateTodo(t *testing.T) {
	uc := newTestUseCase()

	todo, err := uc.CreateTodo("Learn testing")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if todo.ID == "" {
		t.Error("expected ID to be set")
	}
	if todo.Title != "Learn testing" {
		t.Errorf("expected title 'Learn testing', got %s", todo.Title)
	}
	if todo.Done {
		t.Error("expected Done to be false on creation")
	}
}

func TestCreateTodo_EmptyTitle(t *testing.T) {
	uc := newTestUseCase()

	_, err := uc.CreateTodo("")

	if err == nil {
		t.Fatal("expected error for empty title, got nil")
	}
}

func TestCompleteTodo(t *testing.T) {
	uc := newTestUseCase()

	todo, _ := uc.CreateTodo("Complete me")
	err := uc.CompleteTodo(todo.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	todos, _ := uc.ListTodos()
	if !todos[0].Done {
		t.Error("expected Done to be true after completing")
	}
}

func TestDeleteTodo(t *testing.T) {
	uc := newTestUseCase()

	todo, _ := uc.CreateTodo("Delete me")
	err := uc.DeleteTodo(todo.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	todos, _ := uc.ListTodos()
	if len(todos) != 0 {
		t.Errorf("expected 0 todos, got %d", len(todos))
	}
}

func TestCompleteTodo_NotFound(t *testing.T) {
	uc := newTestUseCase()

	err := uc.CompleteTodo("non-existent-id")

	if err == nil {
		t.Fatal("expected error for non-existent todo, got nil")
	}
}

func TestListTodos_Empty(t *testing.T) {
	uc := newTestUseCase()

	todos, err := uc.ListTodos()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(todos) != 0 {
		t.Errorf("expected empty list, got %d todos", len(todos))
	}
}

func TestCreateMultipleTodos(t *testing.T) {
	uc := newTestUseCase()

	uc.CreateTodo("First")
	uc.CreateTodo("Second")
	uc.CreateTodo("Third")

	todos, err := uc.ListTodos()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(todos) != 3 {
		t.Errorf("expected 3 todos, got %d", len(todos))
	}
}
