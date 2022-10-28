# Selecting

* [Creating a SelectDataset](#create)
* Building SQL
  * [`Select`](#select)
  * [`Distinct`](#distinct)
  * [`From`](#from)
  * [`Join`](#joins)
  * [`Where`](#where)
  * [`Limit`](#limit)
  * [`Offset`](#offset)
  * [`GroupBy`](#group_by)
  * [`Having`](#having)
  * [`Window`](#window)
  * [`With`](#with)
  * [`SetError`](#seterror)
  * [`ForUpdate`](#forupdate)
* Executing Queries
  * [`ScanStructs`](#scan-structs) - Scans rows into a slice of structs
  * [`ScanStruct`](#scan-struct) - Scans a row into a slice a struct, returns false if a row wasnt found
  * [`ScanVals`](#scan-vals)- Scans a rows of 1 column into a slice of primitive values
  * [`ScanVal`](#scan-val) - Scans a row of 1 column into a primitive value, returns false if a row wasnt found.
  * [`Scanner`](#scanner) - Allows you to interatively scan rows into structs or values.
  * [`Count`](#count) - Returns the count for the current query
  * [`Pluck`](#pluck) - Selects a single column and stores the results into a slice of primitive values

<a name="create"></a>
To create a [`SelectDataset`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset)  you can use

**[`depiq.From`](https://godoc.org/github.com/orn-id/depiq/#From) and [`depiq.Select`](https://godoc.org/github.com/orn-id/depiq/#Select)**

When you just want to create some quick SQL, this mostly follows the `Postgres` with the exception of placeholders for prepared statements.

```go
sql, _, _ := depiq.From("table").ToSQL()
fmt.Println(sql)

sql, _, _ := depiq.Select(depiq.L("NOW()")).ToSQL()
fmt.Println(sql)
```
Output:
```
SELECT * FROM "table"
SELECT NOW()
```

**[`DialectWrapper.From`](https://godoc.org/github.com/orn-id/depiq/#DialectWrapper.From) and [`DialectWrapper.Select`](https://godoc.org/github.com/orn-id/depiq/#DialectWrapper.Select)**

Use this when you want to create SQL for a specific `dialect`

```go
// import _ "github.com/orn-id/depiq/v9/dialect/mysql"

dialect := depiq.Dialect("mysql")

sql, _, _ := dialect.From("table").ToSQL()
fmt.Println(sql)

sql, _, _ := dialect.Select(depiq.L("NOW()")).ToSQL()
fmt.Println(sql)
```
Output:
```
SELECT * FROM `table`
SELECT NOW()
```

**[`Database.From`](https://godoc.org/github.com/orn-id/depiq/#DialectWrapper.From) and [`Database.Select`](https://godoc.org/github.com/orn-id/depiq/#DialectWrapper.From)**

Use this when you want to execute the SQL or create SQL for the drivers dialect.

```go
// import _ "github.com/orn-id/depiq/v9/dialect/mysql"

mysqlDB := //initialize your db
db := depiq.New("mysql", mysqlDB)

sql, _, _ := db.From("table").ToSQL()
fmt.Println(sql)

sql, _, _ := db.Select(depiq.L("NOW()")).ToSQL()
fmt.Println(sql)
```
Output:
```
SELECT * FROM `table`
SELECT NOW()
```

### Examples

For more examples visit the **[Docs](https://godoc.org/github.com/orn-id/depiq/#SelectDataset)**

<a name="select"></a>
**[`Select`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.Select)**

```go
sql, _, _ := depiq.From("test").Select("a", "b", "c").ToSQL()
fmt.Println(sql)
```

Output:
```sql
SELECT "a", "b", "c" FROM "test"
```

You can also ues another dataset in your select

```go
ds := depiq.From("test")
fromDs := ds.Select("age").Where(depiq.C("age").Gt(10))
sql, _, _ := ds.From().Select(fromDs).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT (SELECT "age" FROM "test" WHERE ("age" > 10))
```

Selecting a literal

```go
sql, _, _ := depiq.From("test").Select(depiq.L("a + b").As("sum")).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT a + b AS "sum" FROM "test"
```

Select aggregate functions

```go
sql, _, _ := depiq.From("test").Select(
	depiq.COUNT("*").As("age_count"),
	depiq.MAX("age").As("max_age"),
	depiq.AVG("age").As("avg_age"),
).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT COUNT(*) AS "age_count", MAX("age") AS "max_age", AVG("age") AS "avg_age" FROM "test"
```

Selecting columns from a struct

```go
ds := depiq.From("test")

type myStruct struct {
	Name         string
	Address      string `db:"address"`
	EmailAddress string `db:"email_address"`
}

// Pass with pointer
sql, _, _ := ds.Select(&myStruct{}).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT "address", "email_address", "name" FROM "test"
```

<a name="distinct"></a>
**[`Distinct`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.Distinct)**

```go
sql, _, _ := depiq.From("test").Select("a", "b").Distinct().ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT DISTINCT "a", "b" FROM "test"
```

If you dialect supports `DISTINCT ON` you provide arguments to the `Distinct` method.

**NOTE** currently only the `postgres` and the default dialects support `DISTINCT ON` clauses

```go
sql, _, _ := depiq.From("test").Distinct("a").ToSQL()
fmt.Println(sql)
```
Output:

```
SELECT DISTINCT ON ("a") * FROM "test"
```

You can also provide other expression arguments

With `depiq.L`

```go
sql, _, _ := depiq.From("test").Distinct(depiq.L("COALESCE(?, ?)", depiq.C("a"), "empty")).ToSQL()
fmt.Println(sql)
```
Output:
```
SELECT DISTINCT ON (COALESCE("a", 'empty')) * FROM "test"
```
With `depiq.Coalesce`
```go
sql, _, _ := depiq.From("test").Distinct(depiq.COALESCE(depiq.C("a"), "empty")).ToSQL()
fmt.Println(sql)
```
Output:
```
SELECT DISTINCT ON (COALESCE("a", 'empty')) * FROM "test"
```

<a name="from"></a>
**[`From`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.From)**

Overriding the original from
```go
ds := depiq.From("test")
sql, _, _ := ds.From("test2").ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT * FROM "test2"
```

From another dataset

```go
ds := depiq.From("test")
fromDs := ds.Where(depiq.C("age").Gt(10))
sql, _, _ := ds.From(fromDs).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT * FROM (SELECT * FROM "test" WHERE ("age" > 10)) AS "t1"
```

From an aliased dataset

```go
ds := depiq.From("test")
fromDs := ds.Where(depiq.C("age").Gt(10))
sql, _, _ := ds.From(fromDs.As("test2")).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT * FROM (SELECT * FROM "test" WHERE ("age" > 10)) AS "test2"
```

Lateral Query

```go
maxEntry := depiq.From("entry").
	Select(depiq.MAX("int").As("max_int")).
	Where(depiq.Ex{"time": depiq.Op{"lt": depiq.I("e.time")}}).
	As("max_entry")

maxId := depiq.From("entry").
	Select("id").
	Where(depiq.Ex{"int": depiq.I("max_entry.max_int")}).
	As("max_id")

ds := depiq.
	Select("e.id", "max_entry.max_int", "max_id.id").
	From(
		depiq.T("entry").As("e"),
		depiq.Lateral(maxEntry),
		depiq.Lateral(maxId),
	)
query, args, _ := ds.ToSQL()
fmt.Println(query, args)

query, args, _ = ds.Prepared(true).ToSQL()
fmt.Println(query, args)
```

Output
```
SELECT "e"."id", "max_entry"."max_int", "max_id"."id" FROM "entry" AS "e", LATERAL (SELECT MAX("int") AS "max_int" FROM "entry" WHERE ("time" < "e"."time")) AS "max_entry", LATERAL (SELECT "id" FROM "entry" WHERE ("int" = "max_entry"."max_int")) AS "max_id" []
SELECT "e"."id", "max_entry"."max_int", "max_id"."id" FROM "entry" AS "e", LATERAL (SELECT MAX("int") AS "max_int" FROM "entry" WHERE ("time" < "e"."time")) AS "max_entry", LATERAL (SELECT "id" FROM "entry" WHERE ("int" = "max_entry"."max_int")) AS "max_id" []
```

<a name="joins"></a>
**[`Join`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.Join)**

```go
sql, _, _ := depiq.From("test").Join(
	depiq.T("test2"),
	depiq.On(depiq.Ex{"test.fkey": depiq.I("test2.Id")}),
).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT * FROM "test" INNER JOIN "test2" ON ("test"."fkey" = "test2"."Id")
```

[`InnerJoin`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.InnerJoin)

```go
sql, _, _ := depiq.From("test").InnerJoin(
	depiq.T("test2"),
	depiq.On(depiq.Ex{"test.fkey": depiq.I("test2.Id")}),
).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT * FROM "test" INNER JOIN "test2" ON ("test"."fkey" = "test2"."Id")
```

[`FullOuterJoin`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.FullOuterJoin)

```go
sql, _, _ := depiq.From("test").FullOuterJoin(
	depiq.T("test2"),
	depiq.On(depiq.Ex{
		"test.fkey": depiq.I("test2.Id"),
	}),
).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT * FROM "test" FULL OUTER JOIN "test2" ON ("test"."fkey" = "test2"."Id")
```

[`RightOuterJoin`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.RightOuterJoin)

```go
sql, _, _ := depiq.From("test").RightOuterJoin(
	depiq.T("test2"),
	depiq.On(depiq.Ex{
		"test.fkey": depiq.I("test2.Id"),
	}),
).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT * FROM "test" RIGHT OUTER JOIN "test2" ON ("test"."fkey" = "test2"."Id")
```

[`LeftOuterJoin`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.LeftOuterJoin)

```go
sql, _, _ := depiq.From("test").LeftOuterJoin(
	depiq.T("test2"),
	depiq.On(depiq.Ex{
		"test.fkey": depiq.I("test2.Id"),
	}),
).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT * FROM "test" LEFT OUTER JOIN "test2" ON ("test"."fkey" = "test2"."Id")
```

[`FullJoin`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.FullJoin)

```go
sql, _, _ := depiq.From("test").FullJoin(
	depiq.T("test2"),
	depiq.On(depiq.Ex{
		"test.fkey": depiq.I("test2.Id"),
	}),
).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT * FROM "test" FULL JOIN "test2" ON ("test"."fkey" = "test2"."Id")
```


[`RightJoin`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.RightJoin)

```go
sql, _, _ := depiq.From("test").RightJoin(
	depiq.T("test2"),
	depiq.On(depiq.Ex{
		"test.fkey": depiq.I("test2.Id"),
	}),
).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT * FROM "test" RIGHT JOIN "test2" ON ("test"."fkey" = "test2"."Id")
```

[`LeftJoin`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.LeftJoin)

```go
sql, _, _ := depiq.From("test").LeftJoin(
	depiq.T("test2"),
	depiq.On(depiq.Ex{
		"test.fkey": depiq.I("test2.Id"),
	}),
).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT * FROM "test" LEFT JOIN "test2" ON ("test"."fkey" = "test2"."Id")
```

[`NaturalJoin`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.NaturalJoin)

```go
sql, _, _ := depiq.From("test").NaturalJoin(depiq.T("test2")).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT * FROM "test" NATURAL JOIN "test2"
```

[`NaturalLeftJoin`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.NaturalLeftJoin)

```go
sql, _, _ := depiq.From("test").NaturalLeftJoin(depiq.T("test2")).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT * FROM "test" NATURAL LEFT JOIN "test2"
```

[`NaturalRightJoin`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.NaturalRightJoin)

```go
sql, _, _ := depiq.From("test").NaturalRightJoin(depiq.T("test2")).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT * FROM "test" NATURAL RIGHT LEFT JOIN "test2"
```

[`NaturalFullJoin`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.NaturalFullJoin)

```go
sql, _, _ := depiq.From("test").NaturalFullJoin(depiq.T("test2")).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT * FROM "test" NATURAL FULL LEFT JOIN "test2"
```

[`CrossJoin`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.CrossJoin)

```go
sql, _, _ := depiq.From("test").CrossJoin(depiq.T("test2")).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT * FROM "test" CROSS JOIN "test2"
```

Join with a Lateral

```go
maxEntry := depiq.From("entry").
	Select(depiq.MAX("int").As("max_int")).
	Where(depiq.Ex{"time": depiq.Op{"lt": depiq.I("e.time")}}).
	As("max_entry")

maxId := depiq.From("entry").
	Select("id").
	Where(depiq.Ex{"int": depiq.I("max_entry.max_int")}).
	As("max_id")

ds := depiq.
	Select("e.id", "max_entry.max_int", "max_id.id").
	From(depiq.T("entry").As("e")).
	Join(depiq.Lateral(maxEntry), depiq.On(depiq.V(true))).
	Join(depiq.Lateral(maxId), depiq.On(depiq.V(true)))
query, args, _ := ds.ToSQL()
fmt.Println(query, args)

query, args, _ = ds.Prepared(true).ToSQL()
fmt.Println(query, args)
```

Output:
```
SELECT "e"."id", "max_entry"."max_int", "max_id"."id" FROM "entry" AS "e" INNER JOIN LATERAL (SELECT MAX("int") AS "max_int" FROM "entry" WHERE ("time" < "e"."time")) AS "max_entry" ON TRUE INNER JOIN LATERAL (SELECT "id" FROM "entry" WHERE ("int" = "max_entry"."max_int")) AS "max_id" ON TRUE []

SELECT "e"."id", "max_entry"."max_int", "max_id"."id" FROM "entry" AS "e" INNER JOIN LATERAL (SELECT MAX("int") AS "max_int" FROM "entry" WHERE ("time" < "e"."time")) AS "max_entry" ON ? INNER JOIN LATERAL (SELECT "id" FROM "entry" WHERE ("int" = "max_entry"."max_int")) AS "max_id" ON ? [true true]
```

<a name="where"></a>
**[`Where`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.Where)**

You can use `depiq.Ex` to create an ANDed condition
```go
sql, _, _ := depiq.From("test").Where(depiq.Ex{
	"a": depiq.Op{"gt": 10},
	"b": depiq.Op{"lt": 10},
	"c": nil,
	"d": []string{"a", "b", "c"},
}).ToSQL()
fmt.Println(sql)
```

Output:

```
SELECT * FROM "test" WHERE (("a" > 10) AND ("b" < 10) AND ("c" IS NULL) AND ("d" IN ('a', 'b', 'c')))
```

You can use `depiq.ExOr` to create an ORed condition

```go
sql, _, _ := depiq.From("test").Where(depiq.ExOr{
	"a": depiq.Op{"gt": 10},
	"b": depiq.Op{"lt": 10},
	"c": nil,
	"d": []string{"a", "b", "c"},
}).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT * FROM "test" WHERE (("a" > 10) OR ("b" < 10) OR ("c" IS NULL) OR ("d" IN ('a', 'b', 'c')))
```

You can use `depiq.Ex` with `depiq.ExOr` for complex expressions

```go
// You can use Or with Ex to Or multiple Ex maps together
sql, _, _ := depiq.From("test").Where(
	depiq.Or(
		depiq.Ex{
			"a": depiq.Op{"gt": 10},
			"b": depiq.Op{"lt": 10},
		},
		depiq.Ex{
			"c": nil,
			"d": []string{"a", "b", "c"},
		},
	),
).ToSQL()
fmt.Println(sql)
```

Output:

```
SELECT * FROM "test" WHERE ((("a" > 10) AND ("b" < 10)) OR (("c" IS NULL) AND ("d" IN ('a', 'b', 'c'))))
```

You can also use identifiers to create your where condition

```go
sql, _, _ := depiq.From("test").Where(
	depiq.C("a").Gt(10),
	depiq.C("b").Lt(10),
	depiq.C("c").IsNull(),
	depiq.C("d").In("a", "b", "c"),
).ToSQL()
fmt.Println(sql)
```

Output:
```
SELECT * FROM "test" WHERE (("a" > 10) AND ("b" < 10) AND ("c" IS NULL) AND ("d" IN ('a', 'b', 'c')))
```

Using `depiq.Or` to create ORed expression

```go
// You can use a combination of Ors and Ands
sql, _, _ := depiq.From("test").Where(
	depiq.Or(
		depiq.C("a").Gt(10),
		depiq.And(
			depiq.C("b").Lt(10),
			depiq.C("c").IsNull(),
		),
	),
).ToSQL()
fmt.Println(sql)
```

Output:

```
SELECT * FROM "test" WHERE (("a" > 10) OR (("b" < 10) AND ("c" IS NULL)))
```

<a name="limit"></a>
**[`Limit`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.Limit)**

```go
ds := depiq.From("test").Limit(10)
sql, _, _ := ds.ToSQL()
fmt.Println(sql)
```

Output:

```
SELECT * FROM "test" LIMIT 10
```

<a name="offset"></a>
**[`Offset`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.Offset)**

```go
ds := depiq.From("test").Offset(2)
sql, _, _ := ds.ToSQL()
fmt.Println(sql)
```

Output:

```
SELECT * FROM "test" OFFSET 2
```

<a name="group_by"></a>
**[`GroupBy`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.GroupBy)**

```go
sql, _, _ := depiq.From("test").
	Select(depiq.SUM("income").As("income_sum")).
	GroupBy("age").
	ToSQL()
fmt.Println(sql)
```

Output:

```
SELECT SUM("income") AS "income_sum" FROM "test" GROUP BY "age"
```

<a name="having"></a>
**[`Having`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.Having)**

```go
sql, _, _ = depiq.From("test").GroupBy("age").Having(depiq.SUM("income").Gt(1000)).ToSQL()
fmt.Println(sql)
```

Output:

```
SELECT * FROM "test" GROUP BY "age" HAVING (SUM("income") > 1000)
```

<a name="with"></a>
**[`With`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.With)**

To use CTEs in `SELECT` statements you can use the `With` method.

Simple Example

```go
sql, _, _ := depiq.From("one").
	With("one", depiq.From().Select(depiq.L("1"))).
	Select(depiq.Star()).
	ToSQL()
fmt.Println(sql)
```

Output:

```
WITH one AS (SELECT 1) SELECT * FROM "one"
```

Dependent `WITH` clauses:

```go
sql, _, _ = depiq.From("derived").
	With("intermed", depiq.From("test").Select(depiq.Star()).Where(depiq.C("x").Gte(5))).
	With("derived", depiq.From("intermed").Select(depiq.Star()).Where(depiq.C("x").Lt(10))).
	Select(depiq.Star()).
	ToSQL()
fmt.Println(sql)
```

Output:
```
WITH intermed AS (SELECT * FROM "test" WHERE ("x" >= 5)), derived AS (SELECT * FROM "intermed" WHERE ("x" < 10)) SELECT * FROM "derived"
```

`WITH` clause with arguments

```go
sql, _, _ = depiq.From("multi").
		With("multi(x,y)", depiq.From().Select(depiq.L("1"), depiq.L("2"))).
		Select(depiq.C("x"), depiq.C("y")).
		ToSQL()
fmt.Println(sql)
```

Output:
```
WITH multi(x,y) AS (SELECT 1, 2) SELECT "x", "y" FROM "multi"
```

Using a `InsertDataset`.

```go
insertDs := depiq.Insert("foo").Rows(depiq.Record{"user_id": 10}).Returning("id")

ds := depiq.From("bar").
	With("ins", insertDs).
	Select("bar_name").
	Where(depiq.Ex{"bar.user_id": depiq.I("ins.user_id")})

sql, _, _ := ds.ToSQL()
fmt.Println(sql)

sql, args, _ := ds.Prepared(true).ToSQL()
fmt.Println(sql, args)
```
Output:
```
WITH ins AS (INSERT INTO "foo" ("user_id") VALUES (10) RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "ins"."user_id")
WITH ins AS (INSERT INTO "foo" ("user_id") VALUES (?) RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "ins"."user_id") [10]
```

Using an `UpdateDataset`

```go
updateDs := depiq.Update("foo").Set(depiq.Record{"bar": "baz"}).Returning("id")

ds := depiq.From("bar").
	With("upd", updateDs).
	Select("bar_name").
	Where(depiq.Ex{"bar.user_id": depiq.I("ins.user_id")})

sql, _, _ := ds.ToSQL()
fmt.Println(sql)

sql, args, _ := ds.Prepared(true).ToSQL()
fmt.Println(sql, args)
```

Output:
```
WITH upd AS (UPDATE "foo" SET "bar"='baz' RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "upd"."user_id")
WITH upd AS (UPDATE "foo" SET "bar"=? RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "upd"."user_id") [baz]
```

Using a `DeleteDataset`

```go
deleteDs := depiq.Delete("foo").Where(depiq.Ex{"bar": "baz"}).Returning("id")

ds := depiq.From("bar").
	With("del", deleteDs).
	Select("bar_name").
	Where(depiq.Ex{"bar.user_id": depiq.I("del.user_id")})

sql, _, _ := ds.ToSQL()
fmt.Println(sql)

sql, args, _ := ds.Prepared(true).ToSQL()
fmt.Println(sql, args)
```

Output:
```
WITH del AS (DELETE FROM "foo" WHERE ("bar" = 'baz') RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "del"."user_id")
WITH del AS (DELETE FROM "foo" WHERE ("bar" = ?) RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "del"."user_id") [baz]
```

<a name="window"></a>
**[`Window Function`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.Window)**

**NOTE** currently only the `postgres`, `mysql8` (NOT `mysql`) and the default dialect support `Window Function`

To use windowing in `SELECT` statements you can use the `Over` method on an `SQLFunction`

```go
sql, _, _ := depiq.From("test").Select(
	depiq.ROW_NUMBER().Over(depiq.W().PartitionBy("a").OrderBy(depiq.I("b").Asc())),
)
fmt.Println(sql)
```

Output:

```
SELECT ROW_NUMBER() OVER (PARTITION BY "a" ORDER BY "b") FROM "test"
```

`depiq` also supports the `WINDOW` clause.

```go
sql, _, _ := depiq.From("test").
	Select(depiq.ROW_NUMBER().OverName(depiq.I("w"))).
	Window(depiq.W("w").PartitionBy("a").OrderBy(depiq.I("b").Asc()))
fmt.Println(sql)
```

Output:

```
SELECT ROW_NUMBER() OVER "w" FROM "test" WINDOW "w" AS (PARTITION BY "a" ORDER BY "b")
```

<a name="seterror"></a>
**[`SetError`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.SetError)**

Sometimes while building up a query with depiq you will encounter situations where certain
preconditions are not met or some end-user contraint has been violated. While you could
track this error case separately, depiq provides a convenient built-in mechanism to set an
error on a dataset if one has not already been set to simplify query building.

Set an Error on a dataset:

```go
func GetSelect(name string) *depiq.SelectDataset {

    var ds = depiq.From("test")

    if len(name) == 0 {
        return ds.SetError(fmt.Errorf("name is empty"))
    }

    return ds.Select(name)
}

```

This error is returned on any subsequent call to `Error` or `ToSQL`:

```go
var name string = ""
ds = GetSelect(name)
fmt.Println(ds.Error())

sql, args, err = ds.ToSQL()
fmt.Println(err)
```

Output:
```
name is empty
name is empty
```

<a name="forupdate"></a>
**[`ForUpdate`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.ForUpdate)**

```go
sql, _, _ := depiq.From("test").ForUpdate(exp.Wait).ToSQL()
fmt.Println(sql)
```

Output:
```sql
SELECT * FROM "test" FOR UPDATE
```

If your dialect supports FOR UPDATE OF you provide tables to be locked as variable arguments to the ForUpdate method.

```go
sql, _, _ := depiq.From("test").ForUpdate(exp.Wait, depiq.T("test")).ToSQL()
fmt.Println(sql)
```

Output:
```sql
SELECT * FROM "test" FOR UPDATE OF "test"
```

## Executing Queries

To execute your query use [`depiq.Database#From`](https://godoc.org/github.com/orn-id/depiq/#Database.From) to create your dataset

<a name="scan-structs"></a>
**[`ScanStructs`](http://godoc.org/github.com/orn-id/depiq#SelectDataset.ScanStructs)**

Scans rows into a slice of structs

**NOTE** [`ScanStructs`](http://godoc.org/github.com/orn-id/depiq#SelectDataset.ScanStructs) will only select the columns that can be scanned in to the structs unless you have explicitly selected certain columns.

 ```go
type User struct{
  FirstName string `db:"first_name"`
  LastName  string `db:"last_name"`
  Age       int    `db:"-"` // a field that shouldn't be selected
}

var users []User
//SELECT "first_name", "last_name" FROM "user";
if err := db.From("user").ScanStructs(&users); err != nil{
  panic(err.Error())
}
fmt.Printf("\n%+v", users)

var users []User
//SELECT "first_name" FROM "user";
if err := db.From("user").Select("first_name").ScanStructs(&users); err != nil{
  panic(err.Error())
}
fmt.Printf("\n%+v", users)
```

`depiq` also supports scanning into multiple structs. In the example below we define a `Role` and `User` struct that could both be used individually to scan into. However, you can also create a new struct that adds both structs as fields that can be populated in a single query.

**NOTE** When calling `ScanStructs` without a select already defined it will automatically only `SELECT` the columns found in the struct, omitting any that are tagged with `db:"-"`

 ```go
type Role struct {
  Id     uint64 `db:"id"`
	UserID uint64 `db:"user_id"`
	Name   string `db:"name"`
}
type User struct {
	Id        uint64 `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}
type UserAndRole struct {
	User User `db:"depiq_user"` // tag as the "depiq_user" table
	Role Role `db:"user_role"` // tag as "user_role" table
}
db := getDb()

ds := db.
	From("depiq_user").
	Join(depiq.T("user_role"), depiq.On(depiq.I("depiq_user.id").Eq(depiq.I("user_role.user_id"))))
var users []UserAndRole
	// Scan structs will auto build the
if err := ds.ScanStructs(&users); err != nil {
	fmt.Println(err.Error())
	return
}
for _, u := range users {
	fmt.Printf("\n%+v", u)
}
```

You can alternatively manually select the columns with the appropriate aliases using the `depiq.C` method to create the alias.

```go
type Role struct {
	UserID uint64 `db:"user_id"`
	Name   string `db:"name"`
}
type User struct {
	Id        uint64 `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Role      Role   `db:"user_role"` // tag as "user_role" table
}
db := getDb()

ds := db.
	Select(
		"depiq_user.id",
		"depiq_user.first_name",
		"depiq_user.last_name",
		// alias the fully qualified identifier `C` is important here so it doesnt parse it
		depiq.I("user_role.user_id").As(depiq.C("user_role.user_id")),
		depiq.I("user_role.name").As(depiq.C("user_role.name")),
	).
	From("depiq_user").
	Join(depiq.T("user_role"), depiq.On(depiq.I("depiq_user.id").Eq(depiq.I("user_role.user_id"))))

var users []User
if err := ds.ScanStructs(&users); err != nil {
	fmt.Println(err.Error())
	return
}
for _, u := range users {
	fmt.Printf("\n%+v", u)
}
```

<a name="scan-struct"></a>
**[`ScanStruct`](http://godoc.org/github.com/orn-id/depiq#SelectDataset.ScanStruct)**

Scans a row into a slice a struct, returns false if a row wasnt found

**NOTE** [`ScanStruct`](http://godoc.org/github.com/orn-id/depiq#SelectDataset.ScanStruct) will only select the columns that can be scanned in to the struct unless you have explicitly selected certain columns.

```go
type User struct{
  FirstName string `db:"first_name"`
  LastName  string `db:"last_name"`
  Age       int    `db:"-"` // a field that shouldn't be selected
}

var user User
// SELECT "first_name", "last_name" FROM "user" LIMIT 1;
found, err := db.From("user").ScanStruct(&user)
if err != nil{
  fmt.Println(err.Error())
  return
}
if !found {
  fmt.Println("No user found")
} else {
  fmt.Printf("\nFound user: %+v", user)
}
```

`depiq` also supports scanning into multiple structs. In the example below we define a `Role` and `User` struct that could both be used individually to scan into. However, you can also create a new struct that adds both structs as fields that can be populated in a single query.

**NOTE** When calling `ScanStruct` without a select already defined it will automatically only `SELECT` the columns found in the struct, omitting any that are tagged with `db:"-"`

 ```go
type Role struct {
	UserID uint64 `db:"user_id"`
	Name   string `db:"name"`
}
type User struct {
	ID        uint64 `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}
type UserAndRole struct {
	User User `db:"depiq_user"` // tag as the "depiq_user" table
	Role Role `db:"user_role"` // tag as "user_role" table
}
db := getDb()
var userAndRole UserAndRole
ds := db.
	From("depiq_user").
	Join(depiq.T("user_role"),depiq.On(depiq.I("depiq_user.id").Eq(depiq.I("user_role.user_id")))).
	Where(depiq.C("first_name").Eq("Bob"))

found, err := ds.ScanStruct(&userAndRole)
if err != nil{
  fmt.Println(err.Error())
  return
}
if !found {
  fmt.Println("No user found")
} else {
  fmt.Printf("\nFound user: %+v", user)
}
```

You can alternatively manually select the columns with the appropriate aliases using the `depiq.C` method to create the alias.

```go
type Role struct {
	UserID uint64 `db:"user_id"`
	Name   string `db:"name"`
}
type User struct {
	ID        uint64 `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Role      Role   `db:"user_role"` // tag as "user_role" table
}
db := getDb()
var userAndRole UserAndRole
ds := db.
	Select(
		"depiq_user.id",
		"depiq_user.first_name",
		"depiq_user.last_name",
		// alias the fully qualified identifier `C` is important here so it doesnt parse it
		depiq.I("user_role.user_id").As(depiq.C("user_role.user_id")),
		depiq.I("user_role.name").As(depiq.C("user_role.name")),
	).
	From("depiq_user").
	Join(depiq.T("user_role"),depiq.On(depiq.I("depiq_user.id").Eq(depiq.I("user_role.user_id")))).
	Where(depiq.C("first_name").Eq("Bob"))

found, err := ds.ScanStruct(&userAndRole)
if err != nil{
  fmt.Println(err.Error())
  return
}
if !found {
  fmt.Println("No user found")
} else {
  fmt.Printf("\nFound user: %+v", user)
}
```


**NOTE** Using the `depiq.SetColumnRenameFunction` function, you can change the function that's used to rename struct fields when struct tags aren't defined

```go
import "strings"

depiq.SetColumnRenameFunction(strings.ToUpper)

type User struct{
  FirstName string
  LastName string
}

var user User
//SELECT "FIRSTNAME", "LASTNAME" FROM "user" LIMIT 1;
found, err := db.From("user").ScanStruct(&user)
// ...
```

**NOTE** Using the `depiq.SetIgnoreUntaggedFields(true)` function, you can cause depiq to ignore any fields that aren't explicitly tagged.

```go
depiq.SetIgnoreUntaggedFields(true)

type User struct{
  FirstName string `db:"first_name"`
  LastName string
}

var user User
//SELECT "first_name" FROM "user" LIMIT 1;
found, err := db.From("user").ScanStruct(&user)
// ...
```


<a name="scan-vals"></a>
**[`ScanVals`](http://godoc.org/github.com/orn-id/depiq#SelectDataset.ScanVals)**

Scans a rows of 1 column into a slice of primitive values

```go
var ids []int64
if err := db.From("user").Select("id").ScanVals(&ids); err != nil{
  fmt.Println(err.Error())
  return
}
fmt.Printf("\n%+v", ids)
```

<a name="scan-val"></a>
[`ScanVal`](http://godoc.org/github.com/orn-id/depiq#SelectDataset.ScanVal)

Scans a row of 1 column into a primitive value, returns false if a row wasnt found.

**Note** when using the dataset a `LIMIT` of 1 is automatically applied.
```go
var id int64
found, err := db.From("user").Select("id").ScanVal(&id)
if err != nil{
  fmt.Println(err.Error())
  return
}
if !found{
  fmt.Println("No id found")
}else{
  fmt.Printf("\nFound id: %d", id)
}
```

<a name="scanner"></a>
**[`Scanner`](http://godoc.org/github.com/orn-id/depiq/exec#Scanner)**

Scanner knows how to scan rows into structs. This is useful when dealing with large result sets where you can have only one item scanned in memory at one time.

In the following example we scan each row into struct.

```go

type User struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}
db := getDb()

scanner, err := db.
  From("depiq_user").
	Select("first_name", "last_name").
	Where(depiq.Ex{
		"last_name": "Yukon",
	}).
	Executor().
	Scanner()

if err != nil {
	fmt.Println(err.Error())
	return
}

defer scanner.Close()

for scanner.Next() {
	u := User{}

	err = scanner.ScanStruct(&u)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("\n%+v", u)
}

if scanner.Err() != nil {
	fmt.Println(scanner.Err().Error())
}
```

In this example we scan each row into a val.
```go
db := getDb()

scanner, err := db.
	From("depiq_user").
	Select("first_name").
	Where(depiq.Ex{
		"last_name": "Yukon",
	}).
	Executor().
	Scanner()

if err != nil {
	fmt.Println(err.Error())
	return
}

defer scanner.Close()

for scanner.Next() {
	name := ""

	err = scanner.ScanVal(&name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(name)
}

if scanner.Err() != nil {
	fmt.Println(scanner.Err().Error())
}
```

<a name="count"></a>
**[`Count`](http://godoc.org/github.com/orn-id/depiq#SelectDataset.Count)**

Returns the count for the current query

```go
count, err := db.From("user").Count()
if err != nil{
  fmt.Println(err.Error())
  return
}
fmt.Printf("\nCount:= %d", count)
```

<a name="pluck"></a>
**[`Pluck`](http://godoc.org/github.com/orn-id/depiq#SelectDataset.Pluck)**

Selects a single column and stores the results into a slice of primitive values

```go
var ids []int64
if err := db.From("user").Pluck(&ids, "id"); err != nil{
  fmt.Println(err.Error())
  return
}
fmt.Printf("\nIds := %+v", ids)
```
