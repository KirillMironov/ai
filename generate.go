package ai

//go:generate protoc --go_out=. --go_opt=module=$MODULE --go-grpc_out=. --go-grpc_opt=module=$MODULE ./api/llm.proto
//go:generate protoc --go_out=. --go_opt=module=$MODULE --go-grpc_out=. --go-grpc_opt=module=$MODULE ./api/ai.proto
//go:generate sqlc generate
