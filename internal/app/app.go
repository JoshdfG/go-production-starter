package app

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"todo-clean/config"
	todohttp "todo-clean/internal/controller/http"
	"todo-clean/internal/repo"
	"todo-clean/internal/usecase"
)

func Run(cfg *config.Config, l zerolog.Logger) error {
	// postgres
	db, err := newPostgres(cfg, l)
	if err != nil {
		return err
	}
	defer db.Close()

	// redis
	rdb, err := newRedis(cfg, l)
	if err != nil {
		return err
	}
	defer rdb.Close()

	// todo layer
	//  todo layer wiring — the chain is now:
	// usecase → LoggingRepo → CachingRepo → PostgresRepo → Postgres
	todoRepo := repo.NewPostgresRepo(db)
	cachedRepo := repo.NewCachingRepo(todoRepo, rdb, l)
	loggedRepo := repo.NewLoggingRepo(cachedRepo, l)
	todoUC := usecase.NewTodoUseCase(loggedRepo)

	// auth layer
	userRepo := repo.NewUserPostgresRepo(db)
	authUC := usecase.NewAuthUseCase(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiryHours)

	// http
	todoHandler := todohttp.NewTodoHandler(todoUC)
	authHandler := todohttp.NewAuthHandler(authUC)
	router := todohttp.NewRouter(l, todoHandler, authHandler, authUC)

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		l.Info().Str("port", cfg.HTTP.Port).Msg("listening")
		if err := router.Run(":" + cfg.HTTP.Port); err != nil {
			l.Error().Err(err).Msg("server stopped")
			quit <- syscall.SIGTERM
		}
	}()

	sig := <-quit
	l.Info().Str("signal", sig.String()).Msg("shutting down")

	// give in-flight requests 5 seconds to finish
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	l.Info().Msg("shutdown complete")
	return nil
}

func newPostgres(cfg *config.Config, l zerolog.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Postgres.DSN())
	if err != nil {
		return nil, fmt.Errorf("postgres open: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("postgres ping: %w", err)
	}

	l.Info().Msg("connected to postgres")
	return db, nil
}

func newRedis(cfg *config.Config, l zerolog.Logger) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr(),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	// warm up pool so first cache write succeeds immediately
	rdb.Do(context.Background(), "PING")

	l.Info().Msg("connected to redis")
	return rdb, nil
}
