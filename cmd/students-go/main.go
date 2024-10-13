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

	"github.com/raunak173/students-go/internal/config"
)

func main() {

	//load config
	cfg := config.MustLoad()

	//setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the student's api"))
	})

	//setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started", slog.String("address", cfg.Addr))

	//We are making a channel to receive signals
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	//we are making a go routine
	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal("failed to start the server")
		}
	}()

	//Blocked until we receive a signal from the channel
	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown the server", slog.String("error", err.Error()))
	}

	slog.Info("Server shut down successfully!")

}
