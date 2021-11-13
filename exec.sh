#!/usr/bin/env bash
docker build -t auth0-demo-api-go .
docker run -it --rm -p 8080:8080 --name auth0-demo-api-go auth0-demo-api-go
