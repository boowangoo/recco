package main

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

type QdrantMoviesPayload struct {
    Title string `json:"title"`
    Genres []string `json:"genres"`
    Year int `json:"year"`
    AverageRating float32 `json:"average_rating"`
}

type QdrantTitlesQueryResponse struct {
    Result struct {
        Points []struct {
            Id int `json:"id"`
            Payload QdrantMoviesPayload `json:"payload"`
        } `json:"points"`
    } `json:"result"`
}

type QdrantTitlesQueryRequest struct {
    Query []float32 `json:"query"`
    Using string `json:"using"`
    Limit int `json:"limit"`
    WithPayload bool `json:"with_payload"`
}

type RecommendRequest struct {
    Ids     []int `json:"ids"`
    Ratings []int `json:"ratings"`
}

type QdrantPointsRequest struct {
    Ids         []int  `json:"ids"`
    WithPayload bool   `json:"with_payload"`
    WithVector  []string `json:"with_vector"`
}

type QdrantVectorResponse struct {
    Result []struct {
        Id     int                    `json:"id"`
        Vector map[string][]float32   `json:"vector"`
    } `json:"result"`
}

type QdrantRatingResponse struct {
    Result struct {
        Id      int `json:"id"`
        Payload struct {
            Original []float32 `json:"original"`
            Refit    []float32 `json:"refit"`
        } `json:"payload"`
    } `json:"result"`
}

type QdrantRecommendQuery struct {
    Query struct {
        Recommend struct {
            Positive [][]float32 `json:"positive,omitempty"`
            Negative [][]float32 `json:"negative,omitempty"`
        } `json:"recommend"`
    } `json:"query"`
    Using       string `json:"using"`
    Limit       int    `json:"limit"`
    WithVector  bool   `json:"with_vector"`
    WithPayload bool   `json:"with_payload"`
}
