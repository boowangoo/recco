#!/bin/bash

# Check if RECCO_IP exists in recco.env
if [ -f "./recco.env" ] && grep -q "RECCO_IP=" ./recco.env; then
    echo "Using existing host: $(grep '^RECCO_IP=' ./recco.env | cut -d'=' -f2-)"
else
    # Find the default network interface from the routing table
    DEFAULT_ITF=$(ip route | awk '/default/ {print $5}')
    if [ -z "$DEFAULT_ITF" ]; then
        echo "No default interface found. Please check your network configuration."
        exit 1
    fi
    # Find the inet address of the default interface, removing the suffix
    RECCO_IP=$(ip addr show "$DEFAULT_ITF" | awk '/inet / {print $2}' | cut -d/ -f1)
    if [ -z "$RECCO_IP" ]; then
        echo "No IP address found for interface $DEFAULT_ITF. Please check your network configuration."
        exit 1
    fi
    echo "" >> ./recco.env
    echo "RECCO_IP=$RECCO_IP" >> ./recco.env
fi

docker compose --env-file ./recco.env up --build -d
