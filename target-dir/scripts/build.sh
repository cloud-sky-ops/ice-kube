#!/bin/bash
set -e

OUTPUT_DIR="bin"
APP_NAME="kube-cleanup"

mkdir -p "$OUTPUT_DIR"
go build -o "$OUTPUT_DIR/$APP_NAME" ./cmd/main.go

echo "Build completed. Binary is located at $OUTPUT_DIR/$APP_NAME"
