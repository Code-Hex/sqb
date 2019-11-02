// +build !go1.10

package pool

import (
	"bytes"
	"sync"
)

var _ Builder = (*bytes.Buffer)(nil)

var globalPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}
