package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"svc-task_master/src/common/config"
	"svc-task_master/src/common/logger"
	"svc-task_master/src/ports_adapters/primary/http_server"
	"svc-task_master/src/ports_adapters/secondary/inmemory/db"
	"svc-task_master/src/ports_adapters/secondary/service/application"
	"syscall"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "svc-task_master/docs"
)

// @title task_master API
// @version 1.0
// @description API для управления задачами
// @host localhost:8080
// @BasePath /
func main() {

	cfg := config.LoadConfig()
	asyncLogeer := logger.NewCustomAsyncLogger(cfg.Logger.BathSize, cfg.Logger.LogLvl)

	asyncLogeer.Info("Application is starting...")
	asyncLogeer.Info("Loaded configuration", slog.Any("config", cfg))

	asyncLogeer.Info("Initializing repository...")
	repo := db.NewRepository(asyncLogeer, cfg.MemoryDB.NumShards, cfg.MemoryDB.TTL)

	asyncLogeer.Info("Initializing application service...")
	app := application.InitApp(repo.InMemoryDB, asyncLogeer)

	asyncLogeer.Info("Initializing HTTP server...")
	s := http_server.NewServer(&app)
	r := http_server.NewRouter()

	r.PUT("/task/:id", s.UpdateStatusTask)
	r.POST("/task", s.CreateTask)
	r.GET("/task/:id", s.GetTaskForId)
	r.GET("/task", s.GetTasksSortStatus)
	r.Handle("GET", "/swagger/*", httpSwagger.WrapHandler)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: r,
	}

	go func() {
		asyncLogeer.Info(fmt.Sprintf("Starting server on port %s", cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			asyncLogeer.Error("Failed to start server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	<-done
	asyncLogeer.Info("Server received shutdown signal")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	asyncLogeer.Info("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		asyncLogeer.Error("Server shutdown failed", slog.String("error", err.Error()))
	} else {
		asyncLogeer.Info("Server shutdown completed successfully")
	}

	asyncLogeer.Info("Shutting down logger...")
	asyncLogeer.Info("Application exited properly")
	asyncLogeer.Shutdown()

}
