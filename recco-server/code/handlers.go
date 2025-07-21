package main

import (
    "fmt"
    "net/http"
    "log"
    "encoding/json"
    "bytes"
)

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
    vec_data.Using = "title"
    vec_data.Limit = 5
    vec_data.WithPayload = true
    json_data, err = json.Marshal(vec_data)
    if !CheckMarshalJson(json_data, err, w) {
        return
    }
    json_vec_reader := bytes.NewReader(json_data)

    url = fmt.Sprintf("http://%s:%s/collections/movies/points/query", host.ip, host.db_port)
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
        payloads := make([]QdrantMoviesPayload, 0, n)
        for _, point := range result.Result.Points {
            if point.Payload.Title == "" {
                continue
            }
            payloads = append(payloads, point.Payload)
        }
        payloads_data, err := json.Marshal(payloads)
        if !CheckMarshalJson(payloads_data, err, w) {
            return
        }
        fmt.Fprintf(w, string(payloads_data))
    }
}

func (host ReccoHost) RecommendMovieHandler(w http.ResponseWriter, r *http.Request) {
    var req RecommendRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        log.Println("Error parsing request body:", err)
        fmt.Fprintf(w, "Request Failed")
        return
    }
	
    if len(req.Ids) != len(req.Ratings) {
		log.Println("Ids and ratings arrays must have the same length")
		fmt.Fprintf(w, "Request Failed")
        return
    }

    // Get feature vectors for the movie IDs
    vectorsReq := QdrantPointsRequest{
        Ids:         req.Ids,
        WithPayload: false,
        WithVector:  []string{"features"},
    }

    json_data, err := json.Marshal(vectorsReq)
    if !CheckMarshalJson(json_data, err, w) {
        return
    }

    url := fmt.Sprintf("http://%s:%s/collections/movies/points", host.ip, host.db_port)
    resp, err := http.Post(url, "application/json", bytes.NewReader(json_data))
    if !CheckResponse(url, resp, err, w) {
        return
    }
    defer resp.Body.Close()

    var vectorsResult QdrantVectorResponse
    err = json.NewDecoder(resp.Body).Decode(&vectorsResult)
    if err != nil {
        log.Println("Error decoding vectors response:", err)
        fmt.Fprintf(w, "Request Failed")
        return
    }

    vectorDim := 0
	if len(vectorsResult.Result) > 0 {
		features, exists := vectorsResult.Result[0].Vector["features"]
		if exists {
			vectorDim = len(features)
		}
	}

    // Get ratings table (weights)
    url = fmt.Sprintf("http://%s:%s/collections/ratings/points/0", host.ip, host.db_port)
    resp, err = http.Get(url)
    if !CheckResponse(url, resp, err, w) {
        return
    }
    defer resp.Body.Close()

    var ratingsResult QdrantRatingResponse
    err = json.NewDecoder(resp.Body).Decode(&ratingsResult)
    if err != nil {
        log.Println("Error decoding ratings response:", err)
        fmt.Fprintf(w, "Request Failed")
        return
    }

    weights := ratingsResult.Result.Payload.Refit

    if vectorDim == 0 {
		log.Println("No feature vectors found to infer dimension")
        fmt.Fprintf(w, "Request Failed")
        return
    }

	// Calculate positive and negative vectors
	var positive []float32
	var negative []float32

    positive = make([]float32, vectorDim)
    negative = make([]float32, vectorDim)

    for i, id := range req.Ids {
        rating := req.Ratings[i]

        // Validate rating index
        if rating < 0 || rating >= len(weights) {
            log.Printf("Invalid rating index %d for ID %d", rating, id)
        	fmt.Fprintf(w, "Request Failed")
            continue
        }

        weight := weights[rating]
        
        // Find the corresponding vector result for this ID
        var features []float32
        var exists bool
        for _, point := range vectorsResult.Result {
            if point.Id == id {
                features, exists = point.Vector["features"]
                break
            }
        }
        
        if !exists {
            log.Printf("No features found for ID %d", id)
        	fmt.Fprintf(w, "Request Failed")
            continue
        }

        scaledFeatures := scaleVector(features, weight)
        
        if weight >= 0 {
            positive = addVectors(positive, scaledFeatures)
        } else {
            negative = addVectors(negative, scaleVector(scaledFeatures, -1))
        }
    }

    if !isZeroVector(positive) {
        positive = normalizeVector(positive)
    }
    if !isZeroVector(negative) {
        negative = normalizeVector(negative)
    }

    // Perform recommendation
    var recommendQuery QdrantRecommendQuery
    recommendQuery.Query.Recommend = struct {
        Positive [][]float32 `json:"positive,omitempty"`
        Negative [][]float32 `json:"negative,omitempty"`
    }{}
    
    if !isZeroVector(positive) {
        recommendQuery.Query.Recommend.Positive = [][]float32{positive}
    }
    if !isZeroVector(negative) {
        recommendQuery.Query.Recommend.Negative = [][]float32{negative}
    }

    recommendQuery.Using = "features"
    recommendQuery.Limit = 5
    recommendQuery.WithVector = false
    recommendQuery.WithPayload = true

    json_data, err = json.Marshal(recommendQuery)
    if !CheckMarshalJson(json_data, err, w) {
        return
    }

    url = fmt.Sprintf("http://%s:%s/collections/movies/points/query", host.ip, host.db_port)
    resp, err = http.Post(url, "application/json", bytes.NewReader(json_data))
    if !CheckResponse(url, resp, err, w) {
        return
    }
    defer resp.Body.Close()

    var result QdrantTitlesQueryResponse
    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
        log.Println("Error decoding recommendation response:", err)
        fmt.Fprintf(w, "Recommendation failed")
        return
    }

    if len(result.Result.Points) == 0 {
        fmt.Fprintf(w, "No recommendations found.")
    } else {
        n := len(result.Result.Points)
        payloads := make([]QdrantMoviesPayload, 0, n)
        for _, point := range result.Result.Points {
            if point.Payload.Title == "" {
                continue
            }
            payloads = append(payloads, point.Payload)
        }
        payloads_data, err := json.Marshal(payloads)
        if !CheckMarshalJson(payloads_data, err, w) {
            return
        }
        fmt.Fprintf(w, string(payloads_data))
    }
}
