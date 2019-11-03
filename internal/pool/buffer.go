package pool

// Buffer is the interface that wraps the basic
// Reset, Cap and WriteString method.
type Buffer interface {
	Reset()
	Cap() int
	WriteString(string) (int, error)
	String() string
}
