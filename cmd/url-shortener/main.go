package main

import (
	"flag"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
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
	stType := flag.String("storage-type", "", "Specify the type of storage to use")
	flag.Parse()

	if *stType == "" {
		log.Fatalln("Argument --storage-type is required")
	}

	conf := config.MustLoad(*stType)

	logger := setupLogger(conf.Env)

	logger.Info("Starting service", slog.String("env", conf.Env))

	storage, err := storage_factory.GetStorage(*stType, conf)
	if err != nil {
		logger.Error("Error creating storage: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Storage started", slog.String("storage", *stType))

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(logger, storage))
	router.Get("/url/{alias}", redirect.New(logger, storage))

	logger.Info("Starting server", slog.String("address", conf.HTTPServer.Addr))

	srv := &http.Server{
		Addr:         conf.HTTPServer.Addr,
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
	var logger *slog.Logger

	switch envType {
	case env.EnvLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case env.EnvDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case env.EnvProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
