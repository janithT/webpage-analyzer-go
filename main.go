package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/janithT/webpage-analyzer/config"
	"github.com/janithT/webpage-analyzer/engine"
)

func main() {
	conf := config.GetAppConfig()

	router := engine.NewRouter() // Gin router - start here

	server := &http.Server{
		Addr:         fmt.Sprintf(":%v", conf.ServicePort),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Run the server in a goroutine so it doesn't block.
	go func() {
		log.Printf("Server is running on port %v. Visit: http://localhost:%v/", conf.ServicePort, conf.ServicePort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on port %v: %v\n", conf.ServicePort, err)
		}
	}()

	// Shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	log.Println("Shutting down server...")

	// Context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting gracefully.")

}
