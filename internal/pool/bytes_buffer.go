// +build !go1.10

package pool

import (
	"bytes"
	"sync"
)

var _ Buffer = (*bytes.Buffer)(nil)

var globalPool = sync.Pool{
	New: func() interface{} {
		return &Builder{
			buf:  new(bytes.Buffer),
			args: make([]interface{}, 0, 3),
		}
	},
}
