package main

import (
	"database/sql"

	"todo-clean/config"
	todohttp "todo-clean/internal/controller/http"
	"todo-clean/internal/repo"
	"todo-clean/internal/usecase"
	"todo-clean/pkg/logger"

	_ "github.com/lib/pq"
)

// @title           Todo Clean API
// @version         1.0.0
// @description     Production-grade Todo API built with Clean Architecture
// @host            localhost:8080
// @BasePath        /v1
// @schemes         http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your JWT token
func main() {
	// 1. Config
	cfg, err := config.New()
	if err != nil {
		panic("config error: " + err.Error())
	}

	// 2. Logger
	l := logger.New(cfg.Log.Level, cfg.App.Name, cfg.App.Version)
	l.Info().Str("env", cfg.App.Env).Msgf("starting %s v%s", cfg.App.Name, cfg.App.Version)

	// 3. Postgres
	db, err := sql.Open("postgres", cfg.Postgres.DSN())
	if err != nil {
		l.Fatal().Err(err).Msg("failed to open db")
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		l.Fatal().Err(err).Msg("failed to connect to db")
	}
	l.Info().Msg("connected to postgres")

	// 4. Repo layer
	r := repo.NewPostgresRepo(db)
	logged := repo.NewLoggingRepo(r, l)

	// 5. Usecase layer
	uc := usecase.NewTodoUseCase(logged)

	// 6. HTTP layer
	todoHandler := todohttp.NewTodoHandler(uc)
	router := todohttp.NewRouter(l, todoHandler)

	// 7. Start
	l.Info().Str("port", cfg.HTTP.Port).Msg("listening")
	if err := router.Run(":" + cfg.HTTP.Port); err != nil {
		l.Fatal().Err(err).Msg("server error")
	}
}

// ========
// package main
//
// import (
// 	"database/sql"
// 	"log"
// 	"net/http"
//
// 	todohttp "todo-clean/internal/controller/http"
// 	"todo-clean/internal/repo"
// 	"todo-clean/internal/usecase"
//
// 	_ "github.com/lib/pq"
// )
//
// func main() {
// 	// 1. Connect to Postgres
// 	db, err := sql.Open("postgres", "postgres://todo:todo@localhost:5432/tododb?sslmode=disable")
// 	if err != nil {
// 		log.Fatalf("failed to open db: %v", err)
// 	}
// 	defer db.Close()
//
// 	if err := db.Ping(); err != nil {
// 		log.Fatalf("failed to connect to db: %v", err)
// 	}
// 	log.Println("connected to postgres")
//
// 	// 2. Build outer layer (repo)
// 	r := repo.NewPostgresRepo(db)
// 	logged := repo.NewLoggingRepo(r)
//
// 	// 3. Inject into inner layer (usecase)
// 	uc := usecase.NewTodoUseCase(logged)
//
// 	// 4. Inject into outer layer (handler)
// 	handler := todohttp.NewTodoHandler(uc)
//
// 	// 5. Start server
// 	log.Println("Listening on :8080")
// 	log.Fatal(http.ListenAndServe(":8080", handler))
// }

// ===========
// package main
//
// import (
// 	"log"
// 	"net/http"
//
// 	todohttp "todo-clean/internal/controller/http"
// 	"todo-clean/internal/repo"
// 	"todo-clean/internal/usecase"
// )
//
// func main() {
// 	// 1. Build outer layer (repo)
// 	r := repo.NewInMemoryRepo()
//
// 	logged := repo.NewLoggingRepo(r) // wrap it
//
// 	// 2. Inject into inner layer (usecase)
// 	uc := usecase.NewTodoUseCase(logged)
//
// 	// 3. Inject into outer layer (handler)
// 	handler := todohttp.NewTodoHandler(uc)
//
// 	// 4. Start server
// 	log.Println("Listening on :8080")
// 	log.Fatal(http.ListenAndServe(":8080", handler))
// }
