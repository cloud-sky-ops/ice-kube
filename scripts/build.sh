#!/bin/bash
set -e

OUTPUT_DIR="bin"
APP_NAME="ice-kube"

mkdir -p "$OUTPUT_DIR"
go build -o "$OUTPUT_DIR/$APP_NAME" ./main.go

echo "Build completed. Binary is located at $OUTPUT_DIR/$APP_NAME"
echo "Moving to /usr/local/bin/"

cp "$OUTPUT_DIR/$APP_NAME" "/usr/local/bin/"

rm -rf $OUTPUT_DIR/$APP_NAME
ice-kube --help
