# DEPIQ
`depiq` is an expressive SQL builder and executor
    
If you are upgrading from an older version please read the [Migrating Between Versions](./docs/version_migration.md) docs.


## Installation

If using go modules.

```sh
go get -u github.com/orn-id/depiq/v9
```

If you are not using go modules...

**NOTE** You should still be able to use this package if you are using go version `>v1.10` but, you will need to drop the version from the package. `import "github.com/orn-id/depiq/v9` -> `import "github.com/orn-id/depiq"`

```sh
go get -u github.com/orn-id/depiq
```

### [Migrating Between Versions](./docs/version_migration.md)

## Features

`depiq` comes with many features but here are a few of the more notable ones

* Query Builder
* Parameter interpolation (e.g `SELECT * FROM "items" WHERE "id" = ?` -> `SELECT * FROM "items" WHERE "id" = 1`)
* Built from the ground up with multiple dialects in mind
* Insert, Multi Insert, Update, and Delete support
* Scanning of rows to struct[s] or primitive value[s]

While depiq may support the scanning of rows into structs it is not intended to be used as an ORM if you are looking for common ORM features like associations,
or hooks I would recommend looking at some of the great ORM libraries such as:

* [gorm](https://github.com/jinzhu/gorm)
* [hood](https://github.com/eaigner/hood)

## Why?

We tried a few other sql builders but each was a thin wrapper around sql fragments that we found error prone. `depiq` was built with the following goals in mind:

* Make the generation of SQL easy and enjoyable
* Create an expressive DSL that would find common errors with SQL at compile time.
* Provide a DSL that accounts for the common SQL expressions, NOT every nuance for each database.
* Provide developers the ability to:
  * Use SQL when desired
  * Easily scan results into primitive values and structs
  * Use the native sql.Db methods when desired

## Docs

* [Dialect](./docs/dialect.md) - Introduction to different dialects (`mysql`, `postgres`, `sqlite3`, `sqlserver` etc) 
* [Expressions](./docs/expressions.md) - Introduction to `depiq` expressions and common examples.
* [Select Dataset](./docs/selecting.md) - Docs and examples about creating and executing SELECT sql statements.
* [Insert Dataset](./docs/inserting.md) - Docs and examples about creating and executing INSERT sql statements.
* [Update Dataset](./docs/updating.md) - Docs and examples about creating and executing UPDATE sql statements.
* [Delete Dataset](./docs/deleting.md) - Docs and examples about creating and executing DELETE sql statements.
* [Prepared Statements](./docs/interpolation.md) - Docs about interpolation and prepared statements in `depiq`.
* [Database](./docs/database.md) - Docs and examples of using a Database to execute queries in `depiq`
* [Working with time.Time](./docs/time.md) - Docs on how to use alternate time locations.

## Quick Examples

### Select

See the [select dataset](./docs/selecting.md) docs for more in depth examples

```go
sql, _, _ := depiq.From("test").ToSQL()
fmt.Println(sql)
```

Output:

```
SELECT * FROM "test"
```

```go
sql, _, _ := depiq.From("test").Where(depiq.Ex{
	"d": []string{"a", "b", "c"},
}).ToSQL()
fmt.Println(sql)
```

Output:

```
SELECT * FROM "test" WHERE ("d" IN ('a', 'b', 'c'))
```

### Insert

See the [insert dataset](./docs/inserting.md) docs for more in depth examples

```go
ds := depiq.Insert("user").
	Cols("first_name", "last_name").
	Vals(
		depiq.Vals{"Greg", "Farley"},
		depiq.Vals{"Jimmy", "Stewart"},
		depiq.Vals{"Jeff", "Jeffers"},
	)
insertSQL, args, _ := ds.ToSQL()
fmt.Println(insertSQL, args)
```

Output: 
```sql
INSERT INTO "user" ("first_name", "last_name") VALUES ('Greg', 'Farley'), ('Jimmy', 'Stewart'), ('Jeff', 'Jeffers') []
```

```go
ds := depiq.Insert("user").Rows(
	depiq.Record{"first_name": "Greg", "last_name": "Farley"},
	depiq.Record{"first_name": "Jimmy", "last_name": "Stewart"},
	depiq.Record{"first_name": "Jeff", "last_name": "Jeffers"},
)
insertSQL, args, _ := ds.ToSQL()
fmt.Println(insertSQL, args)
```

Output:
```
INSERT INTO "user" ("first_name", "last_name") VALUES ('Greg', 'Farley'), ('Jimmy', 'Stewart'), ('Jeff', 'Jeffers') []
```


```go
type User struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}
ds := depiq.Insert("user").Rows(
	User{FirstName: "Greg", LastName: "Farley"},
	User{FirstName: "Jimmy", LastName: "Stewart"},
	User{FirstName: "Jeff", LastName: "Jeffers"},
)
insertSQL, args, _ := ds.ToSQL()
fmt.Println(insertSQL, args)
```

Output:
```
INSERT INTO "user" ("first_name", "last_name") VALUES ('Greg', 'Farley'), ('Jimmy', 'Stewart'), ('Jeff', 'Jeffers') []
```

```go
ds := depiq.Insert("user").Prepared(true).
	FromQuery(depiq.From("other_table"))
insertSQL, args, _ := ds.ToSQL()
fmt.Println(insertSQL, args)
```

Output:
```
INSERT INTO "user" SELECT * FROM "other_table" []
```

```go
ds := depiq.Insert("user").Prepared(true).
	Cols("first_name", "last_name").
	FromQuery(depiq.From("other_table").Select("fn", "ln"))
insertSQL, args, _ := ds.ToSQL()
fmt.Println(insertSQL, args)
```

Output:
```
INSERT INTO "user" ("first_name", "last_name") SELECT "fn", "ln" FROM "other_table" []
```

### Update

See the [update dataset](./docs/updating.md) docs for more in depth examples

```go
sql, args, _ := depiq.Update("items").Set(
	depiq.Record{"name": "Test", "address": "111 Test Addr"},
).ToSQL()
fmt.Println(sql, args)
```

Output:
```
UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
```

```go
type item struct {
	Address string `db:"address"`
	Name    string `db:"name" depiq:"skipupdate"`
}
sql, args, _ := depiq.Update("items").Set(
	item{Name: "Test", Address: "111 Test Addr"},
).ToSQL()
fmt.Println(sql, args)
```

Output:
```
UPDATE "items" SET "address"='111 Test Addr' []
```

```go
sql, _, _ := depiq.Update("test").
	Set(depiq.Record{"foo": "bar"}).
	Where(depiq.Ex{
		"a": depiq.Op{"gt": 10}
	}).ToSQL()
fmt.Println(sql)
```

Output:
```
UPDATE "test" SET "foo"='bar' WHERE ("a" > 10)
```

### Delete

See the [delete dataset](./docs/deleting.md) docs for more in depth examples

```go
ds := depiq.Delete("items")

sql, args, _ := ds.ToSQL()
fmt.Println(sql, args)
```

```go
sql, _, _ := depiq.Delete("test").Where(depiq.Ex{
		"c": nil
	}).ToSQL()
fmt.Println(sql)
```

Output:
```
DELETE FROM "test" WHERE ("c" IS NULL)
```