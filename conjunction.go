package sqb

import (
	"sort"

	"github.com/Code-Hex/sqb/stmt"
)

// Paren creates the expression with parentheses.
func Paren(expr stmt.Expr) *stmt.Paren {
	return &stmt.Paren{
		Expr: expr,
	}
}

// And creates statement for the AND boolean expression.
// If you want to know more details, See at stmt.And.
func And(left, right stmt.Expr, exprs ...stmt.Expr) *stmt.And {
	ret := &stmt.And{
		Left:  left,
		Right: right,
	}
	for _, expr := range exprs {
		ret = &stmt.And{
			Left:  ret,
			Right: expr,
		}
	}
	return ret
}

// Or creates statement for the OR boolean expression with parentheses.
// If you want to know more details, See at stmt.Or.
func Or(left, right stmt.Expr, exprs ...stmt.Expr) *stmt.Or {
	ret := &stmt.Or{
		Left:  left,
		Right: right,
	}
	for _, expr := range exprs {
		ret = &stmt.Or{
			Left:  ret,
			Right: expr,
		}
	}
	return ret
}

// AndFromMap Creates a concatenated string of AND boolean expression from a map.
//
// If there is no first argument then occurs panic.
// If map length is zero it returns nil.
// If map length is 1 it returns *stmt.Condition created by ConditionalFunc.
func AndFromMap(f ConditionalFunc, m map[string]interface{}) stmt.Expr {
	if f == nil {
		panic("unspecified function")
	}
	if len(m) < 2 {
		// Length is zero
		for key, val := range m {
			return f(key, val)
		}
		// return nil if length is zero
		return nil
	}
	exprs := convertMapToStmts(f, m)
	return And(exprs[0], exprs[1], exprs[2:]...)
}

// OrFromMap Creates a concatenated string of OR boolean expression from a map.
//
// If there is no first argument then occurs panic.
// If map length is zero it returns nil.
// If map length is 1 it returns *stmt.Condition created by ConditionalFunc.
func OrFromMap(f ConditionalFunc, m map[string]interface{}) stmt.Expr {
	if f == nil {
		panic("unspecified function")
	}
	if len(m) < 2 {
		// Length is zero
		for key, val := range m {
			return f(key, val)
		}
		// return nil if length is zero
		return nil
	}
	exprs := convertMapToStmts(f, m)
	return Or(exprs[0], exprs[1], exprs[2:]...)
}

func convertMapToStmts(f ConditionalFunc, m map[string]interface{}) []stmt.Expr {
	i, keys := 0, make([]string, len(m))
	for key := range m {
		keys[i] = key
		i++
	}
	// This is to guarantee the order
	// when concatenating strings
	sort.Strings(keys)

	exprs := make([]stmt.Expr, len(m))
	for idx, key := range keys {
		exprs[idx] = f(key, m[key])
	}
	return exprs
}
