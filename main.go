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
    "spotnik/handlers"
    "spotnik/database"
)

func main() {
    fmt.Println("Starting Server...")

    err := database.Connect("postgres://arnavgaur@localhost:5432/spotnik?search_path=spotnik&sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("/chat", handlers.ChatHandler)
    
    server := &http.Server{Addr: ":8080"}

    // Listen for Ctrl+C (just like our CLI!)
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

    // Start server in background
    go func() {
        fmt.Println("Server running on http://localhost:8080")
        if err := server.ListenAndServe(); err != nil {
            log.Println("Server stopped:", err)
        }
    }()

    // Wait for Ctrl+C
    <-stop
    fmt.Println("\nShutting down gracefully...")

    // Give ongoing requests 5 seconds to finish
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    server.Shutdown(ctx)
    fmt.Println("Done. Goodbye!")
}
