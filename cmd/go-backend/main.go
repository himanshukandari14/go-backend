package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/himanshukandari14/go-server/internal/config"
	"github.com/himanshukandari14/go-server/internal/http/handlers/student"
	"github.com/himanshukandari14/go-server/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()
	//db setups
	storage, err := sqlite.New(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initilaized", slog.String("env", cfg.Env), slog.String("version", "1,0"))

	//setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}",student.GetById(storage))
	router.HandleFunc("GET /api/students",student.GetList(storage))

	//setup srver
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {

		err := server.ListenAndServe()

		if err != nil {
			log.Fatal("failed tp start server")
		}
	}()

	<-done //we will be blocked here

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")

}
