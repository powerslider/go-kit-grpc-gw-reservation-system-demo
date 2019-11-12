GRPC_GW_GOOGLEAPIS=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
GRPC_GW_SWAGGER=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway
GOGO_PROTOBUF=${GOPATH}/src/github.com/gogo/protobuf

all: clean test doc build

doc:
	@echo ">>> Generate Swagger API Documentation..."
	swag init --generalInfo cmd/reservation/main.go

build:
	@echo ">>> Building Application..."
	go build -o bin/reservations cmd/reservation/main.go

test:
	@echo ">>> Running Unit Tests..."
	go test -race ./...

cover-test:
	@echo ">>> Running Tests with Coverage..."
	go test -race ./... -coverprofile=coverage.txt -covermode=atomic

clean:
	@echo ">>> Removing binaries..."
	@rm -rf bin/*

generate:
	protoc \
		-I proto \
		-I ${GRPC_GW_GOOGLEAPIS} \
		-I ${GRPC_GW_SWAGGER} \
		-I ${GOGO_PROTOBUF} \
		--proto_path=${HOME}/go/src/github.com/powerslider/go-kit-grpc-reservation-system-demo \
		--gogo_out=plugins=grpc,paths=source_relative:./proto \
		--grpc-gateway_out=allow_patch_feature=false,Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types,logtostderr=true:./proto \
		--swagger_out=./docs \
		proto/*.proto
