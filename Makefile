javadoc2md:
	go build -o javadoc2md cmd/javadoc2md/main.go

test:
	go test ./...
