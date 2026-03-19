@echo off
setlocal
cd /D %~dp0
FOR /F "TOKENS=1 eol=/ DELIMS=/ " %%A IN ('DATE/T') DO SET dd=%%A
FOR /F "TOKENS=1,2 eol=/ DELIMS=/ " %%A IN ('DATE/T') DO SET mm=%%B
FOR /F "TOKENS=1,2,3 eol=/ DELIMS=/ " %%A IN ('DATE/T') DO SET yyyy=%%C

set BUILDDATE=%yyyy%%mm%%dd%
set GO111MODULE=on
set GOFLAGS=-mod=mod
set GOOS=windows
set GOARCH=amd64

::set GOOS=linux
::set GOARCH=arm64

go mod download
:: strip debug info during build
go build -tags "" -ldflags="-s -w -X main.Version=1.0.0 -X main.BuildTime=$BUILDDATE%" -o ads.exe -v cmd/ads/main.go
