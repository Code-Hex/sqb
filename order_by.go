package sqb

import "github.com/Code-Hex/sqb/stmt"

// OrderBy Creates an unary expression for ORDER BY.
// If you want to know more details, See at stmt.OrderBy.
func OrderBy(column string, desc bool) *stmt.OrderBy {
	return &stmt.OrderBy{
		Column: column,
		Desc:   desc,
	}
}

// OrderByList Creates an expression for ORDER BY from multiple *stmt.OrderBy.
//
// This function creates like "<column_name>, <column_name> DESC". The first argument
// is required. If you want to know more details, See at stmt.OrderBy.
func OrderByList(expr *stmt.OrderBy, exprs ...*stmt.OrderBy) *stmt.OrderBy {
	ret := expr
	for _, o := range exprs {
		tmp := &stmt.OrderBy{
			Column: o.Column,
			Desc:   o.Desc,
		}
		ret.Next = tmp
		ret = tmp
	}
	return expr
}
