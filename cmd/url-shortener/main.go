package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"net/http"
	"os"
	"url-service/internal/config"
	"url-service/internal/constants/env"
	"url-service/internal/http-server/handlers/redirect"
	"url-service/internal/http-server/handlers/url/save"
	"url-service/internal/storage/storage-factory"
)

func main() {
	conf := config.MustLoad()

	logger := setupLogger(conf.Env)

	logger.Info("Starting service...", slog.String("env", conf.Env))

	storage, err := storage_factory.GetStorage(conf)
	if err != nil {
		logger.Error("error creating storageType: ", slog.String("error", err.Error()))
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(logger, storage))
	router.Get("/url/{alias}", redirect.New(logger, storage))

	logger.Info("starting server", slog.String("address", conf.Addr))

	srv := &http.Server{
		Addr:         conf.Addr,
		Handler:      router,
		ReadTimeout:  conf.HTTPServer.Timeout,
		WriteTimeout: conf.HTTPServer.Timeout,
		IdleTimeout:  conf.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("error starting server", slog.String("error", err.Error()))
	}

	logger.Error("stopped server")
}

func setupLogger(envType string) *slog.Logger {
	var log *slog.Logger

	switch envType {
	case env.EnvLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case env.EnvDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case env.EnvProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
