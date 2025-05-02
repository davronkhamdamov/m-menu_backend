#!/bin/sh
GOOS=windows GOARCH=amd64 go build -ldflags '-H windowsgui -s -w' -o bin/restaurant-windows-amd64.exe manage.go
# GOOS=windows    GOARCH=arm64    go build -ldflags '-s -w' -o bin/restaurant-windows-arm64.exe  manage.go
# GOOS=darwin    GOARCH=arm64    go build -ldflags '-s -w' -o bin/restaurant-darwin-arm64  manage.go
