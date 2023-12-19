package main

import (
	"log/slog"
	"net/http"
	"os"

	"shortener/internal/config"
	"shortener/internal/lib/logger/sl"
	"shortener/internal/storage/sqlite"

	mwLogger "shortener/internal/http-server/middleware/logger"

	"shortener/internal/lib/logger/handlers/slogpretty"

	"shortener/internal/http-server/handlers/redirect"
	"shortener/internal/http-server/handlers/url/save"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoadConfig()
	log := setupLogger(cfg.Env)

	log.Info("starting shortener", slog.String("env", cfg.Env), slog.String("version", "1"))
	log.Debug("debug messages are enabled", slog.String("env", cfg.Env))

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to create storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log)) // логер, сделанный вручную, отчасти дублирует роутер чи, но он полезен
	router.Use(middleware.Recoverer)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("shortener", map[string]string{
			cfg.HTTPServer.Username: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, storage))
		//ToDo: add DELETE handler
	})

	router.Post("/url", save.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))
	// ToDo: rounter.Delete("/url/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Addr))

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
