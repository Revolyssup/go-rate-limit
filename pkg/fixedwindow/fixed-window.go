package fixedwindow

import (
	"sync"
	"time"

	"github.com/Revolyssup/go-rate-limit/pkg/store"
)

/*
Within a window of t seconds - n requests are allowed
interface -> InitStore(store), Limit(key string) (delay, rejected), GetLimit()
-> Calulate currtime-lasttime > t -> allowed, set passedRequests = 1
-> else ->if passedRequests > allowed -> reject else passedRequests++
*/
type FixedWindow struct {
	windowsize      int
	allowedRequests int
	s               store.Store
	mx              sync.Mutex
	burst           int
}

func NewFixedWindow(windowSizeInSeconds, allowedRequests, burst int) *FixedWindow {
	return &FixedWindow{
		windowsize:      windowSizeInSeconds,
		allowedRequests: allowedRequests,
		burst:           burst,
	}
}

func (f *FixedWindow) GetLimit() int {
	return f.allowedRequests / f.windowsize
}

func (f *FixedWindow) InitStore(s store.Store) {
	f.s = s
}

func (f *FixedWindow) Limit(key string) (delay time.Duration, rejected bool) {
	f.mx.Lock()
	defer f.mx.Unlock()
	currtime := time.Now().Unix()
	windowStart := (currtime / int64(f.windowsize)) * int64(f.windowsize)

	if f.s == nil {
		f.s = store.NewLocalStore()
	}
	state, _ := f.s.Get(key)
	stateWindowStart := (state.LastReq / int64(f.windowsize)) * int64(f.windowsize)
	// Reset counter if window advanced
	if stateWindowStart != windowStart {
		state.Excess = 0
		state.LastReq = currtime
	}

	if state.Excess >= int64(f.allowedRequests) {
		return 0, true // Reject
	}

	state.Excess++
	f.s.Set(key, state)
	return 0, false // Allow
}
