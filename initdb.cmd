@echo off
setlocal
cd /D %~dp0

:: Define variables for Goose commands
set DATA_DIR=".\data\migrations"

set DB_USER=postgres
set DB_PASSWORD=postgres
set DB_NAME=ADS4
set DB_HOST=localhost
set DB_PORT=5432

:: Run Goose up command
echo "Running goose up..."
::goose -dir %DATA_DIR% postgres "user=%DB_USER% password=%DB_PASSWORD% dbname=%DB_NAME% host=%DB_HOST% port=%DB_PORT% sslmode=disable" up
echo Removing .\data\%DB_NAME%.db
::del .\data\%DB_NAME%.db >nul 2> nul
del .\data\seed_complete >nul 2> nul
echo Running goose -dir %DATA_DIR% sqlite ".\data\%DB_NAME%.db"
goose -dir %DATA_DIR% sqlite ".\data\%DB_NAME%.db" down
::goose -dir %DATA_DIR% sqlite ".\data\%DB_NAME%.db" create init sql
goose -dir %DATA_DIR% sqlite ".\data\%DB_NAME%.db" up
echo "Done."