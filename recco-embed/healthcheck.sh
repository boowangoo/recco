#!/bin/bash

RESPONSE=$(curl -i 0.0.0.0:80/embed -H 'Content-Type: application/json' -d '{"inputs":"test"}' | awk '{print $2; exit}')
if [[ "$RESPONSE" == "200" ]]; then
    # The service is healthy
    exit 0
fi
# The service is unhealthy
exit 1
