package main

import (
    "fmt"
    "net/http"
    "log"
    "os"
)

func main() {
	fmt.Println("Starting server")

    args := os.Args
    // Check if command line arguments are provided
    if len(args) < 4 {
        fmt.Println("Usage: go run . <RECCO_IP> <RECCO_EMBED_PORT> <RECCO_DB_PORT>")
        os.Exit(1)
    }
    host := ReccoHost{
        ip: args[1],
        embed_port: args[2],
        db_port: args[3],
    }
    if host.ip == "" || host.embed_port == "" || host.db_port == "" {
        log.Println("Environment variables RECCO_IP, RECCO_EMBED_PORT, and RECCO_DB_PORT must be set in the recco.env file.")
    } else {
        http.HandleFunc("/search", host.SearchMovieHandler)
        http.HandleFunc("/recommend", host.RecommendMovieHandler)
        listen_err := http.ListenAndServe(":80", nil)

        if listen_err != nil {
            log.Println("listen err", listen_err)
        }
    }
    log.Println("Server stopped")
}
