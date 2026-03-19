package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"todo-clean/internal/entity"
	"todo-clean/internal/usecase"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

const (
	todoKeyPrefix = "todo:"
	todoListKey   = "todos:all"
	cacheTTL      = 5 * time.Minute
)

type CachingRepo struct {
	inner  usecase.TodoRepository
	client *redis.Client
	logger zerolog.Logger
}

func NewCachingRepo(inner usecase.TodoRepository, client *redis.Client, logger zerolog.Logger) *CachingRepo {
	return &CachingRepo{
		inner:  inner,
		client: client,
		logger: logger.With().Str("component", "cache").Logger(),
	}
}

func (r *CachingRepo) GetAll() ([]entity.Todo, error) {
	ctx := context.Background()

	cached, err := r.client.Get(ctx, todoListKey).Bytes()
	if err == nil {
		var todos []entity.Todo
		if unmarshalErr := json.Unmarshal(cached, &todos); unmarshalErr != nil {
			r.logger.Error().Err(unmarshalErr).Str("key", todoListKey).Msg("cache unmarshal error")
		} else {
			r.logger.Debug().Str("key", todoListKey).Int("count", len(todos)).Msg("cache HIT")
			return todos, nil
		}
	} else {
		r.logger.Debug().Str("key", todoListKey).Msg("cache MISS")
	}

	todos, err := r.inner.GetAll()
	if err != nil {
		return nil, err
	}

	data, marshalErr := json.Marshal(todos)
	if marshalErr != nil {
		r.logger.Error().Err(marshalErr).Msg("cache marshal error")
		return todos, nil
	}

	if setErr := r.client.Set(ctx, todoListKey, data, cacheTTL).Err(); setErr != nil {
		r.logger.Error().Err(setErr).Str("key", todoListKey).Msg("cache set error")
	} else {
		r.logger.Debug().Str("key", todoListKey).Msg("cache SET")
	}

	return todos, nil
}

func (r *CachingRepo) GetByID(id string) (*entity.Todo, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", todoKeyPrefix, id)

	cached, err := r.client.Get(ctx, key).Bytes()
	if err == nil {
		var todo entity.Todo
		if unmarshalErr := json.Unmarshal(cached, &todo); unmarshalErr != nil {
			r.logger.Error().Err(unmarshalErr).Str("key", key).Msg("cache unmarshal error")
		} else {
			r.logger.Debug().Str("key", key).Msg("cache HIT")
			return &todo, nil
		}
	} else {
		r.logger.Debug().Str("key", key).Msg("cache MISS")
	}

	todo, err := r.inner.GetByID(id)
	if err != nil {
		return nil, err
	}

	data, marshalErr := json.Marshal(todo)
	if marshalErr != nil {
		r.logger.Error().Err(marshalErr).Msg("cache marshal error")
		return todo, nil
	}

	if setErr := r.client.Set(ctx, key, data, cacheTTL).Err(); setErr != nil {
		r.logger.Error().Err(setErr).Str("key", key).Msg("cache set error")
	} else {
		r.logger.Debug().Str("key", key).Msg("cache SET")
	}

	return todo, nil
}

func (r *CachingRepo) Save(todo *entity.Todo) error {
	if err := r.inner.Save(todo); err != nil {
		return err
	}
	ctx := context.Background()
	r.client.Del(ctx, todoListKey)
	r.client.Del(ctx, fmt.Sprintf("%s%s", todoKeyPrefix, todo.ID))
	r.logger.Debug().Str("id", todo.ID).Msg("cache invalidated")
	return nil
}

func (r *CachingRepo) Delete(id string) error {
	if err := r.inner.Delete(id); err != nil {
		return err
	}
	ctx := context.Background()
	r.client.Del(ctx, todoListKey)
	r.client.Del(ctx, fmt.Sprintf("%s%s", todoKeyPrefix, id))
	r.logger.Debug().Str("id", id).Msg("cache invalidated")
	return nil
}
