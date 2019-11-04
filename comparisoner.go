package sqb

import (
	"github.com/Code-Hex/sqb/stmt"
)

// ConditionalFunc indicates function of conditional.
type ConditionalFunc func(column string, value interface{}) *stmt.Condition

// Eq creates condition `column = ?`.
func Eq(column string, value interface{}) *stmt.Condition {
	return Op("=", column, value)
}

// Ge creates condition `column >= ?`.
func Ge(column string, value interface{}) *stmt.Condition {
	return Op(">=", column, value)
}

// Gt creates condition `column > ?`.
func Gt(column string, value interface{}) *stmt.Condition {
	return Op(">", column, value)
}

// Le creates condition `column <= ?`.
func Le(column string, value interface{}) *stmt.Condition {
	return Op("<=", column, value)
}

// Lt creates condition `column < ?`.
func Lt(column string, value interface{}) *stmt.Condition {
	return Op("<", column, value)
}

// Ne creates condition `column != ?`.
func Ne(column string, value interface{}) *stmt.Condition {
	return Op("!=", column, value)
}

// Op creates flexible compare operation.
func Op(op, column string, value interface{}) *stmt.Condition {
	return &stmt.Condition{
		Column: column,
		Compare: &stmt.CompOp{
			Op:    op,
			Value: value,
		},
	}
}

// Like creates condition `column LIKE ?`.
func Like(column string, value interface{}) *stmt.Condition {
	return &stmt.Condition{
		Column: column,
		Compare: &stmt.CompLike{
			Negative: false,
			Value:    value,
		},
	}
}

// NotLike creates condition `column NOT LIKE ?`.
func NotLike(column string, value interface{}) *stmt.Condition {
	return &stmt.Condition{
		Column: column,
		Compare: &stmt.CompLike{
			Negative: true,
			Value:    value,
		},
	}
}

// Between creates condition `column BETWEEN ? AND ?`.
func Between(column string, left, right interface{}) *stmt.Condition {
	return &stmt.Condition{
		Column: column,
		Compare: &stmt.CompBetween{
			Negative: false,
			Left:     left,
			Right:    right,
		},
	}
}

// NotBetween creates condition `column NOT BETWEEN ? AND ?`.
func NotBetween(column string, left, right interface{}) *stmt.Condition {
	return &stmt.Condition{
		Column: column,
		Compare: &stmt.CompBetween{
			Negative: true,
			Left:     left,
			Right:    right,
		},
	}
}

// In creates condition `column IN (?, ?, ?, ...)`.
func In(column string, args ...interface{}) *stmt.Condition {
	return &stmt.Condition{
		Column: column,
		Compare: &stmt.CompIn{
			Negative: false,
			Values:   args,
		},
	}
}

// NotIn creates condition `column NOT IN (?, ?, ?, ...)`.
func NotIn(column string, args ...interface{}) *stmt.Condition {
	return &stmt.Condition{
		Column: column,
		Compare: &stmt.CompIn{
			Negative: true,
			Values:   args,
		},
	}
}
