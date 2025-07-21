package main

import (
    "fmt"
    "net/http"
    "log"
    "math"
)

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

func addVectors(a, b []float32) []float32 {
    result := make([]float32, len(a))
    for i := range a {
        result[i] = a[i] + b[i]
    }
    return result
}

func scaleVector(vector []float32, scale float32) []float32 {
    result := make([]float32, len(vector))
    for i, v := range vector {
        result[i] = v * scale
    }
    return result
}

func normalizeVector(vector []float32) []float32 {
    var magnitude float32
    for _, v := range vector {
        magnitude += v * v
    }
    magnitude = float32(math.Sqrt(float64(magnitude)))
    
    if magnitude == 0 {
        return vector
    }
    
    result := make([]float32, len(vector))
    for i, v := range vector {
        result[i] = v / magnitude
    }
    return result
}

func isZeroVector(vector []float32) bool {
    for _, v := range vector {
        if v != 0 {
            return false
        }
    }
    return true
}