#!/bin/bash

set -e

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

cd "$SCRIPT_DIR/.."

echo "Tide dependencies up"
go mod tidy

echo "Compiling the cloud native buildpack"
rm -rf ./bin
GOOS=linux go build -ldflags="-s -w" -o ./bin/run ./run/main.go
ln -sf run ./bin/detect
ln -sf run ./bin/build