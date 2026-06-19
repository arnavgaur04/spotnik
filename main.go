package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"spotnik/config"
	"spotnik/handlers"
	"syscall"
	"time"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	addr := fmt.Sprintf(":%d", config.Current.Port)
	fmt.Println("Starting Server...")

	http.HandleFunc("/chat", handlers.ChatHandler)

	server := &http.Server{Addr: addr}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("Server running on http://localhost%s\n", addr)
		if err := server.ListenAndServe(); err != nil {
			log.Println("Server stopped:", err)
		}
	}()

	<-stop
	fmt.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Current.ShutdownTimeout)*time.Second)
	defer cancel()

	server.Shutdown(ctx)
	fmt.Println("Done. Goodbye!")
}
