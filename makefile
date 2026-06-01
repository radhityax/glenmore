APP_NAME = glenmore
.PHONY: all build run clean test dev
all: build

build:
	go build -o $(APP_NAME) .

run: build
	./$(APP_NAME)

dev: build
	@./$(APP_NAME)

clean:
	rm -f $(APP_NAME)
	rm -rf data/

test:
	go test ./...

check:
	go vet ./...
