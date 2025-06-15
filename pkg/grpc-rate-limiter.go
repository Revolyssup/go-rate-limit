package pkg

import (
	"context"
	"fmt"

	"github.com/Revolyssup/go-rate-limit/pkg/limiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func NewGRPCRateLimiter(lim limiter.Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		p, ok := peer.FromContext(ctx)
		if !ok {
			return nil, fmt.Errorf("no peer info available")
		}
		ip := p.Addr.String()
		_, rejected := lim.Limit(ip)
		if rejected {
			return nil, fmt.Errorf("rate limited")
		}
		return handler(ctx, req)
	}
}
