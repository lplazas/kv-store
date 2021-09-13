FROM golang:1.16-alpine AS builder
WORKDIR /src
COPY . .
ENV PROTOC_VERSION=3.14.0
ENV PB_REL=https://github.com/protocolbuffers/protobuf/releases
RUN export PROTO_FILES=$(find . -path '*.proto') && apk add curl unzip protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
RUN protoc -I proto --go_out='module=github.com/gc-plazas/kv-store:.' \
     --go-grpc_out='module=github.com/gc-plazas/kv-store:.'  \
    ./proto/kv/server/server_service.proto
RUN go build -o server ./go/cmd/server/server.go && go build -o client ./go/cmd/client/client.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /src/server ./
COPY --from=builder /src/client ./
CMD ["./server"]
