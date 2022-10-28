// nolint:lll // sql statements are long
package depiq_test

import (
	"fmt"
	"regexp"

	"github.com/orn-id/depiq/v9"
	"github.com/orn-id/depiq/v9/exp"
)

func ExampleAVG() {
	ds := depiq.From("test").Select(depiq.AVG("col"))
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT AVG("col") FROM "test" []
	// SELECT AVG("col") FROM "test" []
}

func ExampleAVG_as() {
	sql, _, _ := depiq.From("test").Select(depiq.AVG("a").As("a")).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT AVG("a") AS "a" FROM "test"
}

func ExampleAVG_havingClause() {
	ds := depiq.
		From("test").
		Select(depiq.AVG("a").As("avg")).
		GroupBy("a").
		Having(depiq.AVG("a").Gt(10))

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT AVG("a") AS "avg" FROM "test" GROUP BY "a" HAVING (AVG("a") > 10) []
	// SELECT AVG("a") AS "avg" FROM "test" GROUP BY "a" HAVING (AVG("a") > ?) [10]
}

func ExampleAnd() {
	ds := depiq.From("test").Where(
		depiq.And(
			depiq.C("col").Gt(10),
			depiq.C("col").Lt(20),
		),
	)
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE (("col" > 10) AND ("col" < 20)) []
	// SELECT * FROM "test" WHERE (("col" > ?) AND ("col" < ?)) [10 20]
}

// You can use And with Or to create more complex queries
func ExampleAnd_withOr() {
	ds := depiq.From("test").Where(
		depiq.And(
			depiq.C("col1").IsTrue(),
			depiq.Or(
				depiq.C("col2").Gt(10),
				depiq.C("col2").Lt(20),
			),
		),
	)
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// by default expressions are anded together
	ds = depiq.From("test").Where(
		depiq.C("col1").IsTrue(),
		depiq.Or(
			depiq.C("col2").Gt(10),
			depiq.C("col2").Lt(20),
		),
	)
	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE (("col1" IS TRUE) AND (("col2" > 10) OR ("col2" < 20))) []
	// SELECT * FROM "test" WHERE (("col1" IS TRUE) AND (("col2" > ?) OR ("col2" < ?))) [10 20]
	// SELECT * FROM "test" WHERE (("col1" IS TRUE) AND (("col2" > 10) OR ("col2" < 20))) []
	// SELECT * FROM "test" WHERE (("col1" IS TRUE) AND (("col2" > ?) OR ("col2" < ?))) [10 20]
}

// You can use ExOr inside of And expression lists.
func ExampleAnd_withExOr() {
	// by default expressions are anded together
	ds := depiq.From("test").Where(
		depiq.C("col1").IsTrue(),
		depiq.ExOr{
			"col2": depiq.Op{"gt": 10},
			"col3": depiq.Op{"lt": 20},
		},
	)
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE (("col1" IS TRUE) AND (("col2" > 10) OR ("col3" < 20))) []
	// SELECT * FROM "test" WHERE (("col1" IS TRUE) AND (("col2" > ?) OR ("col3" < ?))) [10 20]
}

func ExampleC() {
	sql, args, _ := depiq.From("test").
		Select(depiq.C("*")).
		ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").
		Select(depiq.C("col1")).
		ToSQL()
	fmt.Println(sql, args)

	ds := depiq.From("test").Where(
		depiq.C("col1").Eq(10),
		depiq.C("col2").In([]int64{1, 2, 3, 4}),
		depiq.C("col3").Like(regexp.MustCompile("^[ab]")),
		depiq.C("col4").IsNull(),
	)

	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" []
	// SELECT "col1" FROM "test" []
	// SELECT * FROM "test" WHERE (("col1" = 10) AND ("col2" IN (1, 2, 3, 4)) AND ("col3" ~ '^[ab]') AND ("col4" IS NULL)) []
	// SELECT * FROM "test" WHERE (("col1" = ?) AND ("col2" IN (?, ?, ?, ?)) AND ("col3" ~ ?) AND ("col4" IS NULL)) [10 1 2 3 4 ^[ab]]
}

func ExampleC_as() {
	sql, _, _ := depiq.From("test").Select(depiq.C("a").As("as_a")).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Select(depiq.C("a").As(depiq.C("as_a"))).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT "a" AS "as_a" FROM "test"
	// SELECT "a" AS "as_a" FROM "test"
}

func ExampleC_ordering() {
	sql, args, _ := depiq.From("test").Order(depiq.C("a").Asc()).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Order(depiq.C("a").Asc().NullsFirst()).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Order(depiq.C("a").Asc().NullsLast()).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Order(depiq.C("a").Desc()).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Order(depiq.C("a").Desc().NullsFirst()).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Order(depiq.C("a").Desc().NullsLast()).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" ORDER BY "a" ASC []
	// SELECT * FROM "test" ORDER BY "a" ASC NULLS FIRST []
	// SELECT * FROM "test" ORDER BY "a" ASC NULLS LAST []
	// SELECT * FROM "test" ORDER BY "a" DESC []
	// SELECT * FROM "test" ORDER BY "a" DESC NULLS FIRST []
	// SELECT * FROM "test" ORDER BY "a" DESC NULLS LAST []
}

func ExampleC_cast() {
	sql, _, _ := depiq.From("test").
		Select(depiq.C("json1").Cast("TEXT").As("json_text")).
		ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(
		depiq.C("json1").Cast("TEXT").Neq(
			depiq.C("json2").Cast("TEXT"),
		),
	).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT CAST("json1" AS TEXT) AS "json_text" FROM "test"
	// SELECT * FROM "test" WHERE (CAST("json1" AS TEXT) != CAST("json2" AS TEXT))
}

func ExampleC_comparisons() {
	// used from an identifier
	sql, _, _ := depiq.From("test").Where(depiq.C("a").Eq(10)).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.C("a").Neq(10)).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.C("a").Gt(10)).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.C("a").Gte(10)).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.C("a").Lt(10)).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.C("a").Lte(10)).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT * FROM "test" WHERE ("a" = 10)
	// SELECT * FROM "test" WHERE ("a" != 10)
	// SELECT * FROM "test" WHERE ("a" > 10)
	// SELECT * FROM "test" WHERE ("a" >= 10)
	// SELECT * FROM "test" WHERE ("a" < 10)
	// SELECT * FROM "test" WHERE ("a" <= 10)
}

func ExampleC_inOperators() {
	// using identifiers
	sql, _, _ := depiq.From("test").Where(depiq.C("a").In("a", "b", "c")).ToSQL()
	fmt.Println(sql)
	// with a slice
	sql, _, _ = depiq.From("test").Where(depiq.C("a").In([]string{"a", "b", "c"})).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.C("a").NotIn("a", "b", "c")).ToSQL()
	fmt.Println(sql)
	// with a slice
	sql, _, _ = depiq.From("test").Where(depiq.C("a").NotIn([]string{"a", "b", "c"})).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT * FROM "test" WHERE ("a" IN ('a', 'b', 'c'))
	// SELECT * FROM "test" WHERE ("a" IN ('a', 'b', 'c'))
	// SELECT * FROM "test" WHERE ("a" NOT IN ('a', 'b', 'c'))
	// SELECT * FROM "test" WHERE ("a" NOT IN ('a', 'b', 'c'))
}

func ExampleC_likeComparisons() {
	// using identifiers
	sql, _, _ := depiq.From("test").Where(depiq.C("a").Like("%a%")).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.C("a").Like(regexp.MustCompile(`[ab]`))).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.C("a").ILike("%a%")).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.C("a").ILike(regexp.MustCompile("[ab]"))).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.C("a").NotLike("%a%")).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.C("a").NotLike(regexp.MustCompile("[ab]"))).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.C("a").NotILike("%a%")).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.C("a").NotILike(regexp.MustCompile(`[ab]`))).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT * FROM "test" WHERE ("a" LIKE '%a%')
	// SELECT * FROM "test" WHERE ("a" ~ '[ab]')
	// SELECT * FROM "test" WHERE ("a" ILIKE '%a%')
	// SELECT * FROM "test" WHERE ("a" ~* '[ab]')
	// SELECT * FROM "test" WHERE ("a" NOT LIKE '%a%')
	// SELECT * FROM "test" WHERE ("a" !~ '[ab]')
	// SELECT * FROM "test" WHERE ("a" NOT ILIKE '%a%')
	// SELECT * FROM "test" WHERE ("a" !~* '[ab]')
}

func ExampleC_isComparisons() {
	sql, args, _ := depiq.From("test").Where(depiq.C("a").Is(nil)).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.C("a").Is(true)).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.C("a").Is(false)).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.C("a").IsNull()).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.C("a").IsTrue()).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.C("a").IsFalse()).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.C("a").IsNot(nil)).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.C("a").IsNot(true)).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.C("a").IsNot(false)).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.C("a").IsNotNull()).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.C("a").IsNotTrue()).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.C("a").IsNotFalse()).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ("a" IS NULL) []
	// SELECT * FROM "test" WHERE ("a" IS TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS FALSE) []
	// SELECT * FROM "test" WHERE ("a" IS NULL) []
	// SELECT * FROM "test" WHERE ("a" IS TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS FALSE) []
	// SELECT * FROM "test" WHERE ("a" IS NOT NULL) []
	// SELECT * FROM "test" WHERE ("a" IS NOT TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS NOT FALSE) []
	// SELECT * FROM "test" WHERE ("a" IS NOT NULL) []
	// SELECT * FROM "test" WHERE ("a" IS NOT TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS NOT FALSE) []
}

func ExampleC_betweenComparisons() {
	ds := depiq.From("test").Where(
		depiq.C("a").Between(depiq.Range(1, 10)),
	)
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(
		depiq.C("a").NotBetween(depiq.Range(1, 10)),
	)
	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ("a" BETWEEN 1 AND 10) []
	// SELECT * FROM "test" WHERE ("a" BETWEEN ? AND ?) [1 10]
	// SELECT * FROM "test" WHERE ("a" NOT BETWEEN 1 AND 10) []
	// SELECT * FROM "test" WHERE ("a" NOT BETWEEN ? AND ?) [1 10]
}

func ExampleCOALESCE() {
	ds := depiq.From("test").Select(
		depiq.COALESCE(depiq.C("a"), "a"),
		depiq.COALESCE(depiq.C("a"), depiq.C("b"), nil),
	)
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT COALESCE("a", 'a'), COALESCE("a", "b", NULL) FROM "test" []
	// SELECT COALESCE("a", ?), COALESCE("a", "b", ?) FROM "test" [a <nil>]
}

func ExampleCOALESCE_as() {
	sql, _, _ := depiq.From("test").Select(depiq.COALESCE(depiq.C("a"), "a").As("a")).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT COALESCE("a", 'a') AS "a" FROM "test"
}

func ExampleCOUNT() {
	ds := depiq.From("test").Select(depiq.COUNT("*"))
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT COUNT(*) FROM "test" []
	// SELECT COUNT(*) FROM "test" []
}

func ExampleCOUNT_as() {
	sql, _, _ := depiq.From("test").Select(depiq.COUNT("*").As("count")).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT COUNT(*) AS "count" FROM "test"
}

func ExampleCOUNT_havingClause() {
	ds := depiq.
		From("test").
		Select(depiq.COUNT("a").As("COUNT")).
		GroupBy("a").
		Having(depiq.COUNT("a").Gt(10))

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT COUNT("a") AS "COUNT" FROM "test" GROUP BY "a" HAVING (COUNT("a") > 10) []
	// SELECT COUNT("a") AS "COUNT" FROM "test" GROUP BY "a" HAVING (COUNT("a") > ?) [10]
}

func ExampleCast() {
	sql, _, _ := depiq.From("test").
		Select(depiq.Cast(depiq.C("json1"), "TEXT").As("json_text")).
		ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(
		depiq.Cast(depiq.C("json1"), "TEXT").Neq(
			depiq.Cast(depiq.C("json2"), "TEXT"),
		),
	).ToSQL()
	fmt.Println(sql)
	// Output:
	// SELECT CAST("json1" AS TEXT) AS "json_text" FROM "test"
	// SELECT * FROM "test" WHERE (CAST("json1" AS TEXT) != CAST("json2" AS TEXT))
}

func ExampleDISTINCT() {
	ds := depiq.From("test").Select(depiq.DISTINCT("col"))
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT DISTINCT("col") FROM "test" []
	// SELECT DISTINCT("col") FROM "test" []
}

func ExampleDISTINCT_as() {
	sql, _, _ := depiq.From("test").Select(depiq.DISTINCT("a").As("distinct_a")).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT DISTINCT("a") AS "distinct_a" FROM "test"
}

func ExampleDefault() {
	ds := depiq.Insert("items")

	sql, args, _ := ds.Rows(depiq.Record{
		"name":    depiq.Default(),
		"address": depiq.Default(),
	}).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Rows(depiq.Record{
		"name":    depiq.Default(),
		"address": depiq.Default(),
	}).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// INSERT INTO "items" ("address", "name") VALUES (DEFAULT, DEFAULT) []
	// INSERT INTO "items" ("address", "name") VALUES (DEFAULT, DEFAULT) []
}

func ExampleDoNothing() {
	ds := depiq.Insert("items")

	sql, args, _ := ds.Rows(depiq.Record{
		"address": "111 Address",
		"name":    "bob",
	}).OnConflict(depiq.DoNothing()).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Rows(depiq.Record{
		"address": "111 Address",
		"name":    "bob",
	}).OnConflict(depiq.DoNothing()).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// INSERT INTO "items" ("address", "name") VALUES ('111 Address', 'bob') ON CONFLICT DO NOTHING []
	// INSERT INTO "items" ("address", "name") VALUES (?, ?) ON CONFLICT DO NOTHING [111 Address bob]
}

func ExampleDoUpdate() {
	ds := depiq.Insert("items")

	sql, args, _ := ds.
		Rows(depiq.Record{"address": "111 Address"}).
		OnConflict(depiq.DoUpdate("address", depiq.C("address").Set(depiq.I("excluded.address")))).
		ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).
		Rows(depiq.Record{"address": "111 Address"}).
		OnConflict(depiq.DoUpdate("address", depiq.C("address").Set(depiq.I("excluded.address")))).
		ToSQL()
	fmt.Println(sql, args)

	// Output:
	// INSERT INTO "items" ("address") VALUES ('111 Address') ON CONFLICT (address) DO UPDATE SET "address"="excluded"."address" []
	// INSERT INTO "items" ("address") VALUES (?) ON CONFLICT (address) DO UPDATE SET "address"="excluded"."address" [111 Address]
}

func ExampleDoUpdate_where() {
	ds := depiq.Insert("items")

	sql, args, _ := ds.
		Rows(depiq.Record{"address": "111 Address"}).
		OnConflict(depiq.DoUpdate(
			"address",
			depiq.C("address").Set(depiq.I("excluded.address"))).Where(depiq.I("items.updated").IsNull()),
		).
		ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).
		Rows(depiq.Record{"address": "111 Address"}).
		OnConflict(depiq.DoUpdate(
			"address",
			depiq.C("address").Set(depiq.I("excluded.address"))).Where(depiq.I("items.updated").IsNull()),
		).
		ToSQL()
	fmt.Println(sql, args)

	// Output:
	// INSERT INTO "items" ("address") VALUES ('111 Address') ON CONFLICT (address) DO UPDATE SET "address"="excluded"."address" WHERE ("items"."updated" IS NULL) []
	// INSERT INTO "items" ("address") VALUES (?) ON CONFLICT (address) DO UPDATE SET "address"="excluded"."address" WHERE ("items"."updated" IS NULL) [111 Address]
}

func ExampleFIRST() {
	ds := depiq.From("test").Select(depiq.FIRST("col"))
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT FIRST("col") FROM "test" []
	// SELECT FIRST("col") FROM "test" []
}

func ExampleFIRST_as() {
	sql, _, _ := depiq.From("test").Select(depiq.FIRST("a").As("a")).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT FIRST("a") AS "a" FROM "test"
}

// This example shows how to create custom SQL Functions
func ExampleFunc() {
	stragg := func(expression exp.Expression, delimiter string) exp.SQLFunctionExpression {
		return depiq.Func("str_agg", expression, depiq.L(delimiter))
	}
	sql, _, _ := depiq.From("test").Select(stragg(depiq.C("col"), "|")).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT str_agg("col", |) FROM "test"
}

func ExampleI() {
	ds := depiq.From("test").
		Select(
			depiq.I("my_schema.table.col1"),
			depiq.I("table.col2"),
			depiq.I("col3"),
		)

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Select(depiq.I("test.*"))

	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT "my_schema"."table"."col1", "table"."col2", "col3" FROM "test" []
	// SELECT "my_schema"."table"."col1", "table"."col2", "col3" FROM "test" []
	// SELECT "test".* FROM "test" []
	// SELECT "test".* FROM "test" []
}

func ExampleL() {
	ds := depiq.From("test").Where(
		// literal with no args
		depiq.L(`"col"::TEXT = ""other_col"::text`),
		// literal with args they will be interpolated into the sql by default
		depiq.L("col IN (?, ?, ?)", "a", "b", "c"),
	)

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "test" WHERE ("col"::TEXT = ""other_col"::text AND col IN ('a', 'b', 'c')) []
	// SELECT * FROM "test" WHERE ("col"::TEXT = ""other_col"::text AND col IN (?, ?, ?)) [a b c]
}

func ExampleL_withArgs() {
	ds := depiq.From("test").Where(
		depiq.L(
			"(? AND ?) OR ?",
			depiq.C("a").Eq(1),
			depiq.C("b").Eq("b"),
			depiq.C("c").In([]string{"a", "b", "c"}),
		),
	)

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "test" WHERE (("a" = 1) AND ("b" = 'b')) OR ("c" IN ('a', 'b', 'c')) []
	// SELECT * FROM "test" WHERE (("a" = ?) AND ("b" = ?)) OR ("c" IN (?, ?, ?)) [1 b a b c]
}

func ExampleL_as() {
	sql, _, _ := depiq.From("test").Select(depiq.L("json_col->>'totalAmount'").As("total_amount")).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT json_col->>'totalAmount' AS "total_amount" FROM "test"
}

func ExampleL_comparisons() {
	// used from a literal expression
	sql, _, _ := depiq.From("test").Where(depiq.L("(a + b)").Eq(10)).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.L("(a + b)").Neq(10)).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.L("(a + b)").Gt(10)).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.L("(a + b)").Gte(10)).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.L("(a + b)").Lt(10)).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.L("(a + b)").Lte(10)).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT * FROM "test" WHERE ((a + b) = 10)
	// SELECT * FROM "test" WHERE ((a + b) != 10)
	// SELECT * FROM "test" WHERE ((a + b) > 10)
	// SELECT * FROM "test" WHERE ((a + b) >= 10)
	// SELECT * FROM "test" WHERE ((a + b) < 10)
	// SELECT * FROM "test" WHERE ((a + b) <= 10)
}

func ExampleL_inOperators() {
	// using identifiers
	sql, _, _ := depiq.From("test").Where(depiq.L("json_col->>'val'").In("a", "b", "c")).ToSQL()
	fmt.Println(sql)
	// with a slice
	sql, _, _ = depiq.From("test").Where(depiq.L("json_col->>'val'").In([]string{"a", "b", "c"})).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.L("json_col->>'val'").NotIn("a", "b", "c")).ToSQL()
	fmt.Println(sql)
	// with a slice
	sql, _, _ = depiq.From("test").Where(depiq.L("json_col->>'val'").NotIn([]string{"a", "b", "c"})).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT * FROM "test" WHERE (json_col->>'val' IN ('a', 'b', 'c'))
	// SELECT * FROM "test" WHERE (json_col->>'val' IN ('a', 'b', 'c'))
	// SELECT * FROM "test" WHERE (json_col->>'val' NOT IN ('a', 'b', 'c'))
	// SELECT * FROM "test" WHERE (json_col->>'val' NOT IN ('a', 'b', 'c'))
}

func ExampleL_likeComparisons() {
	// using identifiers
	sql, _, _ := depiq.From("test").Where(depiq.L("(a::text || 'bar')").Like("%a%")).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(
		depiq.L("(a::text || 'bar')").Like(regexp.MustCompile("[ab]")),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.L("(a::text || 'bar')").ILike("%a%")).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(
		depiq.L("(a::text || 'bar')").ILike(regexp.MustCompile("[ab]")),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.L("(a::text || 'bar')").NotLike("%a%")).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(
		depiq.L("(a::text || 'bar')").NotLike(regexp.MustCompile("[ab]")),
	).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(depiq.L("(a::text || 'bar')").NotILike("%a%")).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("test").Where(
		depiq.L("(a::text || 'bar')").NotILike(regexp.MustCompile("[ab]")),
	).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT * FROM "test" WHERE ((a::text || 'bar') LIKE '%a%')
	// SELECT * FROM "test" WHERE ((a::text || 'bar') ~ '[ab]')
	// SELECT * FROM "test" WHERE ((a::text || 'bar') ILIKE '%a%')
	// SELECT * FROM "test" WHERE ((a::text || 'bar') ~* '[ab]')
	// SELECT * FROM "test" WHERE ((a::text || 'bar') NOT LIKE '%a%')
	// SELECT * FROM "test" WHERE ((a::text || 'bar') !~ '[ab]')
	// SELECT * FROM "test" WHERE ((a::text || 'bar') NOT ILIKE '%a%')
	// SELECT * FROM "test" WHERE ((a::text || 'bar') !~* '[ab]')
}

func ExampleL_isComparisons() {
	sql, args, _ := depiq.From("test").Where(depiq.L("a").Is(nil)).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.L("a").Is(true)).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.L("a").Is(false)).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.L("a").IsNull()).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.L("a").IsTrue()).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.L("a").IsFalse()).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.L("a").IsNot(nil)).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.L("a").IsNot(true)).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.L("a").IsNot(false)).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.L("a").IsNotNull()).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.L("a").IsNotTrue()).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = depiq.From("test").Where(depiq.L("a").IsNotFalse()).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE (a IS NULL) []
	// SELECT * FROM "test" WHERE (a IS TRUE) []
	// SELECT * FROM "test" WHERE (a IS FALSE) []
	// SELECT * FROM "test" WHERE (a IS NULL) []
	// SELECT * FROM "test" WHERE (a IS TRUE) []
	// SELECT * FROM "test" WHERE (a IS FALSE) []
	// SELECT * FROM "test" WHERE (a IS NOT NULL) []
	// SELECT * FROM "test" WHERE (a IS NOT TRUE) []
	// SELECT * FROM "test" WHERE (a IS NOT FALSE) []
	// SELECT * FROM "test" WHERE (a IS NOT NULL) []
	// SELECT * FROM "test" WHERE (a IS NOT TRUE) []
	// SELECT * FROM "test" WHERE (a IS NOT FALSE) []
}

func ExampleL_betweenComparisons() {
	ds := depiq.From("test").Where(
		depiq.L("(a + b)").Between(depiq.Range(1, 10)),
	)
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(
		depiq.L("(a + b)").NotBetween(depiq.Range(1, 10)),
	)
	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ((a + b) BETWEEN 1 AND 10) []
	// SELECT * FROM "test" WHERE ((a + b) BETWEEN ? AND ?) [1 10]
	// SELECT * FROM "test" WHERE ((a + b) NOT BETWEEN 1 AND 10) []
	// SELECT * FROM "test" WHERE ((a + b) NOT BETWEEN ? AND ?) [1 10]
}

func ExampleLAST() {
	ds := depiq.From("test").Select(depiq.LAST("col"))
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT LAST("col") FROM "test" []
	// SELECT LAST("col") FROM "test" []
}

func ExampleLAST_as() {
	sql, _, _ := depiq.From("test").Select(depiq.LAST("a").As("a")).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT LAST("a") AS "a" FROM "test"
}

func ExampleMAX() {
	ds := depiq.From("test").Select(depiq.MAX("col"))
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT MAX("col") FROM "test" []
	// SELECT MAX("col") FROM "test" []
}

func ExampleMAX_as() {
	sql, _, _ := depiq.From("test").Select(depiq.MAX("a").As("a")).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT MAX("a") AS "a" FROM "test"
}

func ExampleMAX_havingClause() {
	ds := depiq.
		From("test").
		Select(depiq.MAX("a").As("MAX")).
		GroupBy("a").
		Having(depiq.MAX("a").Gt(10))

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT MAX("a") AS "MAX" FROM "test" GROUP BY "a" HAVING (MAX("a") > 10) []
	// SELECT MAX("a") AS "MAX" FROM "test" GROUP BY "a" HAVING (MAX("a") > ?) [10]
}

func ExampleMIN() {
	ds := depiq.From("test").Select(depiq.MIN("col"))
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT MIN("col") FROM "test" []
	// SELECT MIN("col") FROM "test" []
}

func ExampleMIN_as() {
	sql, _, _ := depiq.From("test").Select(depiq.MIN("a").As("a")).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT MIN("a") AS "a" FROM "test"
}

func ExampleMIN_havingClause() {
	ds := depiq.
		From("test").
		Select(depiq.MIN("a").As("MIN")).
		GroupBy("a").
		Having(depiq.MIN("a").Gt(10))

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT MIN("a") AS "MIN" FROM "test" GROUP BY "a" HAVING (MIN("a") > 10) []
	// SELECT MIN("a") AS "MIN" FROM "test" GROUP BY "a" HAVING (MIN("a") > ?) [10]
}

func ExampleOn() {
	ds := depiq.From("test").Join(
		depiq.T("my_table"),
		depiq.On(depiq.I("my_table.fkey").Eq(depiq.I("other_table.id"))),
	)

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "test" INNER JOIN "my_table" ON ("my_table"."fkey" = "other_table"."id") []
	// SELECT * FROM "test" INNER JOIN "my_table" ON ("my_table"."fkey" = "other_table"."id") []
}

func ExampleOn_withEx() {
	ds := depiq.From("test").Join(
		depiq.T("my_table"),
		depiq.On(depiq.Ex{"my_table.fkey": depiq.I("other_table.id")}),
	)

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "test" INNER JOIN "my_table" ON ("my_table"."fkey" = "other_table"."id") []
	// SELECT * FROM "test" INNER JOIN "my_table" ON ("my_table"."fkey" = "other_table"."id") []
}

func ExampleOr() {
	ds := depiq.From("test").Where(
		depiq.Or(
			depiq.C("col").Eq(10),
			depiq.C("col").Eq(20),
		),
	)
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE (("col" = 10) OR ("col" = 20)) []
	// SELECT * FROM "test" WHERE (("col" = ?) OR ("col" = ?)) [10 20]
}

func ExampleOr_withAnd() {
	ds := depiq.From("items").Where(
		depiq.Or(
			depiq.C("a").Gt(10),
			depiq.And(
				depiq.C("b").Eq(100),
				depiq.C("c").Neq("test"),
			),
		),
	)
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "items" WHERE (("a" > 10) OR (("b" = 100) AND ("c" != 'test'))) []
	// SELECT * FROM "items" WHERE (("a" > ?) OR (("b" = ?) AND ("c" != ?))) [10 100 test]
}

func ExampleOr_withExMap() {
	ds := depiq.From("test").Where(
		depiq.Or(
			// Ex will be anded together
			depiq.Ex{
				"col1": 1,
				"col2": true,
			},
			depiq.Ex{
				"col3": nil,
				"col4": "foo",
			},
		),
	)
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ((("col1" = 1) AND ("col2" IS TRUE)) OR (("col3" IS NULL) AND ("col4" = 'foo'))) []
	// SELECT * FROM "test" WHERE ((("col1" = ?) AND ("col2" IS TRUE)) OR (("col3" IS NULL) AND ("col4" = ?))) [1 foo]
}

func ExampleRange_numbers() {
	ds := depiq.From("test").Where(
		depiq.C("col").Between(depiq.Range(1, 10)),
	)
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(
		depiq.C("col").NotBetween(depiq.Range(1, 10)),
	)
	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ("col" BETWEEN 1 AND 10) []
	// SELECT * FROM "test" WHERE ("col" BETWEEN ? AND ?) [1 10]
	// SELECT * FROM "test" WHERE ("col" NOT BETWEEN 1 AND 10) []
	// SELECT * FROM "test" WHERE ("col" NOT BETWEEN ? AND ?) [1 10]
}

func ExampleRange_strings() {
	ds := depiq.From("test").Where(
		depiq.C("col").Between(depiq.Range("a", "z")),
	)
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(
		depiq.C("col").NotBetween(depiq.Range("a", "z")),
	)
	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ("col" BETWEEN 'a' AND 'z') []
	// SELECT * FROM "test" WHERE ("col" BETWEEN ? AND ?) [a z]
	// SELECT * FROM "test" WHERE ("col" NOT BETWEEN 'a' AND 'z') []
	// SELECT * FROM "test" WHERE ("col" NOT BETWEEN ? AND ?) [a z]
}

func ExampleRange_identifiers() {
	ds := depiq.From("test").Where(
		depiq.C("col1").Between(depiq.Range(depiq.C("col2"), depiq.C("col3"))),
	)
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(
		depiq.C("col1").NotBetween(depiq.Range(depiq.C("col2"), depiq.C("col3"))),
	)
	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ("col1" BETWEEN "col2" AND "col3") []
	// SELECT * FROM "test" WHERE ("col1" BETWEEN "col2" AND "col3") []
	// SELECT * FROM "test" WHERE ("col1" NOT BETWEEN "col2" AND "col3") []
	// SELECT * FROM "test" WHERE ("col1" NOT BETWEEN "col2" AND "col3") []
}

func ExampleS() {
	s := depiq.S("test_schema")
	t := s.Table("test")
	sql, args, _ := depiq.
		From(t).
		Select(
			t.Col("col1"),
			t.Col("col2"),
			t.Col("col3"),
		).
		ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT "test_schema"."test"."col1", "test_schema"."test"."col2", "test_schema"."test"."col3" FROM "test_schema"."test" []
}

func ExampleSUM() {
	ds := depiq.From("test").Select(depiq.SUM("col"))
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT SUM("col") FROM "test" []
	// SELECT SUM("col") FROM "test" []
}

func ExampleSUM_as() {
	sql, _, _ := depiq.From("test").Select(depiq.SUM("a").As("a")).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT SUM("a") AS "a" FROM "test"
}

func ExampleSUM_havingClause() {
	ds := depiq.
		From("test").
		Select(depiq.SUM("a").As("SUM")).
		GroupBy("a").
		Having(depiq.SUM("a").Gt(10))

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT SUM("a") AS "SUM" FROM "test" GROUP BY "a" HAVING (SUM("a") > 10) []
	// SELECT SUM("a") AS "SUM" FROM "test" GROUP BY "a" HAVING (SUM("a") > ?) [10]
}

func ExampleStar() {
	ds := depiq.From("test").Select(depiq.Star())

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" []
	// SELECT * FROM "test" []
}

func ExampleT() {
	t := depiq.T("test")
	sql, args, _ := depiq.
		From(t).
		Select(
			t.Col("col1"),
			t.Col("col2"),
			t.Col("col3"),
		).
		ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT "test"."col1", "test"."col2", "test"."col3" FROM "test" []
}

func ExampleUsing() {
	ds := depiq.From("test").Join(
		depiq.T("my_table"),
		depiq.Using("fkey"),
	)

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "test" INNER JOIN "my_table" USING ("fkey") []
	// SELECT * FROM "test" INNER JOIN "my_table" USING ("fkey") []
}

func ExampleUsing_withIdentifier() {
	ds := depiq.From("test").Join(
		depiq.T("my_table"),
		depiq.Using(depiq.C("fkey")),
	)

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "test" INNER JOIN "my_table" USING ("fkey") []
	// SELECT * FROM "test" INNER JOIN "my_table" USING ("fkey") []
}

func ExampleEx() {
	ds := depiq.From("items").Where(
		depiq.Ex{
			"col1": "a",
			"col2": 1,
			"col3": true,
			"col4": false,
			"col5": nil,
			"col6": []string{"a", "b", "c"},
		},
	)

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "items" WHERE (("col1" = 'a') AND ("col2" = 1) AND ("col3" IS TRUE) AND ("col4" IS FALSE) AND ("col5" IS NULL) AND ("col6" IN ('a', 'b', 'c'))) []
	// SELECT * FROM "items" WHERE (("col1" = ?) AND ("col2" = ?) AND ("col3" IS TRUE) AND ("col4" IS FALSE) AND ("col5" IS NULL) AND ("col6" IN (?, ?, ?))) [a 1 a b c]
}

func ExampleEx_withOp() {
	sql, args, _ := depiq.From("items").Where(
		depiq.Ex{
			"col1": depiq.Op{"neq": "a"},
			"col3": depiq.Op{"isNot": true},
			"col6": depiq.Op{"notIn": []string{"a", "b", "c"}},
		},
	).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "items" WHERE (("col1" != 'a') AND ("col3" IS NOT TRUE) AND ("col6" NOT IN ('a', 'b', 'c'))) []
}

func ExampleEx_in() {
	// using an Ex expression map
	sql, _, _ := depiq.From("test").Where(depiq.Ex{
		"a": []string{"a", "b", "c"},
	}).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT * FROM "test" WHERE ("a" IN ('a', 'b', 'c'))
}

func ExampleExOr() {
	sql, args, _ := depiq.From("items").Where(
		depiq.ExOr{
			"col1": "a",
			"col2": 1,
			"col3": true,
			"col4": false,
			"col5": nil,
			"col6": []string{"a", "b", "c"},
		},
	).ToSQL()
	fmt.Println(sql, args)

	// nolint:lll // sql statements are long
	// Output:
	// SELECT * FROM "items" WHERE (("col1" = 'a') OR ("col2" = 1) OR ("col3" IS TRUE) OR ("col4" IS FALSE) OR ("col5" IS NULL) OR ("col6" IN ('a', 'b', 'c'))) []
}

func ExampleExOr_withOp() {
	sql, _, _ := depiq.From("items").Where(depiq.ExOr{
		"col1": depiq.Op{"neq": "a"},
		"col3": depiq.Op{"isNot": true},
		"col6": depiq.Op{"notIn": []string{"a", "b", "c"}},
	}).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("items").Where(depiq.ExOr{
		"col1": depiq.Op{"gt": 1},
		"col2": depiq.Op{"gte": 1},
		"col3": depiq.Op{"lt": 1},
		"col4": depiq.Op{"lte": 1},
	}).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("items").Where(depiq.ExOr{
		"col1": depiq.Op{"like": "a%"},
		"col2": depiq.Op{"notLike": "a%"},
		"col3": depiq.Op{"iLike": "a%"},
		"col4": depiq.Op{"notILike": "a%"},
	}).ToSQL()
	fmt.Println(sql)

	sql, _, _ = depiq.From("items").Where(depiq.ExOr{
		"col1": depiq.Op{"like": regexp.MustCompile("^[ab]")},
		"col2": depiq.Op{"notLike": regexp.MustCompile("^[ab]")},
		"col3": depiq.Op{"iLike": regexp.MustCompile("^[ab]")},
		"col4": depiq.Op{"notILike": regexp.MustCompile("^[ab]")},
	}).ToSQL()
	fmt.Println(sql)

	// Output:
	// SELECT * FROM "items" WHERE (("col1" != 'a') OR ("col3" IS NOT TRUE) OR ("col6" NOT IN ('a', 'b', 'c')))
	// SELECT * FROM "items" WHERE (("col1" > 1) OR ("col2" >= 1) OR ("col3" < 1) OR ("col4" <= 1))
	// SELECT * FROM "items" WHERE (("col1" LIKE 'a%') OR ("col2" NOT LIKE 'a%') OR ("col3" ILIKE 'a%') OR ("col4" NOT ILIKE 'a%'))
	// SELECT * FROM "items" WHERE (("col1" ~ '^[ab]') OR ("col2" !~ '^[ab]') OR ("col3" ~* '^[ab]') OR ("col4" !~* '^[ab]'))
}

func ExampleOp_comparisons() {
	ds := depiq.From("test").Where(depiq.Ex{
		"a": 10,
		"b": depiq.Op{"neq": 10},
		"c": depiq.Op{"gte": 10},
		"d": depiq.Op{"lt": 10},
		"e": depiq.Op{"lte": 10},
	})

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE (("a" = 10) AND ("b" != 10) AND ("c" >= 10) AND ("d" < 10) AND ("e" <= 10)) []
	// SELECT * FROM "test" WHERE (("a" = ?) AND ("b" != ?) AND ("c" >= ?) AND ("d" < ?) AND ("e" <= ?)) [10 10 10 10 10]
}

func ExampleOp_inComparisons() {
	// using an Ex expression map
	ds := depiq.From("test").Where(depiq.Ex{
		"a": depiq.Op{"in": []string{"a", "b", "c"}},
	})

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(depiq.Ex{
		"a": depiq.Op{"notIn": []string{"a", "b", "c"}},
	})

	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ("a" IN ('a', 'b', 'c')) []
	// SELECT * FROM "test" WHERE ("a" IN (?, ?, ?)) [a b c]
	// SELECT * FROM "test" WHERE ("a" NOT IN ('a', 'b', 'c')) []
	// SELECT * FROM "test" WHERE ("a" NOT IN (?, ?, ?)) [a b c]
}

func ExampleOp_likeComparisons() {
	// using an Ex expression map
	ds := depiq.From("test").Where(depiq.Ex{
		"a": depiq.Op{"like": "%a%"},
	})
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(depiq.Ex{
		"a": depiq.Op{"like": regexp.MustCompile("[ab]")},
	})

	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(depiq.Ex{
		"a": depiq.Op{"iLike": "%a%"},
	})

	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(depiq.Ex{
		"a": depiq.Op{"iLike": regexp.MustCompile("[ab]")},
	})

	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(depiq.Ex{
		"a": depiq.Op{"notLike": "%a%"},
	})

	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(depiq.Ex{
		"a": depiq.Op{"notLike": regexp.MustCompile("[ab]")},
	})

	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(depiq.Ex{
		"a": depiq.Op{"notILike": "%a%"},
	})

	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(depiq.Ex{
		"a": depiq.Op{"notILike": regexp.MustCompile("[ab]")},
	})
	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ("a" LIKE '%a%') []
	// SELECT * FROM "test" WHERE ("a" LIKE ?) [%a%]
	// SELECT * FROM "test" WHERE ("a" ~ '[ab]') []
	// SELECT * FROM "test" WHERE ("a" ~ ?) [[ab]]
	// SELECT * FROM "test" WHERE ("a" ILIKE '%a%') []
	// SELECT * FROM "test" WHERE ("a" ILIKE ?) [%a%]
	// SELECT * FROM "test" WHERE ("a" ~* '[ab]') []
	// SELECT * FROM "test" WHERE ("a" ~* ?) [[ab]]
	// SELECT * FROM "test" WHERE ("a" NOT LIKE '%a%') []
	// SELECT * FROM "test" WHERE ("a" NOT LIKE ?) [%a%]
	// SELECT * FROM "test" WHERE ("a" !~ '[ab]') []
	// SELECT * FROM "test" WHERE ("a" !~ ?) [[ab]]
	// SELECT * FROM "test" WHERE ("a" NOT ILIKE '%a%') []
	// SELECT * FROM "test" WHERE ("a" NOT ILIKE ?) [%a%]
	// SELECT * FROM "test" WHERE ("a" !~* '[ab]') []
	// SELECT * FROM "test" WHERE ("a" !~* ?) [[ab]]
}

func ExampleOp_isComparisons() {
	// using an Ex expression map
	ds := depiq.From("test").Where(depiq.Ex{
		"a": true,
	})
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)
	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(depiq.Ex{
		"a": depiq.Op{"is": true},
	})
	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)
	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(depiq.Ex{
		"a": false,
	})
	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)
	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(depiq.Ex{
		"a": depiq.Op{"is": false},
	})
	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)
	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(depiq.Ex{
		"a": nil,
	})
	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)
	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(depiq.Ex{
		"a": depiq.Op{"is": nil},
	})
	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)
	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(depiq.Ex{
		"a": depiq.Op{"isNot": true},
	})
	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)
	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(depiq.Ex{
		"a": depiq.Op{"isNot": false},
	})
	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)
	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(depiq.Ex{
		"a": depiq.Op{"isNot": nil},
	})
	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)
	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ("a" IS TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS FALSE) []
	// SELECT * FROM "test" WHERE ("a" IS FALSE) []
	// SELECT * FROM "test" WHERE ("a" IS FALSE) []
	// SELECT * FROM "test" WHERE ("a" IS FALSE) []
	// SELECT * FROM "test" WHERE ("a" IS NULL) []
	// SELECT * FROM "test" WHERE ("a" IS NULL) []
	// SELECT * FROM "test" WHERE ("a" IS NULL) []
	// SELECT * FROM "test" WHERE ("a" IS NULL) []
	// SELECT * FROM "test" WHERE ("a" IS NOT TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS NOT TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS NOT FALSE) []
	// SELECT * FROM "test" WHERE ("a" IS NOT FALSE) []
	// SELECT * FROM "test" WHERE ("a" IS NOT NULL) []
	// SELECT * FROM "test" WHERE ("a" IS NOT NULL) []
}

func ExampleOp_betweenComparisons() {
	ds := depiq.From("test").Where(depiq.Ex{
		"a": depiq.Op{"between": depiq.Range(1, 10)},
	})
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("test").Where(depiq.Ex{
		"a": depiq.Op{"notBetween": depiq.Range(1, 10)},
	})
	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ("a" BETWEEN 1 AND 10) []
	// SELECT * FROM "test" WHERE ("a" BETWEEN ? AND ?) [1 10]
	// SELECT * FROM "test" WHERE ("a" NOT BETWEEN 1 AND 10) []
	// SELECT * FROM "test" WHERE ("a" NOT BETWEEN ? AND ?) [1 10]
}

// When using a single op with multiple keys they are ORed together
func ExampleOp_withMultipleKeys() {
	ds := depiq.From("items").Where(depiq.Ex{
		"col1": depiq.Op{"is": nil, "eq": 10},
	})

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "items" WHERE (("col1" = 10) OR ("col1" IS NULL)) []
	// SELECT * FROM "items" WHERE (("col1" = ?) OR ("col1" IS NULL)) [10]
}

func ExampleRecord_insert() {
	ds := depiq.Insert("test")

	records := []depiq.Record{
		{"col1": 1, "col2": "foo"},
		{"col1": 2, "col2": "bar"},
	}

	sql, args, _ := ds.Rows(records).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Rows(records).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// INSERT INTO "test" ("col1", "col2") VALUES (1, 'foo'), (2, 'bar') []
	// INSERT INTO "test" ("col1", "col2") VALUES (?, ?), (?, ?) [1 foo 2 bar]
}

func ExampleRecord_update() {
	ds := depiq.Update("test")
	update := depiq.Record{"col1": 1, "col2": "foo"}

	sql, args, _ := ds.Set(update).ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Set(update).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// UPDATE "test" SET "col1"=1,"col2"='foo' []
	// UPDATE "test" SET "col1"=?,"col2"=? [1 foo]
}

func ExampleV() {
	ds := depiq.From("user").Select(
		depiq.V(true).As("is_verified"),
		depiq.V(1.2).As("version"),
		"first_name",
		"last_name",
	)

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("user").Where(depiq.V(1).Neq(1))
	sql, args, _ = ds.ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT TRUE AS "is_verified", 1.2 AS "version", "first_name", "last_name" FROM "user" []
	// SELECT * FROM "user" WHERE (1 != 1) []
}

func ExampleV_prepared() {
	ds := depiq.From("user").Select(
		depiq.V(true).As("is_verified"),
		depiq.V(1.2).As("version"),
		"first_name",
		"last_name",
	)

	sql, args, _ := ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	ds = depiq.From("user").Where(depiq.V(1).Neq(1))

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT ? AS "is_verified", ? AS "version", "first_name", "last_name" FROM "user" [true 1.2]
	// SELECT * FROM "user" WHERE (? != ?) [1 1]
}

func ExampleVals() {
	ds := depiq.Insert("user").
		Cols("first_name", "last_name", "is_verified").
		Vals(
			depiq.Vals{"Greg", "Farley", true},
			depiq.Vals{"Jimmy", "Stewart", true},
			depiq.Vals{"Jeff", "Jeffers", false},
		)
	insertSQL, args, _ := ds.ToSQL()
	fmt.Println(insertSQL, args)

	// Output:
	// INSERT INTO "user" ("first_name", "last_name", "is_verified") VALUES ('Greg', 'Farley', TRUE), ('Jimmy', 'Stewart', TRUE), ('Jeff', 'Jeffers', FALSE) []
}

func ExampleW() {
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
	// Output:
	// SELECT ROW_NUMBER() OVER (PARTITION BY "a" ORDER BY "b" ASC) FROM "test" []
	// SELECT ROW_NUMBER() OVER "w" FROM "test" WINDOW "w" AS (PARTITION BY "a" ORDER BY "b" ASC) []
	// SELECT ROW_NUMBER() OVER "w1" FROM "test" WINDOW "w1" AS (PARTITION BY "a"), "w" AS ("w1" ORDER BY "b" ASC) []
	// SELECT ROW_NUMBER() OVER ("w" ORDER BY "b") FROM "test" WINDOW "w" AS (PARTITION BY "a") []
}

func ExampleLateral() {
	maxEntry := depiq.From("entry").
		Select(depiq.MAX("int").As("max_int")).
		Where(depiq.Ex{"time": depiq.Op{"lt": depiq.I("e.time")}}).
		As("max_entry")

	maxID := depiq.From("entry").
		Select("id").
		Where(depiq.Ex{"int": depiq.I("max_entry.max_int")}).
		As("max_id")

	ds := depiq.
		Select("e.id", "max_entry.max_int", "max_id.id").
		From(
			depiq.T("entry").As("e"),
			depiq.Lateral(maxEntry),
			depiq.Lateral(maxID),
		)
	query, args, _ := ds.ToSQL()
	fmt.Println(query, args)

	query, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(query, args)

	// Output:
	// SELECT "e"."id", "max_entry"."max_int", "max_id"."id" FROM "entry" AS "e", LATERAL (SELECT MAX("int") AS "max_int" FROM "entry" WHERE ("time" < "e"."time")) AS "max_entry", LATERAL (SELECT "id" FROM "entry" WHERE ("int" = "max_entry"."max_int")) AS "max_id" []
	// SELECT "e"."id", "max_entry"."max_int", "max_id"."id" FROM "entry" AS "e", LATERAL (SELECT MAX("int") AS "max_int" FROM "entry" WHERE ("time" < "e"."time")) AS "max_entry", LATERAL (SELECT "id" FROM "entry" WHERE ("int" = "max_entry"."max_int")) AS "max_id" []
}

func ExampleLateral_join() {
	maxEntry := depiq.From("entry").
		Select(depiq.MAX("int").As("max_int")).
		Where(depiq.Ex{"time": depiq.Op{"lt": depiq.I("e.time")}}).
		As("max_entry")

	maxID := depiq.From("entry").
		Select("id").
		Where(depiq.Ex{"int": depiq.I("max_entry.max_int")}).
		As("max_id")

	ds := depiq.
		Select("e.id", "max_entry.max_int", "max_id.id").
		From(depiq.T("entry").As("e")).
		Join(depiq.Lateral(maxEntry), depiq.On(depiq.V(true))).
		Join(depiq.Lateral(maxID), depiq.On(depiq.V(true)))
	query, args, _ := ds.ToSQL()
	fmt.Println(query, args)

	query, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(query, args)

	// Output:
	// SELECT "e"."id", "max_entry"."max_int", "max_id"."id" FROM "entry" AS "e" INNER JOIN LATERAL (SELECT MAX("int") AS "max_int" FROM "entry" WHERE ("time" < "e"."time")) AS "max_entry" ON TRUE INNER JOIN LATERAL (SELECT "id" FROM "entry" WHERE ("int" = "max_entry"."max_int")) AS "max_id" ON TRUE []
	// SELECT "e"."id", "max_entry"."max_int", "max_id"."id" FROM "entry" AS "e" INNER JOIN LATERAL (SELECT MAX("int") AS "max_int" FROM "entry" WHERE ("time" < "e"."time")) AS "max_entry" ON ? INNER JOIN LATERAL (SELECT "id" FROM "entry" WHERE ("int" = "max_entry"."max_int")) AS "max_id" ON ? [true true]
}

func ExampleAny() {
	ds := depiq.From("test").Where(depiq.Ex{
		"id": depiq.Any(depiq.From("other").Select("test_id")),
	})
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "test" WHERE ("id" = ANY ((SELECT "test_id" FROM "other"))) []
	// SELECT * FROM "test" WHERE ("id" = ANY ((SELECT "test_id" FROM "other"))) []
}

func ExampleAll() {
	ds := depiq.From("test").Where(depiq.Ex{
		"id": depiq.All(depiq.From("other").Select("test_id")),
	})
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "test" WHERE ("id" = ALL ((SELECT "test_id" FROM "other"))) []
	// SELECT * FROM "test" WHERE ("id" = ALL ((SELECT "test_id" FROM "other"))) []
}

func ExampleCase_search() {
	ds := depiq.From("test").
		Select(
			depiq.C("col"),
			depiq.Case().
				When(depiq.C("col").Gt(0), true).
				When(depiq.C("col").Lte(0), false).
				As("is_gt_zero"),
		)
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT "col", CASE  WHEN ("col" > 0) THEN TRUE WHEN ("col" <= 0) THEN FALSE END AS "is_gt_zero" FROM "test" []
	// SELECT "col", CASE  WHEN ("col" > ?) THEN ? WHEN ("col" <= ?) THEN ? END AS "is_gt_zero" FROM "test" [0 true 0 false]
}

func ExampleCase_searchElse() {
	ds := depiq.From("test").
		Select(
			depiq.C("col"),
			depiq.Case().
				When(depiq.C("col").Gt(10), "Gt 10").
				When(depiq.C("col").Gt(20), "Gt 20").
				Else("Bad Val").
				As("str_val"),
		)
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT "col", CASE  WHEN ("col" > 10) THEN 'Gt 10' WHEN ("col" > 20) THEN 'Gt 20' ELSE 'Bad Val' END AS "str_val" FROM "test" []
	// SELECT "col", CASE  WHEN ("col" > ?) THEN ? WHEN ("col" > ?) THEN ? ELSE ? END AS "str_val" FROM "test" [10 Gt 10 20 Gt 20 Bad Val]
}

func ExampleCase_value() {
	ds := depiq.From("test").
		Select(
			depiq.C("col"),
			depiq.Case().
				Value(depiq.C("str")).
				When("foo", "FOO").
				When("bar", "BAR").
				As("foo_bar_upper"),
		)
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT "col", CASE "str" WHEN 'foo' THEN 'FOO' WHEN 'bar' THEN 'BAR' END AS "foo_bar_upper" FROM "test" []
	// SELECT "col", CASE "str" WHEN ? THEN ? WHEN ? THEN ? END AS "foo_bar_upper" FROM "test" [foo FOO bar BAR]
}

func ExampleCase_valueElse() {
	ds := depiq.From("test").
		Select(
			depiq.C("col"),
			depiq.Case().
				Value(depiq.C("str")).
				When("foo", "FOO").
				When("bar", "BAR").
				Else("Baz").
				As("foo_bar_upper"),
		)
	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).ToSQL()
	fmt.Println(sql, args)
	// Output:
	// SELECT "col", CASE "str" WHEN 'foo' THEN 'FOO' WHEN 'bar' THEN 'BAR' ELSE 'Baz' END AS "foo_bar_upper" FROM "test" []
	// SELECT "col", CASE "str" WHEN ? THEN ? WHEN ? THEN ? ELSE ? END AS "foo_bar_upper" FROM "test" [foo FOO bar BAR Baz]
}
