package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	server2 "github.com/gc-plazas/kv-store/go/external/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net/url"
	"os"
)

const (
	getCommand            = "get"
	putCommand            = "put"
	incorrectUsageMessage = "incorrect usage \n ./client host get key \n ./client host put key value"
)

func main() {
	args := os.Args[1:]

	if len(args) < 3 || len(args) > 4 {
		fmt.Println(incorrectUsageMessage)
		return
	}

	host := args[0]
	url, err := url.ParseRequestURI(host)
	if err != nil {
		fmt.Println(incorrectUsageMessage)
		return
	}

	command := args[1]
	if len(args) == 3 && command != getCommand {
		fmt.Println(incorrectUsageMessage)
		return
	}
	if len(args) == 4 && command != putCommand {
		fmt.Println(incorrectUsageMessage)
		return
	}

	conn, connErr := NewGrpcConn(url.String(), true)
	if connErr != nil {
		fmt.Println("failed connecting to the specified host, err:", connErr.Error())
		return
	}

	ctx := context.Background()

	serverClient := server2.NewServerServiceClient(conn)

	if command == getCommand {
		getKey := args[2]
		response, getErr := serverClient.GetValue(ctx, &server2.GetValueRequest{Key: getKey})
		if getErr != nil {
			fmt.Println("error: ", getErr.Error())
		} else {
			fmt.Println("value: ", response.Value)
		}
	} else if command == putCommand {
		putKey := args[2]
		putValue := args[3]
		_, putErr := serverClient.PutValue(ctx, &server2.PutValueRequest{
			Key:   putKey,
			Value: putValue,
		})
		if putErr != nil {
			fmt.Println("error: ", putErr.Error())
		}
	}
}

func NewGrpcConn(host string, insecure bool) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if host != "" {
		opts = append(opts, grpc.WithAuthority(host))
	}

	if insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs:    systemRoots,
			MinVersion: tls.VersionTLS12,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))

		host = host + ":443"
	}

	return grpc.Dial(host, opts...)
}

func OptionalEnvGet(key, defaultValue string) string {
	value, isSet := os.LookupEnv(key)
	if !isSet {
		return defaultValue
	}
	return value
}
