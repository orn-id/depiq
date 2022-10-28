# Updating

* [Create a UpdateDataset](#create)
* Examples
  * [Set with `depiq.Record`](#set-record)
  * [Set with struct](#set-struct)
  * [Set with map](#set-map)
  * [Multi Table](#from)
  * [Where](#where)
  * [Order](#order)
  * [Limit](#limit)
  * [Returning](#returning)
  * [SetError](#seterror)
  * [Executing](#executing)

<a name="create"></a>
To create a [`UpdateDataset`](https://godoc.org/github.com/orn-id/depiq/#UpdateDataset)  you can use

**[`depiq.Update`](https://godoc.org/github.com/orn-id/depiq/#Update)**

When you just want to create some quick SQL, this mostly follows the `Postgres` with the exception of placeholders for prepared statements.

```go
ds := depiq.Update("user").Set(
    depiq.Record{"first_name": "Greg", "last_name": "Farley"},
)
updateSQL, _, _ := ds.ToSQL()
fmt.Println(insertSQL, args)
```
Output:
```
UPDATE "user" SET "first_name"='Greg', "last_name"='Farley'
```

**[`SelectDataset.Update`](https://godoc.org/github.com/orn-id/depiq/#SelectDataset.Update)**

If you already have a `SelectDataset` you can invoke `Update()` to get a `UpdateDataset`

**NOTE** This method will also copy over the `WITH`, `WHERE`, `ORDER`, and `LIMIT` clauses from the update

```go
ds := depiq.From("user")

updateSQL, _, _ := ds.Update().Set(
    depiq.Record{"first_name": "Greg", "last_name": "Farley"},
).ToSQL()
fmt.Println(insertSQL, args)

updateSQL, _, _ = ds.Where(depiq.C("first_name").Eq("Gregory")).Update().Set(
    depiq.Record{"first_name": "Greg", "last_name": "Farley"},
).ToSQL()
fmt.Println(insertSQL, args)
```
Output:
```
UPDATE "user" SET "first_name"='Greg', "last_name"='Farley'
UPDATE "user" SET "first_name"='Greg', "last_name"='Farley' WHERE "first_name"='Gregory'
```

**[`DialectWrapper.Update`](https://godoc.org/github.com/orn-id/depiq/#DialectWrapper.Update)**

Use this when you want to create SQL for a specific `dialect`

```go
// import _ "github.com/orn-id/depiq/v9/dialect/mysql"

dialect := depiq.Dialect("mysql")

ds := dialect.Update("user").Set(
    depiq.Record{"first_name": "Greg", "last_name": "Farley"},
)
updateSQL, _, _ := ds.ToSQL()
fmt.Println(insertSQL, args)
```
Output:
```
UPDATE `user` SET `first_name`='Greg', `last_name`='Farley'
```

**[`Database.Update`](https://godoc.org/github.com/orn-id/depiq/#DialectWrapper.Update)**

Use this when you want to execute the SQL or create SQL for the drivers dialect.

```go
// import _ "github.com/orn-id/depiq/v9/dialect/mysql"

mysqlDB := //initialize your db
db := depiq.New("mysql", mysqlDB)

ds := db.Update("user").Set(
    depiq.Record{"first_name": "Greg", "last_name": "Farley"},
)
updateSQL, _, _ := ds.ToSQL()
fmt.Println(insertSQL, args)
```
Output:
```
UPDATE `user` SET `first_name`='Greg', `last_name`='Farley'
```

### Examples

For more examples visit the **[Docs](https://godoc.org/github.com/orn-id/depiq/#UpdateDataset)**

<a name="set-record"></a>
**[Set with `depiq.Record`](https://godoc.org/github.com/orn-id/depiq/#UpdateDataset.Set)**

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

<a name="set-struct"></a>
**[Set with Struct](https://godoc.org/github.com/orn-id/depiq/#UpdateDataset.Set)**

```go
type item struct {
	Address string `db:"address"`
	Name    string `db:"name"`
}
sql, args, _ := depiq.Update("items").Set(
	item{Name: "Test", Address: "111 Test Addr"},
).ToSQL()
fmt.Println(sql, args)
```

Output:
```
UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
```

With structs you can also skip fields by using the `skipupdate` tag

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

If you want to use the database `DEFAULT` when the struct field is a zero value you can use the `defaultifempty` tag.

```go
type item struct {
	Address string `db:"address"`
	Name    string `db:"name" depiq:"defaultifempty"`
}
sql, args, _ := depiq.Update("items").Set(
	item{Address: "111 Test Addr"},
).ToSQL()
fmt.Println(sql, args)
```

Output:
```
UPDATE "items" SET "address"='111 Test Addr',"name"=DEFAULT []
```

`depiq` will also use fields in embedded structs when creating an update.

**NOTE** unexported fields will be ignored!

```go
type Address struct {
	Street string `db:"address_street"`
	State  string `db:"address_state"`
}
type User struct {
	Address
	FirstName string
	LastName  string
}
ds := depiq.Update("user").Set(
	User{Address: Address{Street: "111 Street", State: "NY"}, FirstName: "Greg", LastName: "Farley"},
)
updateSQL, args, _ := ds.ToSQL()
fmt.Println(updateSQL, args)
```

Output:
```
UPDATE "user" SET "address_state"='NY',"address_street"='111 Street',"firstname"='Greg',"lastname"='Farley' []
```

**NOTE** When working with embedded pointers if the embedded struct is nil then the fields will be ignored.

```go
type Address struct {
	Street string
	State  string
}
type User struct {
	*Address
	FirstName string
	LastName  string
}
ds := depiq.Update("user").Set(
	User{FirstName: "Greg", LastName: "Farley"},
)
updateSQL, args, _ := ds.ToSQL()
fmt.Println(updateSQL, args)
```

Output:
```
UPDATE "user" SET "firstname"='Greg',"lastname"='Farley' []
```

You can ignore an embedded struct or struct pointer by using `db:"-"`

```go
type Address struct {
	Street string
	State  string
}
type User struct {
	Address   `db:"-"`
	FirstName string
	LastName  string
}
ds := depiq.Update("user").Set(
	User{Address: Address{Street: "111 Street", State: "NY"}, FirstName: "Greg", LastName: "Farley"},
)
updateSQL, args, _ := ds.ToSQL()
fmt.Println(updateSQL, args)
```

Output:
```
UPDATE "user" SET "firstname"='Greg',"lastname"='Farley' []
```


<a name="set-map"></a>
**[Set with Map](https://godoc.org/github.com/orn-id/depiq/#UpdateDataset.Set)**

```go
sql, args, _ := depiq.Update("items").Set(
	map[string]interface{}{"name": "Test", "address": "111 Test Addr"},
).ToSQL()
fmt.Println(sql, args)
```

Output:
```
UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
```

<a name="from"></a>
**[From / Multi Table](https://godoc.org/github.com/orn-id/depiq/#UpdateDataset.From)**

`depiq` allows joining multiple tables in a update clause through `From`.

**NOTE** The `sqlite3` adapter does not support a multi table syntax.

`Postgres` Example

```go
dialect := depiq.Dialect("postgres")

ds := dialect.Update("table_one").
    Set(depiq.Record{"foo": depiq.I("table_two.bar")}).
    From("table_two").
    Where(depiq.Ex{"table_one.id": depiq.I("table_two.id")})

sql, _, _ := ds.ToSQL()
fmt.Println(sql)
```

Output:
```sql
UPDATE "table_one" SET "foo"="table_two"."bar" FROM "table_two" WHERE ("table_one"."id" = "table_two"."id")
```

`MySQL` Example

```go
dialect := depiq.Dialect("mysql")

ds := dialect.Update("table_one").
    Set(depiq.Record{"foo": depiq.I("table_two.bar")}).
    From("table_two").
    Where(depiq.Ex{"table_one.id": depiq.I("table_two.id")})

sql, _, _ := ds.ToSQL()
fmt.Println(sql)
```
Output:
```sql
UPDATE `table_one`,`table_two` SET `foo`=`table_two`.`bar` WHERE (`table_one`.`id` = `table_two`.`id`)
```

<a name="where"></a>
**[Where](https://godoc.org/github.com/orn-id/depiq/#UpdateDataset.Where)**

```go
sql, _, _ := depiq.Update("test").
	Set(depiq.Record{"foo": "bar"}).
	Where(depiq.Ex{
		"a": depiq.Op{"gt": 10},
		"b": depiq.Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).ToSQL()
fmt.Println(sql)
```

Output:
```
UPDATE "test" SET "foo"='bar' WHERE (("a" > 10) AND ("b" < 10) AND ("c" IS NULL) AND ("d" IN ('a', 'b', 'c')))
```

<a name="order"></a>
**[Order](https://godoc.org/github.com/orn-id/depiq/#UpdateDataset.Order)**

**NOTE** This will only work if your dialect supports it

```go
// import _ "github.com/orn-id/depiq/v9/dialect/mysql"

ds := depiq.Dialect("mysql").
	Update("test").
	Set(depiq.Record{"foo": "bar"}).
	Order(depiq.C("a").Asc())
sql, _, _ := ds.ToSQL()
fmt.Println(sql)
```

Output:
```
UPDATE `test` SET `foo`='bar' ORDER BY `a` ASC
```

<a name="limit"></a>
**[Order](https://godoc.org/github.com/orn-id/depiq/#UpdateDataset.Limit)**

**NOTE** This will only work if your dialect supports it

```go
// import _ "github.com/orn-id/depiq/v9/dialect/mysql"

ds := depiq.Dialect("mysql").
	Update("test").
	Set(depiq.Record{"foo": "bar"}).
	Limit(10)
sql, _, _ := ds.ToSQL()
fmt.Println(sql)
```

Output:
```
UPDATE `test` SET `foo`='bar' LIMIT 10
```

<a name="returning"></a>
**[Returning](https://godoc.org/github.com/orn-id/depiq/#UpdateDataset.Returning)**

Returning a single column example.

```go
sql, _, _ := depiq.Update("test").
	Set(depiq.Record{"foo": "bar"}).
	Returning("id").
	ToSQL()
fmt.Println(sql)
```

Output:
```
UPDATE "test" SET "foo"='bar' RETURNING "id"
```

Returning multiple columns

```go
sql, _, _ := depiq.Update("test").
	Set(depiq.Record{"foo": "bar"}).
	Returning("a", "b").
	ToSQL()
fmt.Println(sql)
```

Output:
```
UPDATE "test" SET "foo"='bar' RETURNING "a", "b"
```

Returning all columns

```go
sql, _, _ := depiq.Update("test").
	Set(depiq.Record{"foo": "bar"}).
	Returning(depiq.T("test").All()).
	ToSQL()
fmt.Println(sql)
```

Output:
```
UPDATE "test" SET "foo"='bar' RETURNING "test".*
```

<a name="seterror"></a>
**[`SetError`](https://godoc.org/github.com/orn-id/depiq/#UpdateDataset.SetError)**

Sometimes while building up a query with depiq you will encounter situations where certain
preconditions are not met or some end-user contraint has been violated. While you could
track this error case separately, depiq provides a convenient built-in mechanism to set an
error on a dataset if one has not already been set to simplify query building.

Set an Error on a dataset:

```go
func GetUpdate(name string, value string) *depiq.UpdateDataset {

    var ds = depiq.Update("test")

    if len(name) == 0 {
        return ds.SetError(fmt.Errorf("name is empty"))
    }

    if len(value) == 0 {
        return ds.SetError(fmt.Errorf("value is empty"))
    }

    return ds.Set(depiq.Record{name: value})
}

```

This error is returned on any subsequent call to `Error` or `ToSQL`:

```go
var field, value string
ds = GetUpdate(field, value)
fmt.Println(ds.Error())

sql, args, err = ds.ToSQL()
fmt.Println(err)
```

Output:
```
name is empty
name is empty
```

<a name="executing"></a>
## Executing Updates

To execute Updates use [`depiq.Database#Update`](https://godoc.org/github.com/orn-id/depiq/#Database.Update) to create your dataset

### Examples

**Executing an update**
```go
db := getDb()

update := db.Update("depiq_user").
	Where(depiq.C("first_name").Eq("Bob")).
	Set(depiq.Record{"first_name": "Bobby"}).
	Executor()

if r, err := update.Exec(); err != nil {
	fmt.Println(err.Error())
} else {
	c, _ := r.RowsAffected()
	fmt.Printf("Updated %d users", c)
}
```

Output:

```
Updated 1 users
```

**Executing with Returning**

```go
db := getDb()

update := db.Update("depiq_user").
	Set(depiq.Record{"last_name": "ucon"}).
	Where(depiq.Ex{"last_name": "Yukon"}).
	Returning("id").
	Executor()

var ids []int64
if err := update.ScanVals(&ids); err != nil {
	fmt.Println(err.Error())
} else {
	fmt.Printf("Updated users with ids %+v", ids)
}

```

Output:
```
Updated users with ids [1 2 3]
```
