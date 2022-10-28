// nolint:lll // sql statements are long
package depiq_test

import (
	goSQL "database/sql"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/lib/pq"
	"github.com/orn-id/depiq"
	"github.com/orn-id/depiq/exp"
)

const schema = `
		DROP TABLE IF EXISTS "user_role";
		DROP TABLE IF EXISTS "depiq_user";	
		CREATE  TABLE "depiq_user" (
			"id" SERIAL PRIMARY KEY NOT NULL,
			"first_name" VARCHAR(45) NOT NULL,
			"last_name" VARCHAR(45) NOT NULL,
			"created" TIMESTAMP NOT NULL DEFAULT now()
		);
		CREATE  TABLE "user_role" (
			"id" SERIAL PRIMARY KEY NOT NULL,
			"user_id" BIGINT NOT NULL REFERENCES depiq_user(id) ON DELETE CASCADE,
			"name" VARCHAR(45) NOT NULL,
			"created" TIMESTAMP NOT NULL DEFAULT now()
		); 
    `

const defaultDBURI = "postgres://postgres:@localhost:5435/goqupostgres?sslmode=disable"

var goquDB *depiq.Database

func getDB() *depiq.Database {
	if goquDB == nil {
		dbURI := os.Getenv("PG_URI")
		if dbURI == "" {
			dbURI = defaultDBURI
		}
		uri, err := pq.ParseURL(dbURI)
		if err != nil {
			panic(err)
		}
		pdb, err := goSQL.Open("postgres", uri)
		if err != nil {
			panic(err)
		}
		goquDB = depiq.New("postgres", pdb)
	}
	// reset the db
	if _, err := goquDB.Exec(schema); err != nil {
		panic(err)
	}
	type goquUser struct {
		ID        int64     `db:"id" goqu:"skipinsert"`
		FirstName string    `db:"first_name"`
		LastName  string    `db:"last_name"`
		Created   time.Time `db:"created" goqu:"skipupdate"`
	}

	users := []goquUser{
		{FirstName: "Bob", LastName: "Yukon"},
		{FirstName: "Sally", LastName: "Yukon"},
		{FirstName: "Vinita", LastName: "Yukon"},
		{FirstName: "John", LastName: "Doe"},
	}
	var userIds []int64
	err := goquDB.Insert("depiq_user").Rows(users).Returning("id").Executor().ScanVals(&userIds)
	if err != nil {
		panic(err)
	}
	type userRole struct {
		ID      int64     `db:"id" goqu:"skipinsert"`
		UserID  int64     `db:"user_id"`
		Name    string    `db:"name"`
		Created time.Time `db:"created" goqu:"skipupdate"`
	}

	roles := []userRole{
		{UserID: userIds[0], Name: "Admin"},
		{UserID: userIds[1], Name: "Manager"},
		{UserID: userIds[2], Name: "Manager"},
		{UserID: userIds[3], Name: "User"},
	}
	_, err = goquDB.Insert("user_role").Rows(roles).Executor().Exec()
	if err != nil {
		panic(err)
	}
	return goquDB
}

func ExampleSelectDataset() {
	ds := depiq.From("test").
		Select(depiq.COUNT("*")).
		InnerJoin(depiq.T("test2"), depiq.On(depiq.I("test.fkey").Eq(depiq.I("test2.id")))).
		LeftJoin(depiq.T("test3"), depiq.On(depiq.I("test2.fkey").Eq(depiq.I("test3.id")))).
		Where(
			depiq.Ex{
				"test.name": depiq.Op{
					"like": regexp.MustCompile("^[ab]"),
				},
				"test2.amount": depiq.Op{
					"isNot": nil,
				},
			},
			depiq.ExOr{
				"test3.id":     nil,
				"test3.status": []string{"passed", "active", "registered"},
			}).
		Order(depiq.I("test.created").Desc().NullsLast()).
		GroupBy(depiq.I("test.user_id")).
		Having(depiq.AVG("test3.age").Gt(10))

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// nolint:lll // SQL statements are long
	// Output:
	// SELECT COUNT(*) FROM "test" INNER JOIN "test2" ON ("test"."fkey" = "test2"."id") LEFT JOIN "test3" ON ("test2"."fkey" = "test3"."id") WHERE ((("test"."name" ~ '^[ab]') AND ("test2"."amount" IS NOT NULL)) AND (("test3"."id" IS NULL) OR ("test3"."status" IN ('passed', 'active', 'registered')))) GROUP BY "test"."user_id" HAVING (AVG("test3"."age") > 10) ORDER BY "test"."created" DESC NULLS LAST []
	// SELECT COUNT(*) FROM "test" INNER JOIN "test2" ON ("test"."fkey" = "test2"."id") LEFT JOIN "test3" ON ("test2"."fkey" = "test3"."id") WHERE ((("test"."name" ~ ?) AND ("test2"."amount" IS NOT NULL)) AND (("test3"."id" IS NULL) OR ("test3"."status" IN (?, ?, ?)))) GROUP BY "test"."user_id" HAVING (AVG("test3"."age") > ?) ORDER BY "test"."created" DESC NULLS LAST [^[ab] passed active registered 10]
}

func ExampleSelect() {
	sql, _, _ := depiq.Select(depiq.L("NOW()")).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT NOW()
}

func ExampleFrom() {
	sql, args, _ := depiq.From("test").ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" []
}

func ExampleSelectDataset_As() {
	ds := depiq.From("test").As("t")
	sql, _, _ := depiq.From(ds).ToSQL()
	fmt.Println(sql)
	// Output: SELECT * FROM (SELECT * FROM "test") AS "t"
}

func ExampleSelectDataset_Union() {
	sql, _, _ := depiq.From("test").
		Union(depiq.From("test2")).
		ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").
		Limit(1).
		Union(depiq.From("test2")).
		ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").
		Limit(1).
		Union(depiq.From("test2").
			Order(depiq.C("id").Desc())).
		ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" UNION (SELECT * FROM "test2")
	// SELECT * FROM (SELECT * FROM "test" LIMIT 1) AS "t1" UNION (SELECT * FROM "test2")
	// SELECT * FROM (SELECT * FROM "test" LIMIT 1) AS "t1" UNION (SELECT * FROM (SELECT * FROM "test2" ORDER BY "id" DESC) AS "t1")
}

func ExampleSelectDataset_UnionAll() {
	sql, _, _ := depiq.From("test").
		UnionAll(depiq.From("test2")).
		ToSQL()
	fmt.Println(sql)
	sql, _, _ = depiq.From("test").
		Limit(1).
		UnionAll(depiq.From("test2")).
		ToSQL()
	fmt.Println(sql)
	sql, _, _ = depiq.From("test").
		Limit(1).
		UnionAll(depiq.From("test2").
			Order(depiq.C("id").Desc())).
		ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" UNION ALL (SELECT * FROM "test2")
	// SELECT * FROM (SELECT * FROM "test" LIMIT 1) AS "t1" UNION ALL (SELECT * FROM "test2")
	// SELECT * FROM (SELECT * FROM "test" LIMIT 1) AS "t1" UNION ALL (SELECT * FROM (SELECT * FROM "test2" ORDER BY "id" DESC) AS "t1")
}

func ExampleSelectDataset_With() {
	sql, _, _ := depiq.From("one").
		With("one", depiq.From().Select(depiq.L("1"))).
		Select(depiq.Star()).
		ToSQL()
	fmt.Println(sql)
	sql, _, _ = depiq.From("derived").
		With("intermed", depiq.From("test").Select(depiq.Star()).Where(depiq.C("x").Gte(5))).
		With("derived", depiq.From("intermed").Select(depiq.Star()).Where(depiq.C("x").Lt(10))).
		Select(depiq.Star()).
		ToSQL()
	fmt.Println(sql)
	sql, _, _ = depiq.From("multi").
		With("multi(x,y)", depiq.From().Select(depiq.L("1"), depiq.L("2"))).
		Select(depiq.C("x"), depiq.C("y")).
		ToSQL()
	fmt.Println(sql)

	// Output:
	// WITH one AS (SELECT 1) SELECT * FROM "one"
	// WITH intermed AS (SELECT * FROM "test" WHERE ("x" >= 5)), derived AS (SELECT * FROM "intermed" WHERE ("x" < 10)) SELECT * FROM "derived"
	// WITH multi(x,y) AS (SELECT 1, 2) SELECT "x", "y" FROM "multi"
}

func ExampleSelectDataset_With_insertDataset() {
	insertDs := depiq.Insert("foo").Rows(depiq.Record{"user_id": 10}).Returning("id")

	ds := depiq.From("bar").
		With("ins", insertDs).
		Select("bar_name").
		Where(depiq.Ex{"bar.user_id": depiq.I("ins.user_id")})

	sql, _, _ := ds.ToSQL()
	fmt.Println(sql)

	sql, args, _ := ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// WITH ins AS (INSERT INTO "foo" ("user_id") VALUES (10) RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "ins"."user_id")
	// WITH ins AS (INSERT INTO "foo" ("user_id") VALUES (?) RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "ins"."user_id") [10]
}

func ExampleSelectDataset_With_updateDataset() {
	updateDs := depiq.Update("foo").Set(depiq.Record{"bar": "baz"}).Returning("id")

	ds := depiq.From("bar").
		With("upd", updateDs).
		Select("bar_name").
		Where(depiq.Ex{"bar.user_id": depiq.I("upd.user_id")})

	sql, _, _ := ds.ToSQL()
	fmt.Println(sql)

	sql, args, _ := ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// WITH upd AS (UPDATE "foo" SET "bar"='baz' RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "upd"."user_id")
	// WITH upd AS (UPDATE "foo" SET "bar"=? RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "upd"."user_id") [baz]
}

func ExampleSelectDataset_With_deleteDataset() {
	deleteDs := depiq.Delete("foo").Where(depiq.Ex{"bar": "baz"}).Returning("id")

	ds := depiq.From("bar").
		With("del", deleteDs).
		Select("bar_name").
		Where(depiq.Ex{"bar.user_id": depiq.I("del.user_id")})

	sql, _, _ := ds.ToSQL()
	fmt.Println(sql)

	sql, args, _ := ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// WITH del AS (DELETE FROM "foo" WHERE ("bar" = 'baz') RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "del"."user_id")
	// WITH del AS (DELETE FROM "foo" WHERE ("bar" = ?) RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "del"."user_id") [baz]
}

func ExampleSelectDataset_WithRecursive() {
	sql, _, _ := depiq.From("nums").
		WithRecursive("nums(x)",
			depiq.From().Select(depiq.L("1")).
				UnionAll(depiq.From("nums").
					Select(depiq.L("x+1")).Where(depiq.C("x").Lt(5)))).
		ToSQL()
	fmt.Println(sql)
	// Output:
	// WITH RECURSIVE nums(x) AS (SELECT 1 UNION ALL (SELECT x+1 FROM "nums" WHERE ("x" < 5))) SELECT * FROM "nums"
}

func ExampleSelectDataset_Intersect() {
	sql, _, _ := depiq.From("test").
		Intersect(depiq.From("test2")).
		ToSQL()
	fmt.Println(sql)
	sql, _, _ = depiq.From("test").
		Limit(1).
		Intersect(depiq.From("test2")).
		ToSQL()
	fmt.Println(sql)
	sql, _, _ = depiq.From("test").
		Limit(1).
		Intersect(depiq.From("test2").
			Order(depiq.C("id").Desc())).
		ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" INTERSECT (SELECT * FROM "test2")
	// SELECT * FROM (SELECT * FROM "test" LIMIT 1) AS "t1" INTERSECT (SELECT * FROM "test2")
	// SELECT * FROM (SELECT * FROM "test" LIMIT 1) AS "t1" INTERSECT (SELECT * FROM (SELECT * FROM "test2" ORDER BY "id" DESC) AS "t1")
}

func ExampleSelectDataset_IntersectAll() {
	sql, _, _ := depiq.From("test").
		IntersectAll(depiq.From("test2")).
		ToSQL()
	fmt.Println(sql)
	sql, _, _ = depiq.From("test").
		Limit(1).
		IntersectAll(depiq.From("test2")).
		ToSQL()
	fmt.Println(sql)
	sql, _, _ = depiq.From("test").
		Limit(1).
		IntersectAll(depiq.From("test2").
			Order(depiq.C("id").Desc())).
		ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" INTERSECT ALL (SELECT * FROM "test2")
	// SELECT * FROM (SELECT * FROM "test" LIMIT 1) AS "t1" INTERSECT ALL (SELECT * FROM "test2")
	// SELECT * FROM (SELECT * FROM "test" LIMIT 1) AS "t1" INTERSECT ALL (SELECT * FROM (SELECT * FROM "test2" ORDER BY "id" DESC) AS "t1")
}

func ExampleSelectDataset_ClearOffset() {
	ds := depiq.From("test").
		Offset(2)
	sql, _, _ := ds.
		ClearOffset().
		ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test"
}

func ExampleSelectDataset_Offset() {
	ds := depiq.From("test").Offset(2)
	sql, _, _ := ds.ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" OFFSET 2
}

func ExampleSelectDataset_Limit() {
	ds := depiq.From("test").Limit(10)
	sql, _, _ := ds.ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" LIMIT 10
}

func ExampleSelectDataset_LimitAll() {
	ds := depiq.From("test").LimitAll()
	sql, _, _ := ds.ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" LIMIT ALL
}

func ExampleSelectDataset_ClearLimit() {
	ds := depiq.From("test").Limit(10)
	sql, _, _ := ds.ClearLimit().ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test"
}

func ExampleSelectDataset_Order() {
	ds := depiq.From("test").Order(depiq.C("a").Asc())
	sql, _, _ := ds.ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" ORDER BY "a" ASC
}

func ExampleSelectDataset_Order_caseExpression() {
	ds := depiq.From("test").Order(depiq.Case().When(depiq.C("num").Gt(10), 0).Else(1).Asc())
	sql, _, _ := ds.ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" ORDER BY CASE  WHEN ("num" > 10) THEN 0 ELSE 1 END ASC
}

func ExampleSelectDataset_OrderAppend() {
	ds := depiq.From("test").Order(depiq.C("a").Asc())
	sql, _, _ := ds.OrderAppend(depiq.C("b").Desc().NullsLast()).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" ORDER BY "a" ASC, "b" DESC NULLS LAST
}

func ExampleSelectDataset_OrderPrepend() {
	ds := depiq.From("test").Order(depiq.C("a").Asc())
	sql, _, _ := ds.OrderPrepend(depiq.C("b").Desc().NullsLast()).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" ORDER BY "b" DESC NULLS LAST, "a" ASC
}

func ExampleSelectDataset_ClearOrder() {
	ds := depiq.From("test").Order(depiq.C("a").Asc())
	sql, _, _ := ds.ClearOrder().ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test"
}

func ExampleSelectDataset_GroupBy() {
	sql, _, _ := depiq.From("test").
		Select(depiq.SUM("income").As("income_sum")).
		GroupBy("age").
		ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT SUM("income") AS "income_sum" FROM "test" GROUP BY "age"
}

func ExampleSelectDataset_GroupByAppend() {
	ds := depiq.From("test").
		Select(depiq.SUM("income").As("income_sum")).
		GroupBy("age")
	sql, _, _ := ds.
		GroupByAppend("job").
		ToSQL()
	fmt.Println(sql)
	// the original dataset group by does not change
	sql, _, _ = ds.ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT SUM("income") AS "income_sum" FROM "test" GROUP BY "age", "job"
	// SELECT SUM("income") AS "income_sum" FROM "test" GROUP BY "age"
}

func ExampleSelectDataset_Having() {
	sql, _, _ := depiq.From("test").Having(depiq.SUM("income").Gt(1000)).ToSQL()
	fmt.Println(sql)
	sql, _, _ = depiq.From("test").GroupBy("age").Having(depiq.SUM("income").Gt(1000)).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" HAVING (SUM("income") > 1000)
	// SELECT * FROM "test" GROUP BY "age" HAVING (SUM("income") > 1000)
}

func ExampleSelectDataset_Window() {
	ds := depiq.From("test").
		Select(depiq.ROW_NUMBER().Over(depiq.W().PartitionBy("a").OrderBy(depiq.I("b").Asc())))
	query, args, _ := ds.ToSQL()
	fmt.Println(query, args)

	ds = depiq.From("test").
		Select(depiq.ROW_NUMBER().OverName(depiq.I("w"))).
		Window(depiq.W("w").PartitionBy("a").OrderBy(depiq.I("b").Asc()))
	query, args, _ = ds.ToSQL()
	fmt.Println(query, args)

	ds = depiq.From("test").
		Select(depiq.ROW_NUMBER().OverName(depiq.I("w1"))).
		Window(
			depiq.W("w1").PartitionBy("a"),
			depiq.W("w").Inherit("w1").OrderBy(depiq.I("b").Asc()),
		)
	query, args, _ = ds.ToSQL()
	fmt.Println(query, args)

	ds = depiq.From("test").
		Select(depiq.ROW_NUMBER().Over(depiq.W().Inherit("w").OrderBy("b"))).
		Window(depiq.W("w").PartitionBy("a"))
	query, args, _ = ds.ToSQL()
	fmt.Println(query, args)
	// Output
	// SELECT ROW_NUMBER() OVER (PARTITION BY "a" ORDER BY "b" ASC) FROM "test" []
	// SELECT ROW_NUMBER() OVER "w" FROM "test" WINDOW "w" AS (PARTITION BY "a" ORDER BY "b" ASC) []
	// SELECT ROW_NUMBER() OVER "w" FROM "test" WINDOW "w1" AS (PARTITION BY "a"), "w" AS ("w1" ORDER BY "b" ASC) []
	// SELECT ROW_NUMBER() OVER ("w" ORDER BY "b") FROM "test" WINDOW "w" AS (PARTITION BY "a") []
}

func ExampleSelectDataset_Where() {
	// By default everything is anded together
	sql, _, _ := depiq.From("test").Where(depiq.Ex{
		"a": depiq.Op{"gt": 10},
		"b": depiq.Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).ToSQL()
	fmt.Println(sql)
	// You can use ExOr to get ORed expressions together
	sql, _, _ = depiq.From("test").Where(depiq.ExOr{
		"a": depiq.Op{"gt": 10},
		"b": depiq.Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).ToSQL()
	fmt.Println(sql)
	// You can use Or with Ex to Or multiple Ex maps together
	sql, _, _ = depiq.From("test").Where(
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
	sql, _, _ = depiq.From("test").Where(
		depiq.C("a").Gt(10),
		depiq.C("b").Lt(10),
		depiq.C("c").IsNull(),
		depiq.C("d").In("a", "b", "c"),
	).ToSQL()
	fmt.Println(sql)
	// You can use a combination of Ors and Ands
	sql, _, _ = depiq.From("test").Where(
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
	// SELECT * FROM "test" WHERE (("a" > 10) AND ("b" < 10) AND ("c" IS NULL) AND ("d" IN ('a', 'b', 'c')))
	// SELECT * FROM "test" WHERE (("a" > 10) OR ("b" < 10) OR ("c" IS NULL) OR ("d" IN ('a', 'b', 'c')))
	// SELECT * FROM "test" WHERE ((("a" > 10) AND ("b" < 10)) OR (("c" IS NULL) AND ("d" IN ('a', 'b', 'c'))))
	// SELECT * FROM "test" WHERE (("a" > 10) AND ("b" < 10) AND ("c" IS NULL) AND ("d" IN ('a', 'b', 'c')))
	// SELECT * FROM "test" WHERE (("a" > 10) OR (("b" < 10) AND ("c" IS NULL)))
}

func ExampleSelectDataset_Where_prepared() {
	// By default everything is anded together
	sql, args, _ := depiq.From("test").Prepared(true).Where(depiq.Ex{
		"a": depiq.Op{"gt": 10},
		"b": depiq.Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).ToSQL()
	fmt.Println(sql, args)
	// You can use ExOr to get ORed expressions together
	sql, args, _ = depiq.From("test").Prepared(true).Where(depiq.ExOr{
		"a": depiq.Op{"gt": 10},
		"b": depiq.Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).ToSQL()
	fmt.Println(sql, args)
	// You can use Or with Ex to Or multiple Ex maps together
	sql, args, _ = depiq.From("test").Prepared(true).Where(
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
	sql, args, _ = depiq.From("test").Prepared(true).Where(
		depiq.C("a").Gt(10),
		depiq.C("b").Lt(10),
		depiq.C("c").IsNull(),
		depiq.C("d").In("a", "b", "c"),
	).ToSQL()
	fmt.Println(sql, args)
	// You can use a combination of Ors and Ands
	sql, args, _ = depiq.From("test").Prepared(true).Where(
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
	// SELECT * FROM "test" WHERE (("a" > ?) AND ("b" < ?) AND ("c" IS NULL) AND ("d" IN (?, ?, ?))) [10 10 a b c]
	// SELECT * FROM "test" WHERE (("a" > ?) OR ("b" < ?) OR ("c" IS NULL) OR ("d" IN (?, ?, ?))) [10 10 a b c]
	// SELECT * FROM "test" WHERE ((("a" > ?) AND ("b" < ?)) OR (("c" IS NULL) AND ("d" IN (?, ?, ?)))) [10 10 a b c]
	// SELECT * FROM "test" WHERE (("a" > ?) AND ("b" < ?) AND ("c" IS NULL) AND ("d" IN (?, ?, ?))) [10 10 a b c]
	// SELECT * FROM "test" WHERE (("a" > ?) OR (("b" < ?) AND ("c" IS NULL))) [10 10]
}

func ExampleSelectDataset_ClearWhere() {
	ds := depiq.From("test").Where(
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
	// SELECT * FROM "test"
}

func ExampleSelectDataset_Join() {
	sql, _, _ := depiq.From("test").Join(
		depiq.T("test2"),
		depiq.On(depiq.Ex{"test.fkey": depiq.I("test2.Id")}),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Join(depiq.T("test2"), depiq.Using("common_column")).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Join(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)),
		depiq.On(depiq.I("test.fkey").Eq(depiq.T("test2").Col("Id"))),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Join(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)).As("t"),
		depiq.On(depiq.T("test").Col("fkey").Eq(depiq.T("t").Col("Id"))),
	).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" INNER JOIN "test2" ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" INNER JOIN "test2" USING ("common_column")
	// SELECT * FROM "test" INNER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" INNER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t" ON ("test"."fkey" = "t"."Id")
}

func ExampleSelectDataset_InnerJoin() {
	sql, _, _ := depiq.From("test").InnerJoin(
		depiq.T("test2"),
		depiq.On(depiq.Ex{
			"test.fkey": depiq.I("test2.Id"),
		}),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").InnerJoin(
		depiq.T("test2"),
		depiq.Using("common_column"),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").InnerJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)),
		depiq.On(depiq.I("test.fkey").Eq(depiq.I("test2.Id"))),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").InnerJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)).As("t"),
		depiq.On(depiq.I("test.fkey").Eq(depiq.I("t.Id"))),
	).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" INNER JOIN "test2" ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" INNER JOIN "test2" USING ("common_column")
	// SELECT * FROM "test" INNER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" INNER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t" ON ("test"."fkey" = "t"."Id")
}

func ExampleSelectDataset_FullOuterJoin() {
	sql, _, _ := depiq.From("test").FullOuterJoin(
		depiq.T("test2"),
		depiq.On(depiq.Ex{
			"test.fkey": depiq.I("test2.Id"),
		}),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").FullOuterJoin(
		depiq.T("test2"),
		depiq.Using("common_column"),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").FullOuterJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)),
		depiq.On(depiq.I("test.fkey").Eq(depiq.I("test2.Id"))),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").FullOuterJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)).As("t"),
		depiq.On(depiq.I("test.fkey").Eq(depiq.I("t.Id"))),
	).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" FULL OUTER JOIN "test2" ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" FULL OUTER JOIN "test2" USING ("common_column")
	// SELECT * FROM "test" FULL OUTER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" FULL OUTER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t" ON ("test"."fkey" = "t"."Id")
}

func ExampleSelectDataset_RightOuterJoin() {
	sql, _, _ := depiq.From("test").RightOuterJoin(
		depiq.T("test2"),
		depiq.On(depiq.Ex{
			"test.fkey": depiq.I("test2.Id"),
		}),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").RightOuterJoin(
		depiq.T("test2"),
		depiq.Using("common_column"),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").RightOuterJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)),
		depiq.On(depiq.I("test.fkey").Eq(depiq.I("test2.Id"))),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").RightOuterJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)).As("t"),
		depiq.On(depiq.I("test.fkey").Eq(depiq.I("t.Id"))),
	).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" RIGHT OUTER JOIN "test2" ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" RIGHT OUTER JOIN "test2" USING ("common_column")
	// SELECT * FROM "test" RIGHT OUTER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" RIGHT OUTER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t" ON ("test"."fkey" = "t"."Id")
}

func ExampleSelectDataset_LeftOuterJoin() {
	sql, _, _ := depiq.From("test").LeftOuterJoin(
		depiq.T("test2"),
		depiq.On(depiq.Ex{
			"test.fkey": depiq.I("test2.Id"),
		}),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").LeftOuterJoin(
		depiq.T("test2"),
		depiq.Using("common_column"),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").LeftOuterJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)),
		depiq.On(depiq.I("test.fkey").Eq(depiq.I("test2.Id"))),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").LeftOuterJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)).As("t"),
		depiq.On(depiq.I("test.fkey").Eq(depiq.I("t.Id"))),
	).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" LEFT OUTER JOIN "test2" ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" LEFT OUTER JOIN "test2" USING ("common_column")
	// SELECT * FROM "test" LEFT OUTER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" LEFT OUTER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t" ON ("test"."fkey" = "t"."Id")
}

func ExampleSelectDataset_FullJoin() {
	sql, _, _ := depiq.From("test").FullJoin(
		depiq.T("test2"),
		depiq.On(depiq.Ex{
			"test.fkey": depiq.I("test2.Id"),
		}),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").FullJoin(
		depiq.T("test2"),
		depiq.Using("common_column"),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").FullJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)),
		depiq.On(depiq.I("test.fkey").Eq(depiq.I("test2.Id"))),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").FullJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)).As("t"),
		depiq.On(depiq.I("test.fkey").Eq(depiq.I("t.Id"))),
	).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" FULL JOIN "test2" ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" FULL JOIN "test2" USING ("common_column")
	// SELECT * FROM "test" FULL JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" FULL JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t" ON ("test"."fkey" = "t"."Id")
}

func ExampleSelectDataset_RightJoin() {
	sql, _, _ := depiq.From("test").RightJoin(
		depiq.T("test2"),
		depiq.On(depiq.Ex{
			"test.fkey": depiq.I("test2.Id"),
		}),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").RightJoin(
		depiq.T("test2"),
		depiq.Using("common_column"),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").RightJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)),
		depiq.On(depiq.I("test.fkey").Eq(depiq.I("test2.Id"))),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").RightJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)).As("t"),
		depiq.On(depiq.I("test.fkey").Eq(depiq.I("t.Id"))),
	).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" RIGHT JOIN "test2" ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" RIGHT JOIN "test2" USING ("common_column")
	// SELECT * FROM "test" RIGHT JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" RIGHT JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t" ON ("test"."fkey" = "t"."Id")
}

func ExampleSelectDataset_LeftJoin() {
	sql, _, _ := depiq.From("test").LeftJoin(
		depiq.T("test2"),
		depiq.On(depiq.Ex{
			"test.fkey": depiq.I("test2.Id"),
		}),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").LeftJoin(
		depiq.T("test2"),
		depiq.Using("common_column"),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").LeftJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)),
		depiq.On(depiq.I("test.fkey").Eq(depiq.I("test2.Id"))),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").LeftJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)).As("t"),
		depiq.On(depiq.I("test.fkey").Eq(depiq.I("t.Id"))),
	).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" LEFT JOIN "test2" ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" LEFT JOIN "test2" USING ("common_column")
	// SELECT * FROM "test" LEFT JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" LEFT JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t" ON ("test"."fkey" = "t"."Id")
}

func ExampleSelectDataset_NaturalJoin() {
	sql, _, _ := depiq.From("test").NaturalJoin(depiq.T("test2")).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").NaturalJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").NaturalJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)).As("t"),
	).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" NATURAL JOIN "test2"
	// SELECT * FROM "test" NATURAL JOIN (SELECT * FROM "test2" WHERE ("amount" > 0))
	// SELECT * FROM "test" NATURAL JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t"
}

func ExampleSelectDataset_NaturalLeftJoin() {
	sql, _, _ := depiq.From("test").NaturalLeftJoin(depiq.T("test2")).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").NaturalLeftJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").NaturalLeftJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)).As("t"),
	).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" NATURAL LEFT JOIN "test2"
	// SELECT * FROM "test" NATURAL LEFT JOIN (SELECT * FROM "test2" WHERE ("amount" > 0))
	// SELECT * FROM "test" NATURAL LEFT JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t"
}

func ExampleSelectDataset_NaturalRightJoin() {
	sql, _, _ := depiq.From("test").NaturalRightJoin(depiq.T("test2")).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").NaturalRightJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").NaturalRightJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)).As("t"),
	).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" NATURAL RIGHT JOIN "test2"
	// SELECT * FROM "test" NATURAL RIGHT JOIN (SELECT * FROM "test2" WHERE ("amount" > 0))
	// SELECT * FROM "test" NATURAL RIGHT JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t"
}

func ExampleSelectDataset_NaturalFullJoin() {
	sql, _, _ := depiq.From("test").NaturalFullJoin(depiq.T("test2")).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").NaturalFullJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").NaturalFullJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)).As("t"),
	).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" NATURAL FULL JOIN "test2"
	// SELECT * FROM "test" NATURAL FULL JOIN (SELECT * FROM "test2" WHERE ("amount" > 0))
	// SELECT * FROM "test" NATURAL FULL JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t"
}

func ExampleSelectDataset_CrossJoin() {
	sql, _, _ := depiq.From("test").CrossJoin(depiq.T("test2")).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").CrossJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").CrossJoin(
		depiq.From("test2").Where(depiq.C("amount").Gt(0)).As("t"),
	).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" CROSS JOIN "test2"
	// SELECT * FROM "test" CROSS JOIN (SELECT * FROM "test2" WHERE ("amount" > 0))
	// SELECT * FROM "test" CROSS JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t"
}

func ExampleSelectDataset_FromSelf() {
	sql, _, _ := depiq.From("test").FromSelf().ToSQL()
	fmt.Println(sql)
	sql, _, _ = depiq.From("test").As("my_test_table").FromSelf().ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM (SELECT * FROM "test") AS "t1"
	// SELECT * FROM (SELECT * FROM "test") AS "my_test_table"
}

func ExampleSelectDataset_From() {
	ds := depiq.From("test")
	sql, _, _ := ds.From("test2").ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test2"
}

func ExampleSelectDataset_From_withDataset() {
	ds := depiq.From("test")
	fromDs := ds.Where(depiq.C("age").Gt(10))
	sql, _, _ := ds.From(fromDs).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM (SELECT * FROM "test" WHERE ("age" > 10)) AS "t1"
}

func ExampleSelectDataset_From_withAliasedDataset() {
	ds := depiq.From("test")
	fromDs := ds.Where(depiq.C("age").Gt(10))
	sql, _, _ := ds.From(fromDs.As("test2")).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM (SELECT * FROM "test" WHERE ("age" > 10)) AS "test2"
}

func ExampleSelectDataset_Select() {
	sql, _, _ := depiq.From("test").Select("a", "b", "c").ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT "a", "b", "c" FROM "test"
}

func ExampleSelectDataset_Select_withDataset() {
	ds := depiq.From("test")
	fromDs := ds.Select("age").Where(depiq.C("age").Gt(10))
	sql, _, _ := ds.From().Select(fromDs).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT (SELECT "age" FROM "test" WHERE ("age" > 10))
}

func ExampleSelectDataset_Select_withAliasedDataset() {
	ds := depiq.From("test")
	fromDs := ds.Select("age").Where(depiq.C("age").Gt(10))
	sql, _, _ := ds.From().Select(fromDs.As("ages")).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT (SELECT "age" FROM "test" WHERE ("age" > 10)) AS "ages"
}

func ExampleSelectDataset_Select_withLiteral() {
	sql, _, _ := depiq.From("test").Select(depiq.L("a + b").As("sum")).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT a + b AS "sum" FROM "test"
}

func ExampleSelectDataset_Select_withSQLFunctionExpression() {
	sql, _, _ := depiq.From("test").Select(
		depiq.COUNT("*").As("age_count"),
		depiq.MAX("age").As("max_age"),
		depiq.AVG("age").As("avg_age"),
	).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT COUNT(*) AS "age_count", MAX("age") AS "max_age", AVG("age") AS "avg_age" FROM "test"
}

func ExampleSelectDataset_Select_withStruct() {
	ds := depiq.From("test")

	type myStruct struct {
		Name         string
		Address      string `db:"address"`
		EmailAddress string `db:"email_address"`
	}

	// Pass with pointer
	sql, _, _ := ds.Select(&myStruct{}).ToSQL()
	fmt.Println(sql)

	// Pass instance of
	sql, _, _ = ds.Select(myStruct{}).ToSQL()
	fmt.Println(sql)

	type myStruct2 struct {
		myStruct
		Zipcode string `db:"zipcode"`
	}

	// Pass pointer to struct with embedded struct
	sql, _, _ = ds.Select(&myStruct2{}).ToSQL()
	fmt.Println(sql)

	// Pass instance of struct with embedded struct
	sql, _, _ = ds.Select(myStruct2{}).ToSQL()
	fmt.Println(sql)

	var myStructs []myStruct

	// Pass slice of structs, will only select columns from underlying type
	sql, _, _ = ds.Select(myStructs).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT "address", "email_address", "name" FROM "test"
	// SELECT "address", "email_address", "name" FROM "test"
	// SELECT "address", "email_address", "name", "zipcode" FROM "test"
	// SELECT "address", "email_address", "name", "zipcode" FROM "test"
	// SELECT "address", "email_address", "name" FROM "test"
}

func ExampleSelectDataset_Distinct() {
	sql, _, _ := depiq.From("test").Select("a", "b").Distinct().ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT DISTINCT "a", "b" FROM "test"
}

func ExampleSelectDataset_Distinct_on() {
	sql, _, _ := depiq.From("test").Distinct("a").ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT DISTINCT ON ("a") * FROM "test"
}

func ExampleSelectDataset_Distinct_onWithLiteral() {
	sql, _, _ := depiq.From("test").Distinct(depiq.L("COALESCE(?, ?)", depiq.C("a"), "empty")).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT DISTINCT ON (COALESCE("a", 'empty')) * FROM "test"
}

func ExampleSelectDataset_Distinct_onCoalesce() {
	sql, _, _ := depiq.From("test").Distinct(depiq.COALESCE(depiq.C("a"), "empty")).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT DISTINCT ON (COALESCE("a", 'empty')) * FROM "test"
}

func ExampleSelectDataset_SelectAppend() {
	ds := depiq.From("test").Select("a", "b")
	sql, _, _ := ds.SelectAppend("c").ToSQL()
	fmt.Println(sql)
	ds = depiq.From("test").Select("a", "b").Distinct()
	sql, _, _ = ds.SelectAppend("c").ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT "a", "b", "c" FROM "test"
	// SELECT DISTINCT "a", "b", "c" FROM "test"
}

func ExampleSelectDataset_ClearSelect() {
	ds := depiq.From("test").Select("a", "b")
	sql, _, _ := ds.ClearSelect().ToSQL()
	fmt.Println(sql)
	ds = depiq.From("test").Select("a", "b").Distinct()
	sql, _, _ = ds.ClearSelect().ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test"
	// SELECT * FROM "test"
}

func ExampleSelectDataset_ToSQL() {
	sql, args, _ := depiq.From("items").Where(depiq.Ex{"a": 1}).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "items" WHERE ("a" = 1) []
}

func ExampleSelectDataset_ToSQL_prepared() {
	sql, args, _ := depiq.From("items").Where(depiq.Ex{"a": 1}).Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "items" WHERE ("a" = ?) [1]
}

func ExampleSelectDataset_Update() {
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name"`
	}
	sql, args, _ := depiq.From("items").Update().Set(
		item{Name: "Test", Address: "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("items").Update().Set(
		depiq.Record{"name": "Test", "address": "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("items").Update().Set(
		map[string]interface{}{"name": "Test", "address": "111 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
}

func ExampleSelectDataset_Insert() {
	type item struct {
		ID      uint32 `db:"id" goqu:"skipinsert"`
		Address string `db:"address"`
		Name    string `db:"name"`
	}
	sql, args, _ := depiq.From("items").Insert().Rows(
		item{Name: "Test1", Address: "111 Test Addr"},
		item{Name: "Test2", Address: "112 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("items").Insert().Rows(
		depiq.Record{"name": "Test1", "address": "111 Test Addr"},
		depiq.Record{"name": "Test2", "address": "112 Test Addr"},
	).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("items").Insert().Rows(
		[]item{
			{Name: "Test1", Address: "111 Test Addr"},
			{Name: "Test2", Address: "112 Test Addr"},
		}).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("items").Insert().Rows(
		[]depiq.Record{
			{"name": "Test1", "address": "111 Test Addr"},
			{"name": "Test2", "address": "112 Test Addr"},
		}).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// INSERT INTO "items" ("address", "name") VALUES ('111 Test Addr', 'Test1'), ('112 Test Addr', 'Test2') []
	// INSERT INTO "items" ("address", "name") VALUES ('111 Test Addr', 'Test1'), ('112 Test Addr', 'Test2') []
	// INSERT INTO "items" ("address", "name") VALUES ('111 Test Addr', 'Test1'), ('112 Test Addr', 'Test2') []
	// INSERT INTO "items" ("address", "name") VALUES ('111 Test Addr', 'Test1'), ('112 Test Addr', 'Test2') []
}

func ExampleSelectDataset_Delete() {
	sql, args, _ := depiq.From("items").Delete().ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("items").
		Where(depiq.Ex{"id": depiq.Op{"gt": 10}}).
		Delete().
		ToSQL()
	fmt.Println(sql, args)

	// Output:
	// DELETE FROM "items" []
	// DELETE FROM "items" WHERE ("id" > 10) []
}

func ExampleSelectDataset_Truncate() {
	sql, args, _ := depiq.From("items").Truncate().ToSQL()
	fmt.Println(sql, args)
	// Output:
	// TRUNCATE "items" []
}

func ExampleSelectDataset_Prepared() {
	sql, args, _ := depiq.From("items").Prepared(true).Where(depiq.Ex{
		"col1": "a",
		"col2": 1,
		"col3": true,
		"col4": false,
		"col5": []string{"a", "b", "c"},
	}).ToSQL()
	fmt.Println(sql, args)
	// nolint:lll // sql statements are long
	// Output:
	// SELECT * FROM "items" WHERE (("col1" = ?) AND ("col2" = ?) AND ("col3" IS TRUE) AND ("col4" IS FALSE) AND ("col5" IN (?, ?, ?))) [a 1 a b c]
}

func ExampleSelectDataset_ScanStructs() {
	type User struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
	}
	db := getDB()
	var users []User
	if err := db.From("depiq_user").Fetch(&users); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("\n%+v", users)

	users = users[0:0]
	if err := db.From("depiq_user").Select("first_name").Fetch(&users); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("\n%+v", users)

	// Output:
	// [{FirstName:Bob LastName:Yukon} {FirstName:Sally LastName:Yukon} {FirstName:Vinita LastName:Yukon} {FirstName:John LastName:Doe}]
	// [{FirstName:Bob LastName:} {FirstName:Sally LastName:} {FirstName:Vinita LastName:} {FirstName:John LastName:}]
}

func ExampleSelectDataset_ScanStructs_prepared() {
	type User struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
	}
	db := getDB()

	ds := db.From("depiq_user").
		Prepared(true).
		Where(depiq.Ex{
			"last_name": "Yukon",
		})

	var users []User
	if err := ds.Fetch(&users); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("\n%+v", users)

	// Output:
	// [{FirstName:Bob LastName:Yukon} {FirstName:Sally LastName:Yukon} {FirstName:Vinita LastName:Yukon}]
}

// In this example we create a new struct that has two structs that represent two table
// the User and Role fields are tagged with the table name
func ExampleSelectDataset_ScanStructs_withJoinAutoSelect() {
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
		Role Role `db:"user_role"`  // tag as "user_role" table
	}
	db := getDB()

	ds := db.
		From("depiq_user").
		Join(depiq.T("user_role"), depiq.On(depiq.I("depiq_user.id").Eq(depiq.I("user_role.user_id"))))
	var users []UserAndRole
	// Scan structs will auto build the
	if err := ds.Fetch(&users); err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, u := range users {
		fmt.Printf("\n%+v", u)
	}
	// Output:
	// {User:{ID:1 FirstName:Bob LastName:Yukon} Role:{UserID:1 Name:Admin}}
	// {User:{ID:2 FirstName:Sally LastName:Yukon} Role:{UserID:2 Name:Manager}}
	// {User:{ID:3 FirstName:Vinita LastName:Yukon} Role:{UserID:3 Name:Manager}}
	// {User:{ID:4 FirstName:John LastName:Doe} Role:{UserID:4 Name:User}}
}

// In this example we create a new struct that has the user properties as well as a nested
// Role struct from the join table
func ExampleSelectDataset_ScanStructs_withJoinManualSelect() {
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
	db := getDB()

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
	if err := ds.Fetch(&users); err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, u := range users {
		fmt.Printf("\n%+v", u)
	}

	// Output:
	// {ID:1 FirstName:Bob LastName:Yukon Role:{UserID:1 Name:Admin}}
	// {ID:2 FirstName:Sally LastName:Yukon Role:{UserID:2 Name:Manager}}
	// {ID:3 FirstName:Vinita LastName:Yukon Role:{UserID:3 Name:Manager}}
	// {ID:4 FirstName:John LastName:Doe Role:{UserID:4 Name:User}}
}

func ExampleSelectDataset_ScanStruct() {
	type User struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
	}
	db := getDB()
	findUserByName := func(name string) {
		var user User
		ds := db.From("depiq_user").Where(depiq.C("first_name").Eq(name))
		found, err := ds.FecthRow(&user)
		switch {
		case err != nil:
			fmt.Println(err.Error())
		case !found:
			fmt.Printf("No user found for first_name %s\n", name)
		default:
			fmt.Printf("Found user: %+v\n", user)
		}
	}

	findUserByName("Bob")
	findUserByName("Zeb")

	// Output:
	// Found user: {FirstName:Bob LastName:Yukon}
	// No user found for first_name Zeb
}

// In this example we create a new struct that has two structs that represent two table
// the User and Role fields are tagged with the table name
func ExampleSelectDataset_ScanStruct_withJoinAutoSelect() {
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
		Role Role `db:"user_role"`  // tag as "user_role" table
	}
	db := getDB()
	findUserAndRoleByName := func(name string) {
		var userAndRole UserAndRole
		ds := db.
			From("depiq_user").
			Join(
				depiq.T("user_role"),
				depiq.On(depiq.I("depiq_user.id").Eq(depiq.I("user_role.user_id"))),
			).
			Where(depiq.C("first_name").Eq(name))
		found, err := ds.FecthRow(&userAndRole)
		switch {
		case err != nil:
			fmt.Println(err.Error())
		case !found:
			fmt.Printf("No user found for first_name %s\n", name)
		default:
			fmt.Printf("Found user and role: %+v\n", userAndRole)
		}
	}

	findUserAndRoleByName("Bob")
	findUserAndRoleByName("Zeb")
	// Output:
	// Found user and role: {User:{ID:1 FirstName:Bob LastName:Yukon} Role:{UserID:1 Name:Admin}}
	// No user found for first_name Zeb
}

// In this example we create a new struct that has the user properties as well as a nested
// Role struct from the join table
func ExampleSelectDataset_ScanStruct_withJoinManualSelect() {
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
	db := getDB()
	findUserByName := func(name string) {
		var userAndRole User
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
			Join(
				depiq.T("user_role"),
				depiq.On(depiq.I("depiq_user.id").Eq(depiq.I("user_role.user_id"))),
			).
			Where(depiq.C("first_name").Eq(name))
		found, err := ds.FecthRow(&userAndRole)
		switch {
		case err != nil:
			fmt.Println(err.Error())
		case !found:
			fmt.Printf("No user found for first_name %s\n", name)
		default:
			fmt.Printf("Found user and role: %+v\n", userAndRole)
		}
	}

	findUserByName("Bob")
	findUserByName("Zeb")

	// Output:
	// Found user and role: {ID:1 FirstName:Bob LastName:Yukon Role:{UserID:1 Name:Admin}}
	// No user found for first_name Zeb
}

func ExampleSelectDataset_ScanVals() {
	var ids []int64
	if err := getDB().From("depiq_user").Select("id").ScanVals(&ids); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("UserIds = %+v", ids)

	// Output:
	// UserIds = [1 2 3 4]
}

func ExampleSelectDataset_ScanVal() {
	db := getDB()
	findUserIDByName := func(name string) {
		var id int64
		ds := db.From("depiq_user").
			Select("id").
			Where(depiq.C("first_name").Eq(name))

		found, err := ds.ScanVal(&id)
		switch {
		case err != nil:
			fmt.Println(err.Error())
		case !found:
			fmt.Printf("No id found for user %s", name)
		default:
			fmt.Printf("\nFound userId: %+v\n", id)
		}
	}

	findUserIDByName("Bob")
	findUserIDByName("Zeb")
	// Output:
	// Found userId: 1
	// No id found for user Zeb
}

func ExampleSelectDataset_Count() {
	count, err := getDB().From("depiq_user").Count()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Count is %d", count)

	// Output:
	// Count is 4
}

func ExampleSelectDataset_Pluck() {
	var lastNames []string
	if err := getDB().From("depiq_user").Pluck(&lastNames, "last_name"); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("LastNames = %+v", lastNames)

	// Output:
	// LastNames = [Yukon Yukon Yukon Doe]
}

func ExampleSelectDataset_Executor_scannerScanStruct() {
	type User struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
	}
	db := getDB()

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

	// Output:
	// {FirstName:Bob LastName:Yukon}
	// {FirstName:Sally LastName:Yukon}
	// {FirstName:Vinita LastName:Yukon}
}

func ExampleSelectDataset_Executor_scannerScanVal() {
	db := getDB()

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

	// Output:
	// Bob
	// Sally
	// Vinita
}

func ExampleForUpdate() {
	sql, args, _ := depiq.From("test").ForUpdate(exp.Wait).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" FOR UPDATE  []
}

func ExampleForUpdate_of() {
	sql, args, _ := depiq.From("test").ForUpdate(exp.Wait, depiq.T("test")).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" FOR UPDATE OF "test"  []
}

func ExampleForUpdate_ofMultiple() {
	sql, args, _ := depiq.From("table1").Join(
		depiq.T("table2"),
		depiq.On(depiq.I("table2.id").Eq(depiq.I("table1.id"))),
	).ForUpdate(
		exp.Wait,
		depiq.T("table1"),
		depiq.T("table2"),
	).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "table1" INNER JOIN "table2" ON ("table2"."id" = "table1"."id") FOR UPDATE OF "table1", "table2"  []
}
