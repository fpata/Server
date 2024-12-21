package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"clinic_server/config"
)

func main() {
	// Load configuration once
	cfg := config.GetConfiguration()
	urlp := cfg.Server.ServerUrl + ":" + cfg.Server.Port

	// Initialize the router
	router := SetupRouter()

	// Use context for request handling
	srv := &http.Server{
		Addr:    urlp,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
