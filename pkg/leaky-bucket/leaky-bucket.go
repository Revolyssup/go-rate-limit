package leakybucket

import (
	"sync"
	"time"

	"github.com/Revolyssup/go-rate-limit/pkg/store"
)

type LeakyBucket struct {
	rate    int64
	burst   int64
	buckets store.Store
	mx      sync.Mutex
}

const factor = 1000_000

func NewLeakyBucket(rate int, burst int) *LeakyBucket {
	return &LeakyBucket{
		rate:  int64(rate) * factor,
		burst: int64(rate) * factor,
	}
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func (lb *LeakyBucket) GetLimit() int {
	return int(lb.rate) / factor
}

func (lb *LeakyBucket) InitStore(s store.Store) {
	lb.buckets = s
}
func (lb *LeakyBucket) Limit(key string) (delay time.Duration, rejected bool) {
	lb.mx.Lock()
	defer lb.mx.Unlock()
	now := time.Now().UnixMicro()
	var elapsed int64
	if lb.buckets == nil {
		lb.buckets = store.NewLocalStore()
	}
	s, _ := lb.buckets.Get(key)
	// Calculate time since last request
	if s.LastReq != 0 {
		elapsed = now - s.LastReq
	}
	// check leakage (How much has leaked)
	leakage := (int64(lb.rate) * elapsed) / factor
	prevExcess := s.Excess
	excess := max((prevExcess-leakage), 0) + factor
	if excess > int64(lb.burst) {
		return 0, true
	}
	lb.buckets.Set(key, store.State{
		LastReq: now,
		Excess:  excess,
	})
	delaySeconds := float64(prevExcess) / float64(lb.rate)
	return time.Duration(delaySeconds * float64(time.Second)), false
}
