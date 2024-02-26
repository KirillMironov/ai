package ai

// grpc
//go:generate protoc --go_out=. --go_opt=module=$MODULE --go-grpc_out=. --go-grpc_opt=module=$MODULE ./api/llm.proto
//go:generate protoc --go_out=. --go_opt=module=$MODULE --go-grpc_out=. --go-grpc_opt=module=$MODULE ./api/ai.proto

// sqlc
//go:generate sqlc generate

// mock
//go:generate moq -pkg mock -skip-ensure -out ./internal/mock/authenticator.go ./internal/service usersStorage:UsersStorage tokenManager:TokenManager
