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
)


func main() {
	//load config
	cfg:=config.MustLoad()
	//db setup
	//setup router
	router :=http.NewServeMux()
	router.HandleFunc("POST /api/students",student.New())
	
	//setup srver
	server :=http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started",slog.String("address",cfg.Addr))

	done:=make(chan os.Signal, 1)

	signal.Notify(done,os.Interrupt, syscall.SIGINT, syscall.SIGTERM)


	go func(){

		err :=server.ListenAndServe()
	
		if err!=nil {
			log.Fatal("failed tp start server")
		}
	}()

	<- done //we will be blocked here 

	slog.Info("shutting down the server")

	ctx, cancel :=context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()


	err :=server.Shutdown(ctx)

	if err!=nil {
		slog.Error("Failed to shutdownserver", slog.String("error", err.Error()))
	}
	slog.Info("Server shutodown sucessfully")



}