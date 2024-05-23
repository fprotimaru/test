protoc-gen:
	  protoc -I="./api/proto" \
        --go_out="." \
        --go_opt=Mapi.proto="./internal/transport/grpc_contract" \
        --go-grpc_out=require_unimplemented_servers=false:"." \
        --go-grpc_opt=Mapi.proto="./internal/transport/grpc_contract" \
        api/proto/api.proto

build: protoc-gen
	mkdir -p build
	go build -ldflags="-w -s" -o build/bin cmd/main.go

test:
	go test .

docker-build: protoc-gen
	docker build .

run: protoc-gen
	go run cmd/main.go

lint:
	golangci-lint run