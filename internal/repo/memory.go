package repo

import (
	"fmt"
	"sync"

	"todo-clean/internal/entity"
)

type InMemoryRepo struct {
	mu    sync.RWMutex
	store map[string]*entity.Todo
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{store: make(map[string]*entity.Todo)}
}

func (r *InMemoryRepo) GetAll() ([]entity.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	todos := make([]entity.Todo, 0, len(r.store))
	for _, t := range r.store {
		todos = append(todos, *t)
	}
	return todos, nil
}

func (r *InMemoryRepo) GetByID(id string) (*entity.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("todo %s not found", id)
	}
	return t, nil
}

func (r *InMemoryRepo) Save(todo *entity.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[todo.ID] = todo
	return nil
}

func (r *InMemoryRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.store, id)
	return nil
}
