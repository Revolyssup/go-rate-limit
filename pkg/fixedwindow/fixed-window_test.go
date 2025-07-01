package fixedwindow

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixedWindow(t *testing.T) {
	testcases := []struct {
		fw        *FixedWindow
		operation func(*FixedWindow, string)
		key       string
	}{{
		//Limit 2 requests in 5 second window
		fw: NewFixedWindow(2, 2, 0),
		operation: func(fw *FixedWindow, key string) {
			_, rejected := fw.Limit(key)
			assert.Equal(t, rejected, false)
			_, rejected = fw.Limit(key)
			assert.Equal(t, rejected, false)
			_, rejected = fw.Limit(key)
			assert.Equal(t, rejected, true)

			//Reset
			time.Sleep(3 * time.Second)
			_, rejected = fw.Limit(key)
			assert.Equal(t, rejected, false)
			_, rejected = fw.Limit(key)
			assert.Equal(t, rejected, false)
			_, rejected = fw.Limit(key)
			assert.Equal(t, rejected, true)
		},
	}}

	for _, tt := range testcases {
		tt.operation(tt.fw, tt.key)
	}
}
