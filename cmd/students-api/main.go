package main

import (
	"log"
	"net/http"

	"github.com/rohansinghprogrammer/sudents-api/internals/config"
)

func main()  {
	// Load Confg
	cfg := config.MustLoadConfig()
	// Setup DB

	// Setup Routes
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	
	// Listen Server
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	log.Printf("Starting server on %s", cfg.Address)
	// Start Server
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}
}