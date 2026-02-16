@echo off
setlocal
cd /D %~dp0

set GO111MODULE=on
set GOFLAGS=-mod=mod
go mod download
:: strip debug info during build
go build -tags "" -ldflags="-s -w" -o ads.exe -v cmd/ads/main.go
