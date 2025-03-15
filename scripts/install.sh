#!/bin/bash
set -e

BIN_DIR="/usr/local/bin"
APP_NAME="ice-kube"

if [ -f "$BIN_DIR/$APP_NAME" ]; then
    echo "$APP_NAME already installed. Removing old version..."
    sudo rm -f "$BIN_DIR/$APP_NAME"
fi

sudo cp ./bin/$APP_NAME "$BIN_DIR/"
sudo chmod +x "$BIN_DIR/$APP_NAME"

echo "$APP_NAME installed successfully."
