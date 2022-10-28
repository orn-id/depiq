package depiq_test

import (
	"fmt"

	"github.com/orn-id/depiq"
	_ "github.com/orn-id/depiq/dialect/mysql"
)

func ExampleDelete() {
	ds := depiq.Delete("items")

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	// Output:
	// DELETE FROM "items" []
}

func ExampleDeleteDataset_Executor() {
	db := getDB()

	de := db.Delete("goqu_user").
		Where(depiq.Ex{"first_name": "Bob"}).
		Executor()
	if r, err := de.Exec(); err != nil {
		fmt.Println(err.Error())
	} else {
		c, _ := r.RowsAffected()
		fmt.Printf("Deleted %d users", c)
	}

	// Output:
	// Deleted 1 users
}

func ExampleDeleteDataset_Executor_returning() {
	db := getDB()

	de := db.Delete("goqu_user").
		Where(depiq.C("last_name").Eq("Yukon")).
		Returning(depiq.C("id")).
		Executor()

	var ids []int64
	if err := de.ScanVals(&ids); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Deleted users [ids:=%+v]", ids)
	}

	// Output:
	// Deleted users [ids:=[1 2 3]]
}

func ExampleDeleteDataset_With() {
	sql, _, _ := depiq.Delete("test").
		With("check_vals(val)", depiq.From().Select(depiq.L("123"))).
		Where(depiq.C("val").Eq(depiq.From("check_vals").Select("val"))).
		ToSQL()
	fmt.Println(sql)

	// Output:
	// WITH check_vals(val) AS (SELECT 123) DELETE FROM "test" WHERE ("val" IN (SELECT "val" FROM "check_vals"))
}

func ExampleDeleteDataset_WithRecursive() {
	sql, _, _ := depiq.Delete("nums").
		WithRecursive("nums(x)",
			depiq.From().Select(depiq.L("1")).
				UnionAll(depiq.From("nums").
					Select(depiq.L("x+1")).Where(depiq.C("x").Lt(5)))).
		ToSQL()
	fmt.Println(sql)
	// Output:
	// WITH RECURSIVE nums(x) AS (SELECT 1 UNION ALL (SELECT x+1 FROM "nums" WHERE ("x" < 5))) DELETE FROM "nums"
}

func ExampleDeleteDataset_Where() {
	// By default everything is anded together
	sql, _, _ := depiq.Delete("test").Where(depiq.Ex{
		"a": depiq.Op{"gt": 10},
		"b": depiq.Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).ToSQL()
	fmt.Println(sql)
	// You can use ExOr to get ORed expressions together
	sql, _, _ = depiq.Delete("test").Where(depiq.ExOr{
		"a": depiq.Op{"gt": 10},
		"b": depiq.Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).ToSQL()
	fmt.Println(sql)
	// You can use Or with Ex to Or multiple Ex maps together
	sql, _, _ = depiq.Delete("test").Where(
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
	sql, _, _ = depiq.Delete("test").Where(
		depiq.C("a").Gt(10),
		depiq.C("b").Lt(10),
		depiq.C("c").IsNull(),
		depiq.C("d").In("a", "b", "c"),
	).ToSQL()
	fmt.Println(sql)
	// You can use a combination of Ors and Ands
	sql, _, _ = depiq.Delete("test").Where(
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
	// DELETE FROM "test" WHERE (("a" > 10) AND ("b" < 10) AND ("c" IS NULL) AND ("d" IN ('a', 'b', 'c')))
	// DELETE FROM "test" WHERE (("a" > 10) OR ("b" < 10) OR ("c" IS NULL) OR ("d" IN ('a', 'b', 'c')))
	// DELETE FROM "test" WHERE ((("a" > 10) AND ("b" < 10)) OR (("c" IS NULL) AND ("d" IN ('a', 'b', 'c'))))
	// DELETE FROM "test" WHERE (("a" > 10) AND ("b" < 10) AND ("c" IS NULL) AND ("d" IN ('a', 'b', 'c')))
	// DELETE FROM "test" WHERE (("a" > 10) OR (("b" < 10) AND ("c" IS NULL)))
}

func ExampleDeleteDataset_Where_prepared() {
	// By default everything is anded together
	sql, args, _ := depiq.Delete("test").Prepared(true).Where(depiq.Ex{
		"a": depiq.Op{"gt": 10},
		"b": depiq.Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).ToSQL()
	fmt.Println(sql, args)
	// You can use ExOr to get ORed expressions together
	sql, args, _ = depiq.Delete("test").Prepared(true).Where(depiq.ExOr{
		"a": depiq.Op{"gt": 10},
		"b": depiq.Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).ToSQL()
	fmt.Println(sql, args)
	// You can use Or with Ex to Or multiple Ex maps together
	sql, args, _ = depiq.Delete("test").Prepared(true).Where(
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
	sql, args, _ = depiq.Delete("test").Prepared(true).Where(
		depiq.C("a").Gt(10),
		depiq.C("b").Lt(10),
		depiq.C("c").IsNull(),
		depiq.C("d").In("a", "b", "c"),
	).ToSQL()
	fmt.Println(sql, args)
	// You can use a combination of Ors and Ands
	sql, args, _ = depiq.Delete("test").Prepared(true).Where(
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
	// DELETE FROM "test" WHERE (("a" > ?) AND ("b" < ?) AND ("c" IS NULL) AND ("d" IN (?, ?, ?))) [10 10 a b c]
	// DELETE FROM "test" WHERE (("a" > ?) OR ("b" < ?) OR ("c" IS NULL) OR ("d" IN (?, ?, ?))) [10 10 a b c]
	// DELETE FROM "test" WHERE ((("a" > ?) AND ("b" < ?)) OR (("c" IS NULL) AND ("d" IN (?, ?, ?)))) [10 10 a b c]
	// DELETE FROM "test" WHERE (("a" > ?) AND ("b" < ?) AND ("c" IS NULL) AND ("d" IN (?, ?, ?))) [10 10 a b c]
	// DELETE FROM "test" WHERE (("a" > ?) OR (("b" < ?) AND ("c" IS NULL))) [10 10]
}

func ExampleDeleteDataset_ClearWhere() {
	ds := depiq.Delete("test").Where(
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
	// DELETE FROM "test"
}

func ExampleDeleteDataset_Limit() {
	ds := depiq.Dialect("mysql").Delete("test").Limit(10)
	sql, _, _ := ds.ToSQL()
	fmt.Println(sql)
	// Output:
	// DELETE FROM `test` LIMIT 10
}

func ExampleDeleteDataset_LimitAll() {
	// Using mysql dialect because it supports limit on delete
	ds := depiq.Dialect("mysql").Delete("test").LimitAll()
	sql, _, _ := ds.ToSQL()
	fmt.Println(sql)
	// Output:
	// DELETE FROM `test` LIMIT ALL
}

func ExampleDeleteDataset_ClearLimit() {
	// Using mysql dialect because it supports limit on delete
	ds := depiq.Dialect("mysql").Delete("test").Limit(10)
	sql, _, _ := ds.ClearLimit().ToSQL()
	fmt.Println(sql)
	// Output:
	// DELETE `test` FROM `test`
}

func ExampleDeleteDataset_Order() {
	// use mysql dialect because it supports order by on deletes
	ds := depiq.Dialect("mysql").Delete("test").Order(depiq.C("a").Asc())
	sql, _, _ := ds.ToSQL()
	fmt.Println(sql)
	// Output:
	// DELETE FROM `test` ORDER BY `a` ASC
}

func ExampleDeleteDataset_OrderAppend() {
	// use mysql dialect because it supports order by on deletes
	ds := depiq.Dialect("mysql").Delete("test").Order(depiq.C("a").Asc())
	sql, _, _ := ds.OrderAppend(depiq.C("b").Desc().NullsLast()).ToSQL()
	fmt.Println(sql)
	// Output:
	// DELETE FROM `test` ORDER BY `a` ASC, `b` DESC NULLS LAST
}

func ExampleDeleteDataset_OrderPrepend() {
	// use mysql dialect because it supports order by on deletes
	ds := depiq.Dialect("mysql").Delete("test").Order(depiq.C("a").Asc())
	sql, _, _ := ds.OrderPrepend(depiq.C("b").Desc().NullsLast()).ToSQL()
	fmt.Println(sql)
	// Output:
	// DELETE FROM `test` ORDER BY `b` DESC NULLS LAST, `a` ASC
}

func ExampleDeleteDataset_ClearOrder() {
	ds := depiq.Delete("test").Order(depiq.C("a").Asc())
	sql, _, _ := ds.ClearOrder().ToSQL()
	fmt.Println(sql)
	// Output:
	// DELETE FROM "test"
}

func ExampleDeleteDataset_ToSQL() {
	sql, args, _ := depiq.Delete("items").ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.Delete("items").
		Where(depiq.Ex{"id": depiq.Op{"gt": 10}}).
		ToSQL()
	fmt.Println(sql, args)

	// Output:
	// DELETE FROM "items" []
	// DELETE FROM "items" WHERE ("id" > 10) []
}

func ExampleDeleteDataset_Prepared() {
	sql, args, _ := depiq.Delete("items").Prepared(true).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.Delete("items").
		Prepared(true).
		Where(depiq.Ex{"id": depiq.Op{"gt": 10}}).
		ToSQL()
	fmt.Println(sql, args)

	// Output:
	// DELETE FROM "items" []
	// DELETE FROM "items" WHERE ("id" > ?) [10]
}

func ExampleDeleteDataset_Returning() {
	ds := depiq.Delete("items")
	sql, args, _ := ds.Returning("id").ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Returning("id").Where(depiq.C("id").IsNotNull()).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// DELETE FROM "items" RETURNING "id" []
	// DELETE FROM "items" WHERE ("id" IS NOT NULL) RETURNING "id" []
}
