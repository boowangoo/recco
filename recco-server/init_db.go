package main

import (
    // "fmt"
    "net/http"
    "log"
    "encoding/json"
    "encoding/csv"
    "bytes"
    "io"
    "strconv"
    "os"
    // "os/exec"

    // "github.com/sugarme/tokenizer/pretrained"
)

func CheckDB() {
    url := "http://recco-db:6333/collections/movie-titles/exists"
    resp, err := http.Get(url)
    if err != nil {
        log.Println(err)
        return
    }
    if resp.StatusCode != 200 {
        log.Println(resp.Status)
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
        return
    }

    if !result.Result.Exists {
        log.Println("The collection \"movie-titles\" does not exist.")
        CreateMovieTitlesCollection()
    }
}

func CreateMovieTitlesCollection() {
    // Create collection
    url := "http://recco-db:6333/collections/movie-titles/"
    data := map[string]interface{}{
        "dense-vector-name": map[string]interface{}{
            "size": 1536,
            "distance": "Cosine",
        },
    }
    jsonData, err := json.Marshal(data)
    if err != nil {
        log.Println(err)
        return
    }

    req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
        return
	}

    log.Println("Creating collection")
    client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Println(err)
        return
	}

    LoadMovieTitles()
}


func LoadMovieTitles() {
    csv_file, err := os.Open("/recco-dataset/movies_metadata.csv")
    if err != nil {
		log.Println(err)
        return
    }
    defer csv_file.Close()

    csv_reader := csv.NewReader(csv_file)
    header, err := csv_reader.Read()
    if err != nil {
		log.Printf("Cannot read header: %s", err)
        return
    }
    id_idx := -1
    title_idx := -1
    for i, col := range header {
        if col == "id" {
            id_idx = i
        } else if col == "title" {
            title_idx = i
        }
    }
    if id_idx == -1 || title_idx == -1 {
		log.Println("Missing column(s)")
        return
    }

    var ids []uint64
    var titles []string
    var payloads []map[string]string

    BATCH_SIZE := 256

    for {
        row, err := csv_reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return
        }

        id, err := strconv.ParseUint(row[id_idx], 10, 64)
        if err != nil {
            log.Printf("Invalid ID: %s", err)
            return
        }
        title := row[title_idx]

        ids = append(ids, id)
        titles = append(titles, title)
        payloads = append(payloads, map[string]string { "title": title } )

        if len(ids) >= BATCH_SIZE {
		    log.Println("Creating vectors")
            vectors, err := EmbedText(titles, BATCH_SIZE)
            if err != nil {
                log.Println(err)
                return
            }
            resp, err := BatchUpsertPointsDB(ids, vectors, payloads)
            if err != nil {
                log.Println(err)
                return
            }

            if resp.StatusCode != http.StatusOK {
                log.Println("Update to DB failed: %d", resp.StatusCode)
                return
            }

            ids, titles, payloads = nil, nil, nil
        }
    }
}

func BatchUpsertPointsDB(ids []uint64, vectors [][]float32, payloads []map[string]string) (*http.Response, error) {
    url := "http://recco-db:6333/collections/movie-titles/points"
    data := map[string]interface{}{
        "batch": map[string]interface{}{
            "ids": ids,
            "vectors": vectors,
            "payloads": payloads,
        },
    }
    jsonData, err := json.Marshal(data)
    if err != nil {
        log.Println(err)
        return nil, err
    }

    req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
        return nil, err
	}

    log.Println("Creating collection")
    client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
        return nil, err
	}

    return resp, err
}