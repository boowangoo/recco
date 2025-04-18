package main

import "github.com/anush008/fastembed-go"

func EmbedText(documents []string, batch_size int) ([][]float32, error) {
	model, err := fastembed.NewFlagEmbedding(nil)
	if err != nil {
		return nil, err
	}
	defer model.Destroy()
	embeddings, err := model.Embed(documents, batch_size)
	return embeddings, err
}