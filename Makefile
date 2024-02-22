MODULE := github.com/KirillMironov/ai

test:
	go test ./...

generate:
	MODULE=$(MODULE) go generate

lint:
	golangci-lint run
	nilaway -include-pkgs="$(MODULE)" ./...

tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.1
	go install go.uber.org/nilaway/cmd/nilaway@latest
