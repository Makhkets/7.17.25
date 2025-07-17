package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/makhkets/7.17.25/internal/api"
	"github.com/makhkets/7.17.25/internal/config"
	"github.com/makhkets/7.17.25/internal/repository"
	"github.com/makhkets/7.17.25/internal/service"
	logging "github.com/makhkets/7.17.25/pkg/logger"
)

func main() {
	cfg := config.MustLoad("local.json")
	logging.SetupLogger()
	slog.Info("starting application", slog.Any("config", cfg))

	repo := repository.NewRepo()
	srvc := service.NewService(repo)
	server := api.InitServer(srvc, cfg)

	defer shutdown(server)
}

func shutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", slog.String("error", err.Error()))
	}

	slog.Info("server exited")
	os.Exit(0)
}
