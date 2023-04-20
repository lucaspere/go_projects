compile:
	protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  ./api/v1/log.proto

test:
	go test -race ./...
