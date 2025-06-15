package pkg

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Revolyssup/go-rate-limit/pkg/limiter"
	"github.com/Revolyssup/go-rate-limit/pkg/store"
	redis "github.com/redis/go-redis/v9"
)

type HTTPRateLimiter struct {
	limiter     limiter.Limiter
	key         Key
	redisClient *redis.ClusterClient
}

type Key string

const (
	IPKey     Key = "IPKEY"
	HEADERKey Key = "HEADERKEY"
)

func (k Key) GetValue(r *http.Request) string {
	switch k {
	case IPKey:
		return string(k) + r.RemoteAddr
	case HEADERKey:
		return string(k) + strings.Join(r.Header.Values(string(k)), ",")
	}
	return "" // will never happen
}

func (k Key) IsSupported() bool {
	switch k {
	case IPKey:
		return true
	case HEADERKey:
		return true
	}
	return false
}

type Options interface {
	Set(rl *HTTPRateLimiter)
}

type RedisOptions struct {
	Addrs []string
}

func (r *RedisOptions) Set(rl *HTTPRateLimiter) {
	rl.redisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: r.Addrs,
	})
	rl.limiter.InitStore(store.NewRedisStore(rl.redisClient))
}

func NewHTTPRateLimiter(lim limiter.Limiter, key Key, opts ...Options) (*HTTPRateLimiter, error) {
	if !key.IsSupported() {
		return nil, fmt.Errorf("cannot rate limit on key type: %s", string(key))
	}

	rl := &HTTPRateLimiter{
		limiter: lim,
		key:     key,
	}

	for _, op := range opts {
		op.Set(rl)
	}
	return rl, nil
}

const (
	REMAINING  = "X-Ratelimit-Remaining"
	LIMIT      = "X-Ratelimit-Limit"
	RetryAfter = "X-Ratelimit-Retry-After"
)

func (rl *HTTPRateLimiter) RateLimit(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			delay    time.Duration
			rejected bool
		)
		val := rl.key.GetValue(r)
		delay, rejected = rl.limiter.Limit(val)
		if rejected {
			w.Header().Set(LIMIT, fmt.Sprintf("%d", rl.limiter.GetLimit()))
			w.WriteHeader(429)
			w.Write([]byte("rate limit exceeded"))
			return
		} else {
			if delay != 0 {
				time.Sleep(delay)
			}
			h.ServeHTTP(w, r)
		}
	})
}
