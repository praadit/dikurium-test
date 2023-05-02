package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/praadit/dikurium-test/pkg"
	"github.com/praadit/dikurium-test/pkg/config"
	"github.com/praadit/dikurium-test/pkg/middlewares"
	"github.com/rs/cors"
	migrate "github.com/rubenv/sql-migrate"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Migrate(db *sql.DB) {
	migrations := &migrate.FileMigrationSource{
		Dir: "../migrations",
	}
	migrate.SetTable("migrations")

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		fmt.Printf("Error applying db migration!\n%s", err)
	}

	if n > 0 {
		fmt.Printf("Applied %d migrations!\n", n)
	}
}

func main() {
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := fx.New(
		fx.Provide(
			NewLogger,
			NewChi,
			NewDb,
		),
		fx.Invoke(
			RegisterRoutes,
		),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	)

	app.Run()
}

func NewChi(lc fx.Lifecycle, logger *zap.Logger) *chi.Mux {
	logger.Info("Executing NewMux.")

	r := chi.NewRouter()

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "HEAD", "OPTIONS"},
		Debug:            true,
	}).Handler)
	r.Use(middlewares.AuthenticationMiddleware)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", "0.0.0.0", config.Config.Port),
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("Starting HTTP server.")
			go func() {
				if err := server.ListenAndServe(); err != nil {
					logger.Sugar().Error(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping HTTP server.")
			return server.Shutdown(ctx)
		},
	})

	return r
}

func NewDb(lc fx.Lifecycle, logger *zap.Logger) *gorm.DB {
	db := config.InitDb()
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}

	Migrate(sqlDb)

	return db
}

func NewLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	logger.Info("Executing NewLogger.")
	return logger
}

func RegisterRoutes(router *chi.Mux, db *gorm.DB, logger *zap.Logger) {
	server := pkg.NewServer(db, logger)

	router.Handle("/", server.PlaygroundHandler())
	router.Handle("/query", server.GraphqlHandler())
}
