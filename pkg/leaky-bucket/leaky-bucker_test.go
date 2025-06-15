package leakybucket

import (
	"fmt"
	"testing"
	"time"
)

func TestLeakyBucket(t *testing.T) {
	testCases := []struct {
		lb        *LeakyBucket
		operation func(lb *LeakyBucket, key string) error
		key       string
	}{{
		//Limit 1 request per second
		lb: NewLeakyBucket(1, 1),
		operation: func(lb *LeakyBucket, key string) error {
			delay, rejected := lb.Limit(key)
			if rejected {
				return fmt.Errorf("should not reject first request")
			}
			if delay != 0 {
				return fmt.Errorf("there should be no delay")
			}
			_, rejected = lb.Limit(key)
			if !rejected {
				return fmt.Errorf("should reject second request")
			}
			time.Sleep(1 * time.Second)
			_, rejected = lb.Limit(key)
			if rejected {
				return fmt.Errorf("should not reject third request")
			}
			return nil
		},
	}}

	for _, tc := range testCases {
		if err := tc.operation(tc.lb, tc.key); err != nil {
			t.Error(err.Error())
		}
	}
}
