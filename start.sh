#!/bin/bash

# Download text embedding model for recco-embed
if [[ "$@" == *"--download-embed"* ]]; then
    pushd ./recco-embed/data
    # Requires git-lfs to be installed
    git clone https://huggingface.co/BAAI/bge-large-en-v1.5
    popd
fi

docker compose --env-file ./recco.env up --build -d

# Load dataset into recco-db
if [[ "$@" == *"--download-db"* ]]; then
    python db_load.py
fi