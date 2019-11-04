package stmt

var _ error = (*BuildError)(nil)

// BuildError is the error type usually returned by functions in the stmt
// package. It describes the current operation, occurred of
// an error.
type BuildError struct {
	Op  string
	Err error
}

func (b *BuildError) Error() string {
	if b == nil {
		return "<nil>"
	}
	return b.Op + ": " + b.Err.Error()
}

// Unwrap unwraps the wrapped error.
// This method is implemented to satisfy an interface of errors.Unwap.
func (b *BuildError) Unwrap() error {
	return b.Err
}
