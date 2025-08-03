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

	"github.com/rohansinghprogrammer/sudents-api/internals/config"
	"github.com/rohansinghprogrammer/sudents-api/internals/http/handlers/student"
)

func main() {
	// Load Confg
	cfg := config.MustLoadConfig()
	// Setup DB

	// Setup Routes
	router := http.NewServeMux()

	router.HandleFunc("POST /students", student.New())

	// Listen Server
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	log.Printf("Starting server on %s", cfg.Address)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to start server: %s", err.Error())
		}
	}()

	<- done

	log.Println("Shutting down server...")

	ctx , cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}
}
