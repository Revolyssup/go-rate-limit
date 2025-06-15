package limiter

import (
	"time"

	"github.com/Revolyssup/go-rate-limit/pkg/store"
)

type Limiter interface {
	Limit(key string) (delay time.Duration, rejected bool)
	GetLimit() int
	InitStore(store.Store)
}
