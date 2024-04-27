set -e
docker compose up -d
docker exec api go run main.go migrate up
docker exec api go run main.go seed all fresh
