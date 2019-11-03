package pool

// Builder is the interface that wraps the basic
// Reset, Cap and WriteString method.
type Builder struct {
	buf  Buffer
	args []interface{}
}

// String returns appended the contents.
func (b *Builder) String() string {
	return b.buf.String()
}

// Args return appended args.
func (b *Builder) Args() []interface{} {
	return b.args
}

// WriteString appends the contents of s to Buffer.
func (b *Builder) WriteString(s string) {
	b.buf.WriteString(s)
}

// AppendArgs appends the args.
func (b *Builder) AppendArgs(args ...interface{}) {
	b.args = append(b.args, args...)
}

// Reset resets Builder.
func (b *Builder) Reset() {
	b.buf.Reset()
	b.args = b.args[:0]
}

// Get allocates a new strings.Builder or grabs a cached one.
func Get() *Builder {
	return globalPool.Get().(*Builder)
}

// Put saves used strings.Builder; avoids an allocation per invocation.
func Put(b *Builder) {
	// Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized buffer, we add a hard limit on the maximum buffer
	// to place back in the pool.
	//
	// See https://golang.org/issue/23199
	if b.buf.Cap() > 64<<10 {
		return
	}
	b.Reset()
	globalPool.Put(b)
}
