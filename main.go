package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	channels "github.com/janithT/webpage-analyzer/channel"
	"github.com/janithT/webpage-analyzer/config"
	"github.com/janithT/webpage-analyzer/engine"
)

func main() {
	// Get the app configuration from app.yaml
	conf := config.GetAppConfig()

	// Start thread pool with 10 workers = 10 set to app.yaml
	channels.InitializetPageUrlWorkerThreadPool(conf.ThreadCount)

	// Gin router and server
	router := engine.NewRouter()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%v", conf.ServicePort),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		log.Printf("Server is running on port %v. Visit: http://localhost:%v/", conf.ServicePort, conf.ServicePort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on port %v: %v\n", conf.ServicePort, err)
		}
	}()

	// Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	log.Println("Shutting down server...")

	// Timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting gracefully.")

}
