package main

import (
    "fmt"
    "net/http"
    "log"
    "io"
)

func main() {
	fmt.Println("Starting server")
    http.HandleFunc("/", HelloServer)
    http.HandleFunc("/google-query", GoogleQueryHandler)
    listen_err := http.ListenAndServe(":8080", nil)
    if listen_err != nil {
        fmt.Println("listen err", listen_err)
    }
	fmt.Println("Server stopped")

}

func HelloServer(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, world")
}

func GoogleQueryHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    query_url := fmt.Sprintf("https://google.com/search?q=%s", query)
    resp, err := http.Get(query_url)
    if err != nil {
        log.Println(err)
        fmt.Fprintf(w, "Request Failed")
        return
    } 
    // fmt.Fprintf(w, resp)
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Println("Query to Google failed: %d", resp.StatusCode)
        fmt.Fprintf(w, "Request Failed")
        return
    }

    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("Body Read Failed: ", err)
        fmt.Fprintf(w, "Request Failed")
        return
    }
    fmt.Fprintf(w, string(bodyBytes))
}