package main

import (
	"haddibanga/internal/config"
	"haddibanga/internal/server"
)

func main() {
	cfg := config.LoadEnv()
	db := config.ConnectDatabase(cfg)
	server.Start(db, cfg)
}
