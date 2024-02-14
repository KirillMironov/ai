test:
	go test ./...

generate:
	go generate

lint:
	golangci-lint run
