# Introduction

This project using golang languange & database PostgreSQL

## Installation

- Replace and change name config.yml.example to config.yml
- Setup database configuration
- Run migration with command `goose package`
- Run service with command `go run main.go start`, default port using 8081
- Import API Docs on Postman app

## Migrate Command

- `go run main.go migrate up`
- `go run main.go migrate down`
- `go run main.go migrate create table_name`
- `go run main.go seed all fresh`
- `go run main.go seed seed_name`

## Docker

- Run command `sh bash.sh`
