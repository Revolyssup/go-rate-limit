package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/Revolyssup/go-rate-limit/examples/grpc-rate-limiter/greeter"
	"github.com/Revolyssup/go-rate-limit/pkg"
	leakybucket "github.com/Revolyssup/go-rate-limit/pkg/leaky-bucket"
	"google.golang.org/grpc"
)

func main() {
	lb := leakybucket.NewLeakyBucket(1, 2)
	flag.Parse()
	listenOn := fmt.Sprintf(":%s", flag.Arg(0))
	lis, err := net.Listen("tcp", listenOn)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(pkg.NewGRPCRateLimiter(lb)))

	greeter.RegisterGreeterServer(server, &greeter.UnimplementedGreeterServer{})
	fmt.Println("listening on ", listenOn)
	if err = server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
