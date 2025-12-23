.PHONY: run build test clean

# Run the application
run:
	go run .

# Build the binary
build:
	go build -o bin/url-shortener .

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Format code
fmt:
	go fmt ./...

# Run with race detection
run-race:
	go run -race .
