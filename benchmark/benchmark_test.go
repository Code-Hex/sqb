package benchmark

import (
	"testing"

	"github.com/Code-Hex/sqb"
	sq "github.com/Masterminds/squirrel"
	"github.com/huandu/go-sqlbuilder"
)

func BenchmarkSqbAndFromMap(b *testing.B) {
	const want = "SELECT * FROM users WHERE (col1 = ? AND col2 = ? AND col3 = ?)"
	for i := 0; i < b.N; i++ {
		builder := sqb.New("SELECT * FROM users WHERE ?")
		sql, _, err := builder.Bind(
			sqb.Paren(
				sqb.AndFromMap(sqb.Eq, map[string]interface{}{
					"col1": "world",
					"col2": 100,
					"col3": true,
				}),
			),
		).Build()
		if err != nil {
			b.Fatal(err)
		}
		if want != sql {
			b.Fatalf("\nwant %q\ngot %q", want, sql)
		}
	}
}

func BenchmarkSqbAnd(b *testing.B) {
	const want = "SELECT * FROM users WHERE (col1 = ? AND col2 = ? AND col3 = ?)"
	for i := 0; i < b.N; i++ {
		builder := sqb.New("SELECT * FROM users WHERE ?")
		sql, _, err := builder.Bind(
			sqb.Paren(
				sqb.And(
					sqb.Eq("col1", "world"),
					sqb.Eq("col2", 100),
					sqb.Eq("col3", true),
				),
			),
		).Build()
		if err != nil {
			b.Fatal(err)
		}
		if want != sql {
			b.Fatalf("\nwant %q\ngot %q", want, sql)
		}
	}
}

func BenchmarkSquirrel(b *testing.B) {
	const want = "SELECT * FROM users WHERE (col1 = ? AND col2 = ? AND col3 = ?)"
	for i := 0; i < b.N; i++ {
		users := sq.Select("*").From("users")
		active := users.Where(
			sq.And{
				sq.Eq{"col1": "world"},
				sq.Eq{"col2": 100},
				sq.Eq{"col3": true},
			},
		)
		sql, _, err := active.ToSql()
		if err != nil {
			b.Fatal(err)
		}
		if want != sql {
			b.Fatalf("\nwant %q\ngot %q", want, sql)
		}
	}
}

func BenchmarkSqlbuilder(b *testing.B) {
	const want = "SELECT * FROM users WHERE (col1 = ? AND col2 = ? AND col3 = ?)"
	for i := 0; i < b.N; i++ {
		sb := sqlbuilder.NewSelectBuilder()
		sb.Select("*")
		sb.From("users")
		sb.Where(
			sb.And(
				sb.Equal("col1", "world"),
				sb.Equal("col2", 100),
				sb.Equal("col3", true),
			),
		)
		sql, _ := sb.Build()
		if want != sql {
			b.Fatalf("\nwant %q\ngot %q", want, sql)
		}
	}
}
