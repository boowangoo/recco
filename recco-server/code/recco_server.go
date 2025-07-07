package main

import (
    "fmt"
    "net/http"
    "log"
    "encoding/json"
    "bytes"
    "os"
)

type ReccoHost struct {
    ip          string
    embed_port  string
    db_port     string
}

type QdrantCollectionExistsResponse struct {
    Result struct {
        Exists bool
    }
}
type QdrantTitlesQueryResponse struct {
    Result struct {
        Points []struct {
            Score   float32 `json:"score"`
            Payload struct {
                Title string `json:"title"`
            } `json:"payload"`
        } `json:"points"`
    } `json:"result"`
}

type QdrantTitlesQueryRequest struct {
    Query      []float32 `json:"query"`
    WithPayload bool      `json:"with_payload"`
}

func CheckResponse(url string, resp *http.Response, err error, w http.ResponseWriter) bool {
    if err != nil {
        log.Println("Error sending request: ", err)
        fmt.Fprintf(w, "Request Failed")
        return false
    }
    if resp == nil || resp.StatusCode != http.StatusOK {
        log.Println("Error response from server: ", resp.Status)
        log.Println("Request URL: ", url)
        fmt.Fprintf(w, "Request Failed")
        return false
    }
    return true
}

func CheckMarshalJson(data []byte, err error, w http.ResponseWriter) bool {
    if err != nil {
        log.Println("Error marshaling JSON:", err)
        fmt.Fprintf(w, "Request Failed")
    }
    return err == nil 
}

func (host ReccoHost) SearchMovieHandler(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query().Get("q")

    data := map[string]string{"inputs": q}
    json_data, err := json.Marshal(data)
    if !CheckMarshalJson(json_data, err, w) {
        return
    }
    json_reader := bytes.NewReader(json_data)

    url := fmt.Sprintf("http://%s:%s/embed", host.ip, host.embed_port)
    resp, err := http.Post(url, "application/json", json_reader)
    if !CheckResponse(url, resp, err, w) {
        return
    }
    defer resp.Body.Close()

    // Decode the response
    var emb_result [][]float32
    dec := json.NewDecoder(resp.Body)
    err = dec.Decode(&emb_result)
    if err != nil{
        log.Println(err)
        fmt.Fprintf(w, "Request Failed")
        return
    }

    // Query vector DB
    var vec_data QdrantTitlesQueryRequest
    vec_data.Query = emb_result[0]
    vec_data.WithPayload = true
    json_data, err = json.Marshal(vec_data)
    if !CheckMarshalJson(json_data, err, w) {
        return
    }
    json_vec_reader := bytes.NewReader(json_data)

    url = fmt.Sprintf("http://%s:%s/collections/movie_titles/points/query", host.ip, host.db_port)
    resp, err = http.Post(url, "application/json", json_vec_reader)
    // Check for errors in the HTTP request
    if !CheckResponse(url, resp, err, w) {
        return
    }
    defer resp.Body.Close()

    // Decode the response
    var result QdrantTitlesQueryResponse
    dec = json.NewDecoder(resp.Body)
    err = dec.Decode(&result)
    if err != nil{
        log.Println("Error decoding response from vector DB: ", err)
        fmt.Fprintf(w, "Search Failed")
        return
    }

    if len(result.Result.Points) == 0 {
        fmt.Fprintf(w, "Search failed.")
    } else {
        n := len(result.Result.Points)
        titles := make([]string, 0, n)
        for _, point := range result.Result.Points {
            if point.Payload.Title == "" {
                continue
            }
            titles = append(titles, point.Payload.Title)
        }
        titles_data, err := json.Marshal(titles)
        if !CheckMarshalJson(titles_data, err, w) {
            return
        }
        fmt.Fprintf(w, string(titles_data))
    }
}

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
        listen_err := http.ListenAndServe(":80", nil)

        if listen_err != nil {
            log.Println("listen err", listen_err)
        }
    }
    log.Println("Server stopped")
}
