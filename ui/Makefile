generate:
	protoc --dart_out=grpc:lib/api --proto_path=../api ai.proto
	protoc --dart_out=grpc:lib/api --proto_path=$(PROTO_DIR) google/protobuf/timestamp.proto
	protoc --dart_out=grpc:lib/api --proto_path=$(PROTO_DIR) google/protobuf/empty.proto
