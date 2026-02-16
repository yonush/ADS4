@echo off
setlocal
cd /D %~dp0

set GO111MODULE=on
set GOFLAGS=-mod=mod
::go mod download
:: strip debug info during build
::go build -ldflags="-s -w"  -o importer.exe -v cmd/importer/main.go
go run -ldflags="-s -w" cmd/importer/main.go
