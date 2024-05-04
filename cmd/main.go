package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"ozinshe/cmd/server"
	"ozinshe/internal/config"
	"ozinshe/internal/handler"
	_ "ozinshe/internal/models"
	"ozinshe/internal/service"
	storage "ozinshe/internal/storage/postgresql"
	"syscall"
)

// @title Ozinshe API
// @version 1.0
// description API Server for Ozinshe Application

// @host localhost:8080
// @BasePath /

// @securityDefinitions.cookie CookieAuth
// @in cookie
// @name JWT
// @description Type "JWT" token received from the server

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	db, err := storage.NewPostgres(cfg)
	if err != nil {
		panic(err)
	}

	storage := storage.New(db)
	service := service.New(storage)
	handler := handler.New(service, log)

	srv := new(server.Server)
	err = srv.Run(cfg.Port, handler.InitRoutes())

	if err != nil {
		log.Error("error while running http server", err)
		return
	}

	log.Info("server started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("Stopping application", slog.String("signal", sign.String()))

	srv.MustShutdown(context.Background())

	log.Info("Application stopped")
}
