package main

import (
	"fmt"
	"github.com/gc-plazas/kv-store/external/server"
	"github.com/gc-plazas/kv-store/internal"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	c, err := internal.NewCluster(5, 10)
	if err != nil {
		panic(err)
	}
	clusterServer := server.NewClusterServer(c)
	g := grpc.NewServer()
	server.RegisterServerServiceServer(g, clusterServer)

	address := fmt.Sprintf(":%v", OptionalEnvGet("PORT", "1338"))
	listen, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	fmt.Println("Server API listening to", address)

	if _err := g.Serve(listen); _err != nil {
		panic(_err)
	}
}

func OptionalEnvGet(key, defaultValue string) string {
	value, isSet := os.LookupEnv(key)
	if !isSet {
		return defaultValue
	}
	return value
}
