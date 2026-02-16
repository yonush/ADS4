@echo off
set GO111MODULE=on
set GOFLAGS=-mod=mod
::go mod download
:: strip debug info during build
go run -ldflags="-s -w"  ./cmd/ads/main.go
