package limiter

import "time"

type Limiter interface {
	Limit(key string) (delay time.Duration, rejected bool)
	GetLimit() int
}
