build:
	@go mod tidy
	@mkdir -p bin
	@echo "Building the application..."
	@go build -o bin/api
run: build
	@echo "Running the application..."
	@./bin/api