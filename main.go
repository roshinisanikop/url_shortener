package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	store := NewURLStore()
	handler := NewHandler(store)

	http.HandleFunc("/", handler.HandleRedirect)
	http.HandleFunc("/shorten", handler.HandleShorten)
	http.HandleFunc("/api/urls", handler.HandleListURLs)

	port := ":8080"
	fmt.Printf("URL Shortener running on http://localhost%s\n", port)
	fmt.Println("Endpoints:")
	fmt.Println("  POST /shorten - Create a short URL")
	fmt.Println("  GET  /{code}  - Redirect to original URL")
	fmt.Println("  GET  /api/urls - List all URLs")

	srv := &http.Server{
		Addr:         port,
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
