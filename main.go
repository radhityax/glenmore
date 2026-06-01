package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"fmt"

	"glenmore/internal/db"
	"glenmore/internal/server"
)

func main() {
	cfg, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	database, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer database.Close()

	if err := db.Migrate(database); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	slog.Info("database migrated")

	router := server.NewRouter(database, server.Config{
		Host: cfg.Host,
		Port: cfg.Port,
	})

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT,
	syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("server strated", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down..")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}
