// nolint:lll // sql statements are long
package depiq_test

import (
	"fmt"

	"github.com/orn-id/depiq/v9"
	_ "github.com/orn-id/depiq/v9/dialect/mysql"
)

func ExampleUpdate_withStruct() {
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name"`
	}
	sql, args, _ := depiq.Update("items").Set(
		item{Name: "Test", Address: "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
}

func ExampleUpdate_withGoquRecord() {
	sql, args, _ := depiq.Update("items").Set(
		depiq.Record{"name": "Test", "address": "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
}

func ExampleUpdate_withMap() {
	sql, args, _ := depiq.Update("items").Set(
		map[string]interface{}{"name": "Test", "address": "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
}

func ExampleUpdate_withSkipUpdateTag() {
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name" goqu:"skipupdate"`
	}
	sql, args, _ := depiq.Update("items").Set(
		item{Name: "Test", Address: "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr' []
}

func ExampleUpdateDataset_Executor() {
	db := getDB()
	update := db.Update("goqu_user").
		Where(depiq.C("first_name").Eq("Bob")).
		Set(depiq.Record{"first_name": "Bobby"}).
		Executor()

	if r, err := update.Exec(); err != nil {
		fmt.Println(err.Error())
	} else {
		c, _ := r.RowsAffected()
		fmt.Printf("Updated %d users", c)
	}

	// Output:
	// Updated 1 users
}

func ExampleUpdateDataset_Executor_returning() {
	db := getDB()
	var ids []int64
	update := db.Update("goqu_user").
		Set(depiq.Record{"last_name": "ucon"}).
		Where(depiq.Ex{"last_name": "Yukon"}).
		Returning("id").
		Executor()
	if err := update.ScanVals(&ids); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Updated users with ids %+v", ids)
	}

	// Output:
	// Updated users with ids [1 2 3]
}

func ExampleUpdateDataset_Returning() {
	sql, _, _ := depiq.Update("test").
		Set(depiq.Record{"foo": "bar"}).
		Returning("id").
		ToSQL()
	fmt.Println(sql)
	sql, _, _ = depiq.Update("test").
		Set(depiq.Record{"foo": "bar"}).
		Returning(depiq.T("test").All()).
		ToSQL()
	fmt.Println(sql)
	sql, _, _ = depiq.Update("test").
		Set(depiq.Record{"foo": "bar"}).
		Returning("a", "b").
		ToSQL()
	fmt.Println(sql)
	// Output:
	// UPDATE "test" SET "foo"='bar' RETURNING "id"
	// UPDATE "test" SET "foo"='bar' RETURNING "test".*
	// UPDATE "test" SET "foo"='bar' RETURNING "a", "b"
}

func ExampleUpdateDataset_With() {
	sql, _, _ := depiq.Update("test").
		With("some_vals(val)", depiq.From().Select(depiq.L("123"))).
		Where(depiq.C("val").Eq(depiq.From("some_vals").Select("val"))).
		Set(depiq.Record{"name": "Test"}).ToSQL()
	fmt.Println(sql)

	// Output:
	// WITH some_vals(val) AS (SELECT 123) UPDATE "test" SET "name"='Test' WHERE ("val" IN (SELECT "val" FROM "some_vals"))
}

func ExampleUpdateDataset_WithRecursive() {
	sql, _, _ := depiq.Update("nums").
		WithRecursive("nums(x)", depiq.From().Select(depiq.L("1").As("num")).
			UnionAll(depiq.From("nums").
				Select(depiq.L("x+1").As("num")).Where(depiq.C("x").Lt(5)))).
		Set(depiq.Record{"foo": depiq.T("nums").Col("num")}).
		ToSQL()
	fmt.Println(sql)
	// Output:
	// WITH RECURSIVE nums(x) AS (SELECT 1 AS "num" UNION ALL (SELECT x+1 AS "num" FROM "nums" WHERE ("x" < 5))) UPDATE "nums" SET "foo"="nums"."num"
}

func ExampleUpdateDataset_Limit() {
	ds := depiq.Dialect("mysql").
		Update("test").
		Set(depiq.Record{"foo": "bar"}).
		Limit(10)
	sql, _, _ := ds.ToSQL()
	fmt.Println(sql)
	// Output:
	// UPDATE `test` SET `foo`='bar' LIMIT 10
}

func ExampleUpdateDataset_LimitAll() {
	ds := depiq.Dialect("mysql").
		Update("test").
		Set(depiq.Record{"foo": "bar"}).
		LimitAll()
	sql, _, _ := ds.ToSQL()
	fmt.Println(sql)
	// Output:
	// UPDATE `test` SET `foo`='bar' LIMIT ALL
}

func ExampleUpdateDataset_ClearLimit() {
	ds := depiq.Dialect("mysql").
		Update("test").
		Set(depiq.Record{"foo": "bar"}).
		Limit(10)
	sql, _, _ := ds.ClearLimit().ToSQL()
	fmt.Println(sql)
	// Output:
	// UPDATE `test` SET `foo`='bar'
}

func ExampleUpdateDataset_Order() {
	ds := depiq.Dialect("mysql").
		Update("test").
		Set(depiq.Record{"foo": "bar"}).
		Order(depiq.C("a").Asc())
	sql, _, _ := ds.ToSQL()
	fmt.Println(sql)
	// Output:
	// UPDATE `test` SET `foo`='bar' ORDER BY `a` ASC
}

func ExampleUpdateDataset_OrderAppend() {
	ds := depiq.Dialect("mysql").
		Update("test").
		Set(depiq.Record{"foo": "bar"}).
		Order(depiq.C("a").Asc())
	sql, _, _ := ds.OrderAppend(depiq.C("b").Desc().NullsLast()).ToSQL()
	fmt.Println(sql)
	// Output:
	// UPDATE `test` SET `foo`='bar' ORDER BY `a` ASC, `b` DESC NULLS LAST
}

func ExampleUpdateDataset_OrderPrepend() {
	ds := depiq.Dialect("mysql").
		Update("test").
		Set(depiq.Record{"foo": "bar"}).
		Order(depiq.C("a").Asc())

	sql, _, _ := ds.OrderPrepend(depiq.C("b").Desc().NullsLast()).ToSQL()
	fmt.Println(sql)
	// Output:
	// UPDATE `test` SET `foo`='bar' ORDER BY `b` DESC NULLS LAST, `a` ASC
}

func ExampleUpdateDataset_ClearOrder() {
	ds := depiq.Dialect("mysql").
		Update("test").
		Set(depiq.Record{"foo": "bar"}).
		Order(depiq.C("a").Asc())
	sql, _, _ := ds.ClearOrder().ToSQL()
	fmt.Println(sql)
	// Output:
	// UPDATE `test` SET `foo`='bar'
}

func ExampleUpdateDataset_From() {
	ds := depiq.Update("table_one").
		Set(depiq.Record{"foo": depiq.I("table_two.bar")}).
		From("table_two").
		Where(depiq.Ex{"table_one.id": depiq.I("table_two.id")})

	sql, _, _ := ds.ToSQL()
	fmt.Println(sql)
	// Output:
	// UPDATE "table_one" SET "foo"="table_two"."bar" FROM "table_two" WHERE ("table_one"."id" = "table_two"."id")
}

func ExampleUpdateDataset_From_postgres() {
	dialect := depiq.Dialect("postgres")

	ds := dialect.Update("table_one").
		Set(depiq.Record{"foo": depiq.I("table_two.bar")}).
		From("table_two").
		Where(depiq.Ex{"table_one.id": depiq.I("table_two.id")})

	sql, _, _ := ds.ToSQL()
	fmt.Println(sql)
	// Output:
	// UPDATE "table_one" SET "foo"="table_two"."bar" FROM "table_two" WHERE ("table_one"."id" = "table_two"."id")
}

func ExampleUpdateDataset_From_mysql() {
	dialect := depiq.Dialect("mysql")

	ds := dialect.Update("table_one").
		Set(depiq.Record{"foo": depiq.I("table_two.bar")}).
		From("table_two").
		Where(depiq.Ex{"table_one.id": depiq.I("table_two.id")})

	sql, _, _ := ds.ToSQL()
	fmt.Println(sql)
	// Output:
	// UPDATE `table_one`,`table_two` SET `foo`=`table_two`.`bar` WHERE (`table_one`.`id` = `table_two`.`id`)
}

func ExampleUpdateDataset_Where() {
	// By default everything is anded together
	sql, _, _ := depiq.Update("test").
		Set(depiq.Record{"foo": "bar"}).
		Where(depiq.Ex{
			"a": depiq.Op{"gt": 10},
			"b": depiq.Op{"lt": 10},
			"c": nil,
			"d": []string{"a", "b", "c"},
		}).ToSQL()
	fmt.Println(sql)
	// You can use ExOr to get ORed expressions together
	sql, _, _ = depiq.Update("test").
		Set(depiq.Record{"foo": "bar"}).
		Where(depiq.ExOr{
			"a": depiq.Op{"gt": 10},
			"b": depiq.Op{"lt": 10},
			"c": nil,
			"d": []string{"a", "b", "c"},
		}).ToSQL()
	fmt.Println(sql)
	// You can use Or with Ex to Or multiple Ex maps together
	sql, _, _ = depiq.Update("test").
		Set(depiq.Record{"foo": "bar"}).
		Where(
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
	// By default everything is anded together
	sql, _, _ = depiq.Update("test").
		Set(depiq.Record{"foo": "bar"}).
		Where(
			depiq.C("a").Gt(10),
			depiq.C("b").Lt(10),
			depiq.C("c").IsNull(),
			depiq.C("d").In("a", "b", "c"),
		).ToSQL()
	fmt.Println(sql)
	// You can use a combination of Ors and Ands
	sql, _, _ = depiq.Update("test").
		Set(depiq.Record{"foo": "bar"}).
		Where(
			depiq.Or(
				depiq.C("a").Gt(10),
				depiq.And(
					depiq.C("b").Lt(10),
					depiq.C("c").IsNull(),
				),
			),
		).ToSQL()
	fmt.Println(sql)
	// Output:
	// UPDATE "test" SET "foo"='bar' WHERE (("a" > 10) AND ("b" < 10) AND ("c" IS NULL) AND ("d" IN ('a', 'b', 'c')))
	// UPDATE "test" SET "foo"='bar' WHERE (("a" > 10) OR ("b" < 10) OR ("c" IS NULL) OR ("d" IN ('a', 'b', 'c')))
	// UPDATE "test" SET "foo"='bar' WHERE ((("a" > 10) AND ("b" < 10)) OR (("c" IS NULL) AND ("d" IN ('a', 'b', 'c'))))
	// UPDATE "test" SET "foo"='bar' WHERE (("a" > 10) AND ("b" < 10) AND ("c" IS NULL) AND ("d" IN ('a', 'b', 'c')))
	// UPDATE "test" SET "foo"='bar' WHERE (("a" > 10) OR (("b" < 10) AND ("c" IS NULL)))
}

func ExampleUpdateDataset_Where_prepared() {
	// By default everything is anded together
	sql, args, _ := depiq.Update("test").
		Prepared(true).
		Set(depiq.Record{"foo": "bar"}).
		Where(depiq.Ex{
			"a": depiq.Op{"gt": 10},
			"b": depiq.Op{"lt": 10},
			"c": nil,
			"d": []string{"a", "b", "c"},
		}).ToSQL()
	fmt.Println(sql, args)
	// You can use ExOr to get ORed expressions together
	sql, args, _ = depiq.Update("test").Prepared(true).
		Set(depiq.Record{"foo": "bar"}).
		Where(depiq.ExOr{
			"a": depiq.Op{"gt": 10},
			"b": depiq.Op{"lt": 10},
			"c": nil,
			"d": []string{"a", "b", "c"},
		}).ToSQL()
	fmt.Println(sql, args)
	// You can use Or with Ex to Or multiple Ex maps together
	sql, args, _ = depiq.Update("test").Prepared(true).
		Set(depiq.Record{"foo": "bar"}).
		Where(
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
	fmt.Println(sql, args)
	// By default everything is anded together
	sql, args, _ = depiq.Update("test").Prepared(true).
		Set(depiq.Record{"foo": "bar"}).
		Where(
			depiq.C("a").Gt(10),
			depiq.C("b").Lt(10),
			depiq.C("c").IsNull(),
			depiq.C("d").In("a", "b", "c"),
		).ToSQL()
	fmt.Println(sql, args)
	// You can use a combination of Ors and Ands
	sql, args, _ = depiq.Update("test").Prepared(true).
		Set(depiq.Record{"foo": "bar"}).
		Where(
			depiq.Or(
				depiq.C("a").Gt(10),
				depiq.And(
					depiq.C("b").Lt(10),
					depiq.C("c").IsNull(),
				),
			),
		).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// UPDATE "test" SET "foo"=? WHERE (("a" > ?) AND ("b" < ?) AND ("c" IS NULL) AND ("d" IN (?, ?, ?))) [bar 10 10 a b c]
	// UPDATE "test" SET "foo"=? WHERE (("a" > ?) OR ("b" < ?) OR ("c" IS NULL) OR ("d" IN (?, ?, ?))) [bar 10 10 a b c]
	// UPDATE "test" SET "foo"=? WHERE ((("a" > ?) AND ("b" < ?)) OR (("c" IS NULL) AND ("d" IN (?, ?, ?)))) [bar 10 10 a b c]
	// UPDATE "test" SET "foo"=? WHERE (("a" > ?) AND ("b" < ?) AND ("c" IS NULL) AND ("d" IN (?, ?, ?))) [bar 10 10 a b c]
	// UPDATE "test" SET "foo"=? WHERE (("a" > ?) OR (("b" < ?) AND ("c" IS NULL))) [bar 10 10]
}

func ExampleUpdateDataset_ClearWhere() {
	ds := depiq.
		Update("test").
		Set(depiq.Record{"foo": "bar"}).
		Where(
			depiq.Or(
				depiq.C("a").Gt(10),
				depiq.And(
					depiq.C("b").Lt(10),
					depiq.C("c").IsNull(),
				),
			),
		)
	sql, _, _ := ds.ClearWhere().ToSQL()
	fmt.Println(sql)
	// Output:
	// UPDATE "test" SET "foo"='bar'
}

func ExampleUpdateDataset_Table() {
	ds := depiq.Update("test")
	sql, _, _ := ds.Table("test2").Set(depiq.Record{"foo": "bar"}).ToSQL()
	fmt.Println(sql)
	// Output:
	// UPDATE "test2" SET "foo"='bar'
}

func ExampleUpdateDataset_Table_aliased() {
	ds := depiq.Update("test")
	sql, _, _ := ds.Table(depiq.T("test").As("t")).Set(depiq.Record{"foo": "bar"}).ToSQL()
	fmt.Println(sql)
	// Output:
	// UPDATE "test" AS "t" SET "foo"='bar'
}

func ExampleUpdateDataset_Set() {
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name"`
	}
	sql, args, _ := depiq.Update("items").Set(
		item{Name: "Test", Address: "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.Update("items").Set(
		depiq.Record{"name": "Test", "address": "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.Update("items").Set(
		map[string]interface{}{"name": "Test", "address": "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
}

func ExampleUpdateDataset_Set_struct() {
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name"`
	}
	sql, args, _ := depiq.Update("items").Set(
		item{Name: "Test", Address: "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
}

func ExampleUpdateDataset_Set_goquRecord() {
	sql, args, _ := depiq.Update("items").Set(
		depiq.Record{"name": "Test", "address": "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
}

func ExampleUpdateDataset_Set_map() {
	sql, args, _ := depiq.Update("items").Set(
		map[string]interface{}{"name": "Test", "address": "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
}

func ExampleUpdateDataset_Set_withSkipUpdateTag() {
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name" goqu:"skipupdate"`
	}
	sql, args, _ := depiq.Update("items").Set(
		item{Name: "Test", Address: "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr' []
}

func ExampleUpdateDataset_Set_withDefaultIfEmptyTag() {
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name" goqu:"defaultifempty"`
	}
	sql, args, _ := depiq.Update("items").Set(
		item{Address: "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.Update("items").Set(
		item{Name: "Bob Yukon", Address: "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"=DEFAULT []
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Bob Yukon' []
}

func ExampleUpdateDataset_Set_withNoTags() {
	type item struct {
		Address string
		Name    string
	}
	sql, args, _ := depiq.Update("items").Set(
		item{Name: "Test", Address: "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
}

func ExampleUpdateDataset_Set_withEmbeddedStruct() {
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

	// Output:
	// UPDATE "user" SET "address_state"='NY',"address_street"='111 Street',"firstname"='Greg',"lastname"='Farley' []
}

func ExampleUpdateDataset_Set_withIgnoredEmbedded() {
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

	// Output:
	// UPDATE "user" SET "firstname"='Greg',"lastname"='Farley' []
}

func ExampleUpdateDataset_Set_withNilEmbeddedPointer() {
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

	// Output:
	// UPDATE "user" SET "firstname"='Greg',"lastname"='Farley' []
}

func ExampleUpdateDataset_ToSQL_prepared() {
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name"`
	}

	sql, args, _ := depiq.From("items").Prepared(true).Update().Set(
		item{Name: "Test", Address: "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("items").Prepared(true).Update().Set(
		depiq.Record{"name": "Test", "address": "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("items").Prepared(true).Update().Set(
		map[string]interface{}{"name": "Test", "address": "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// UPDATE "items" SET "address"=?,"name"=? [111 Test Addr Test]
	// UPDATE "items" SET "address"=?,"name"=? [111 Test Addr Test]
	// UPDATE "items" SET "address"=?,"name"=? [111 Test Addr Test]
}

func ExampleUpdateDataset_Prepared() {
	sql, args, _ := depiq.Update("items").Prepared(true).Set(
		depiq.Record{"name": "Test", "address": "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"=?,"name"=? [111 Test Addr Test]
}
