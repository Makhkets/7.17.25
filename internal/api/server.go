package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	api "github.com/makhkets/7.17.25/internal/api/gen"
	"github.com/makhkets/7.17.25/internal/api/handlers"
	"github.com/makhkets/7.17.25/internal/config"
)

func InitServer(service handlers.Service, config *config.Config) *http.Server {
	serverApi := handlers.NewServerAPI(service, config)
	ogenServer, err := api.NewServer(serverApi)
	if err != nil {
		slog.Error("cannot create server", slog.String("error", err.Error()))
		panic(err)
	}

	addr := config.App.Address + ":" + strconv.Itoa(config.App.Port)
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      ogenServer,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("starting server", slog.String("address", addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", slog.String("error", err.Error()))
			panic(err)
		}
		slog.Info("server stopped listening")
	}()

	return httpServer
}
