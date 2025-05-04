#!/bin/bash

rm -rf dist
mkdir -p dist/cicd-utilities-amd64
CGO_ENABLED=0 go build -ldflags="-s -w" -o ./dist/cicd-utilities-amd64/cicd-utilities ./cmd/cicd-utilitiies
