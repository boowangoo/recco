package main

import (
    "fmt"
    "net/http"
    "log"
    "io"
    "encoding/json"
)

func main() {
	fmt.Println("Starting server")

    // check if database is running and data is loaded
    CheckDB()

    http.HandleFunc("/", HelloServer)
    http.HandleFunc("/google-query", GoogleQueryHandler)
    http.HandleFunc("/qdrant-collection-exists", QdrantCollectionExitsHandler)
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

func QdrantCollectionExitsHandler(w http.ResponseWriter, r *http.Request) {
    collection := r.URL.Query().Get("collection")
    url := fmt.Sprintf("http://recco-db:6333/collections/%s/exists", collection)
    resp, err := http.Get(url)
    if err != nil {
        log.Println(err)
        fmt.Fprintf(w, "Request Failed")
        return
    }
    if resp.StatusCode != 200 {
        log.Println(resp.Status)
        fmt.Fprintf(w, "Request Failed")
        return
    }

    type exists_resp struct {
        Result struct {
            Exists bool
        }
    }
    var result exists_resp
    dec := json.NewDecoder(resp.Body)
    err = dec.Decode(&result)
    if err != nil{
        log.Println(err)
        fmt.Fprintf(w, "Request Failed")
        return
    }

    if !result.Result.Exists {
        fmt.Fprintf(w, "The collection %s does not exist", collection)
    } else {
        fmt.Fprintf(w, "The collection %s exist", collection)
    }
}
