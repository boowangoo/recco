#!/bin/bash

# Download text embedding model for recco-embed
if [[ "$@" == *"--download-embed"* ]]; then
    pushd ./recco-embed/data
    # Requires git-lfs to be installed
    git clone https://huggingface.co/BAAI/bge-large-en-v1.5
    popd
fi

# Check if RECCO_IP exists in recco.env
if [ -f "./recco.env" ] && grep -q "RECCO_IP=" ./recco.env; then
    echo "Using existing host: $RECCO_IP"
else
    # Find the default network interface from the routing table
    DEFAULT_ITF=$(ip route | awk '/default/ {print $5}')
    if [ -z "$DEFAULT_ITF" ]; then
        echo "No default interface found. Please check your network configuration."
        exit 1
    fi
    # Find the inet address of the default interface, removing the suffix
    RECCO_IP=$(ip addr show "$DEFAULT_ITF" | grep 'inet ' | awk '{print $2}' | cut -d/ -f1)
    if [ -z "$RECCO_IP" ]; then
        echo "No IP address found for interface $DEFAULT_ITF. Please check your network configuration."
        exit 1
    fi
    echo "RECCO_IP=$RECCO_IP" >> ./recco.env
fi

docker compose --env-file ./recco.env up --build -d

# Load dataset into recco-db
if [[ "$@" == *"--download-db"* ]]; then
    python db_load.py
fi