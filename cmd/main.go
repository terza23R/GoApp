package main

import (
	"fmt"
	"goapp/internal/pkg/api"
	"goapp/internal/pkg/config"
	"goapp/internal/pkg/database"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.New(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to create database service: %v", err)
	}
	defer db.Close()

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	myApi := api.NewApi(addr, db, cfg.Templates.Path)
	myApi.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	sig := <-stop
	log.Printf("received signal: %s", sig.String())

	myApi.Stop()
	log.Println("shutdown complete")
}
