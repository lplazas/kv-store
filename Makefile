PB_REL := https://github.com/protocolbuffers/protobuf/releases
PBGO_REL := https://github.com/protocolbuffers/protobuf-go/releases
GRPC_WEB_REL := https://github.com/grpc/grpc-web/releases
PROTOC_VERSION := 3.14.0
PROTOC_GO_GEN_VERSION := 1.25.0
OS := $(shell uname | tr '[:upper:]' '[:lower:]')
export GOBIN := $(shell pwd)/bin
export GOPATH := $(shell go env GOPATH)

#TODO have both use same base image
build-client:
	docker build . -f client.Dockerfile --tag kv-store-client:latest

build-server:
	docker build . -f server.Dockerfile --tag kv-store-server:latest

PROTOC_OS=$(OS)
protoc-%-x86_64.zip:
	curl -s -LO $(PB_REL)/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(PROTOC_OS)-x86_64.zip
	echo "-s -LO $(PB_REL)/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(PROTOC_OS)-x86_64.zip"

/usr/local/include/google: protoc-%-x86_64.zip
	unzip -qo protoc-$(PROTOC_VERSION)-$(PROTOC_OS)-x86_64.zip -d /usr/local 'include/*'

PROTOC_OS=$(OS)
ifeq ($(OS),darwin)
	override PROTOC_OS=osx
endif
bin/protoc: protoc-%-x86_64.zip /usr/local/include/google
	unzip -qo protoc-$(PROTOC_VERSION)-$(PROTOC_OS)-x86_64.zip -d . bin/protoc
	rm -f protoc-$(PROTOC_VERSION)-$(PROTOC_OS)-x86_64.zip
	echo "Installed bin/protoc"

bin/protoc-gen-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	echo "Installed bin/protoc-gen-go"

bin/protoc-gen-go-grpc:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	echo "Installed bin/protoc-gen-go-grpc"

install-proto-tools: bin/protoc  bin/protoc-gen-go bin/protoc-gen-go-grpc

.PHONY: protoc-go
PROTO_FILES=$(shell find . -path '*.proto')
protoc-go: install-proto-tools # Compile .proto files
	$(GOBIN)/protoc \
		-I proto \
		-I /usr/local/include \
		--go_out='module=github.com/gc-plazas/kv-store:.' \
		--go-grpc_out='module=github.com/gc-plazas/kv-store:.' \
		$(PROTO_FILES)

build: protoc-go
	go build -o ./server ./go/cmd/server/server.go
	go build -o ./client ./go/cmd/client/client.go
