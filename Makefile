.PHONY: build run clean docker test

APP_NAME=feedvackbot
BUILD_DIR=./cmd/bot

build:
	@echo "Building $(APP_NAME)..."
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(APP_NAME) $(BUILD_DIR)
	@echo "Build complete: ./$(APP_NAME)"

run: build
	@echo "Starting $(APP_NAME)..."
	./$(APP_NAME)

clean:
	@rm -f $(APP_NAME)
	@echo "Cleaned."

docker:
	docker compose up --build -d

docker-stop:
	docker compose down

docker-logs:
	docker compose logs -f bot

test:
	go test ./... -v

tidy:
	go mod tidy
