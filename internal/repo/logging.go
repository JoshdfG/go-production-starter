package repo

import (
	"todo-clean/internal/entity"
	"todo-clean/internal/usecase"

	"github.com/rs/zerolog"
)

type LoggingRepo struct {
	inner  usecase.TodoRepository
	logger zerolog.Logger
}

func NewLoggingRepo(inner usecase.TodoRepository, logger zerolog.Logger) *LoggingRepo {
	return &LoggingRepo{
		inner:  inner,
		logger: logger.With().Str("component", "repo").Logger(),
	}
}

func (r *LoggingRepo) GetAll() ([]entity.Todo, error) {
	r.logger.Debug().Msg("GetAll called")
	todos, err := r.inner.GetAll()
	r.logger.Debug().Int("count", len(todos)).Err(err).Msg("GetAll done")
	return todos, err
}

func (r *LoggingRepo) GetByID(id string) (*entity.Todo, error) {
	r.logger.Debug().Str("id", id).Msg("GetByID called")
	todo, err := r.inner.GetByID(id)
	r.logger.Debug().Str("id", id).Err(err).Msg("GetByID done")
	return todo, err
}

func (r *LoggingRepo) Save(todo *entity.Todo) error {
	r.logger.Debug().Str("id", todo.ID).Str("title", todo.Title).Msg("Save called")
	err := r.inner.Save(todo)
	r.logger.Debug().Str("id", todo.ID).Err(err).Msg("Save done")
	return err
}

func (r *LoggingRepo) Delete(id string) error {
	r.logger.Debug().Str("id", id).Msg("Delete called")
	err := r.inner.Delete(id)
	r.logger.Debug().Str("id", id).Err(err).Msg("Delete done")
	return err
}
