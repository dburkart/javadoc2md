javadoc2md:
	go build -o javadoc2md cmd/javadoc2md/main.go

test:
	@echo "Running unit tests..."
	go test ./...
	@echo
	@echo "Running end-2-end tests..."
	bash ./tests/e2e/run.sh
