// +build go1.10

package pool

import (
	"strings"
	"sync"
)

var _ Buffer = (*strings.Builder)(nil)

var globalPool = sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
}
