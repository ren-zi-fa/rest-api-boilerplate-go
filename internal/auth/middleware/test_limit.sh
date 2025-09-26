#!/bin/bash

URL="http://localhost:8080/api/v1/posts"
for i in {1..20}; do
    echo "Request #$i"
    curl -i $URL
    echo ""
done
