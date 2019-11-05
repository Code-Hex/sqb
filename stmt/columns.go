package stmt

import "errors"

// Columns represents columns field.
type Columns []string

// Write writes a string with concatenated column names.
// i.e. "col1",  "col1, col2, col3"
func (c Columns) Write(b Builder) error {
	switch len(c) {
	case 0:
		return errors.New("unspecified columns")
	case 1:
		b.WriteString(c[0])
		return nil
	}
	b.WriteString(c[0])
	for _, column := range c[1:] {
		b.WriteString(", ")
		b.WriteString(column)
	}
	return nil
}
