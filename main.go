package main

import (
	"net/http"

	"github.com/Revolyssup/go-rate-limit/pkg"
	leakybucket "github.com/Revolyssup/go-rate-limit/pkg/leaky-bucket"
)

func main() {
	lb := leakybucket.NewLeakyBucket(1, 2)
	h := http.NewServeMux()
	h.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	}))
	// opts := &pkg.RedisOptions{
	// 	Addrs: []string{
	// 		"localhost:5000",
	// 		"localhost:5002",
	// 		"localhost:5003",
	// 		"localhost:5004",
	// 		"localhost:5005",
	// 		"localhost:5006",
	// 	},
	// }
	rl, err := pkg.NewRateLimiter(lb, pkg.HEADERKey)
	if err != nil {
		panic(err)
	}
	http.ListenAndServe(":8080", rl.RateLimit(h))
}
