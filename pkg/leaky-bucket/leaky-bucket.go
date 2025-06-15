package leakybucket

import (
	"sync"
	"time"
)

type state struct {
	lastReq int64
	excess  int64
}
type LeakyBucket struct {
	rate    int64
	burst   int64
	buckets map[string]*state
	mx      sync.Mutex
}

const factor = 1000_000

func NewLeakyBucket(rate int, burst int) *LeakyBucket {
	return &LeakyBucket{
		rate:    int64(rate) * factor,
		burst:   int64(rate) * factor,
		buckets: make(map[string]*state),
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
func (lb *LeakyBucket) Limit(key string) (delay time.Duration, rejected bool) {
	lb.mx.Lock()
	defer lb.mx.Unlock()
	now := time.Now().UnixMicro()
	var elapsed int64
	if lb.buckets[key] == nil {
		lb.buckets[key] = &state{
			lastReq: 0,
			excess:  0,
		}
	}
	// Calculate time since last request
	if lb.buckets[key].lastReq != 0 {
		elapsed = now - lb.buckets[key].lastReq
	}
	// check leakage (How much has leaked)
	leakage := (int64(lb.rate) * elapsed) / factor
	prevExcess := lb.buckets[key].excess
	excess := max((lb.buckets[key].excess-leakage), 0) + factor
	if excess > int64(lb.burst) {
		return 0, true
	}
	lb.buckets[key].lastReq = now
	lb.buckets[key].excess = excess

	delaySeconds := float64(prevExcess) / float64(lb.rate)
	return time.Duration(delaySeconds * float64(time.Second)), false
}
