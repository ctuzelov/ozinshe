package main

import (
	"ozinshe/internal/config"
	storage "ozinshe/internal/storage/postgresql"
)

func main() {
	cfg := config.MustLoad()
	storage.New(cfg)

	// TODO: initialize service logic

	// TODO: initialize http server

	// TODO: initialize graceful shutdown

	// TODO: run the http server
}
