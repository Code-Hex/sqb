package pool

// Builder is the interface that wraps the basic
// Reset, Cap and WriteString method.
type Builder interface {
	Reset()
	Cap() int
	WriteString(string) (int, error)
}

// Get allocates a new strings.Builder or grabs a cached one.
func Get() Builder {
	return globalPool.Get().(Builder)
}

// Put saves used strings.Builder; avoids an allocation per invocation.
func Put(builder Builder) {
	// Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized buffer, we add a hard limit on the maximum buffer
	// to place back in the pool.
	//
	// See https://golang.org/issue/23199
	if builder.Cap() > 64<<10 {
		return
	}
	builder.Reset()
	globalPool.Put(builder)
}
