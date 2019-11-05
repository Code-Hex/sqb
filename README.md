# SQL Query Builder

> Blazing fast, Flexible, SQL Query Builder for Go

[![GoDoc](https://godoc.org/github.com/Code-Hex/sqb?status.svg)](https://godoc.org/github.com/Code-Hex/sqb) [![CircleCI](https://circleci.com/gh/Code-Hex/sqb.svg?style=svg&circle-token=0ff0570576e90eb3a10e017f7ca1279748565daf)](https://circleci.com/gh/Code-Hex/sqb) [![codecov](https://codecov.io/gh/Code-Hex/sqb/branch/master/graph/badge.svg?token=xjioT8q5f5)](https://codecov.io/gh/Code-Hex/sqb) [![Go Report Card](https://goreportcard.com/badge/github.com/Code-Hex/sqb)](https://goreportcard.com/report/github.com/Code-Hex/sqb)

## Synopsis

When used normally

```go
const sqlstr = "SELECT * FROM tables WHERE ?"
builder := sqb.New(sqlstr).Bind(sqb.Eq("category", 1))
query, args, err := builder.Build()
// query => "SELECT * FROM tables WHERE category = ?",
// args  => []interface{}{1}
```

<details>
<summary>When you want to use build cache</summary>


```go
const sqlstr = "SELECT * FROM tables WHERE ? AND ?"
cached := sqb.New(sqlstr).Bind(sqb.Eq("category", 1))

for _, col := range columns {
    builder := cached.Bind(sqb.Eq(col, "value"))
    query, args, err  := builder.Build()
    // query => "SELECT * FROM tables WHERE category = ? AND " + col + " = ?",
    // args  => []interface{}{1, "value"}
}
```
</details>

<details>
<summary>Error case</summary>


```go
const sqlstr = "SELECT * FROM tables WHERE ? OR ?"
builder := sqb.New(sqlstr).Bind(sqb.Eq("category", 1))
query, args, err  := builder.Build()
// query => "",
// args  => nil
// err   => "number of bindVars exceeds replaceable statements"
```
</details>

## Install

Use `go get` to install this package.

    go get -u github.com/Code-Hex/sqb