FROM golang:1.16-alpine AS builder
WORKDIR /src
COPY . .
RUN apk add curl unzip protobuf &&  go install google.golang.org/protobuf/cmd/protoc-gen-go && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
RUN protoc -I proto --go_out='module=github.com/gc-plazas/kv-store:.' \
     --go-grpc_out='module=github.com/gc-plazas/kv-store:.'  \
    ./proto/kv/server/server_service.proto
RUN go build -o client ./go/cmd/client/client.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /src/client ./
ENTRYPOINT ["./client"]
CMD ["localhost:1338", "get", "oslo"]
