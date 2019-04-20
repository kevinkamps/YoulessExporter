#!/usr/bin/env bash

echo "$(date): Starting"

export GO111MODULE=on

mkdir bin
rm -rf bin/*

echo "$(date): Building for Linux"
GOOS=linux GOARCH=386 go build -o ./bin/youless_linux_386 ./main.go
GOOS=linux GOARCH=amd64 go build -o ./bin/youless_linux_amd64 ./main.go

echo "$(date): Building for Linux arm"
GOOS=linux GOARCH=arm go build -o ./bin/youless_linux_arm ./main.go

echo "$(date): Done"