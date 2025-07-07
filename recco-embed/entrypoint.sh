#!/bin/bash
# Download the model if it does not exist
if [ ! -f /data/bge-large-en-v1.5/model.safetensors ]; then
    echo "Downloading BGE model..."
    git clone https://huggingface.co/BAAI/bge-large-en-v1.5 /data/bge-large-en-v1.5
    echo "BGE model downloaded successfully."
else
    echo "BGE model already exists."
fi

# Start the Text Embeddings Inference service
echo "Starting Text Embeddings Inference service..."
/usr/local/bin/text-embeddings-router
