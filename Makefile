PB_REL := https://github.com/protocolbuffers/protobuf/releases
PBGO_REL := https://github.com/protocolbuffers/protobuf-go/releases
GRPC_WEB_REL := https://github.com/grpc/grpc-web/releases
PROTOC_VERSION := 3.14.0
PROTOC_GO_GEN_VERSION := 1.25.0
OS := $(shell uname | tr '[:upper:]' '[:lower:]')
export GOBIN := $(shell pwd)/bin
export GOPATH := $(shell go env GOPATH)

PROTOC_OS=$(OS)
protoc-%-x86_64.zip:
	curl -s -LO $(PB_REL)/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(PROTOC_OS)-x86_64.zip

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
protoc-go: install-proto-tools
	$(GOBIN)/protoc \
		-I proto \
		-I /usr/local/include \
		--go_out='module=github.com/gc-plazas/kv-store:.' \
		--go-grpc_out='module=github.com/gc-plazas/kv-store:.' \
		$(PROTO_FILES)

server: ## Run server
	GRPC_GO_LOG_SEVERITY_LEVEL=debug GRPC_GO_LOG_VERBOSITY_LEVEL=2 go run cmd/server/server.go

client: ## Run client
	GRPC_GO_LOG_SEVERITY_LEVEL=debug GRPC_GO_LOG_VERBOSITY_LEVEL=2 go run cmd/client/client.go localhost:1338 get oslo
