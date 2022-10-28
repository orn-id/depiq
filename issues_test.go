package depiq_test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/orn-id/depiq"
	"github.com/orn-id/depiq/exp"
	"github.com/stretchr/testify/suite"
)

type githubIssuesSuite struct {
	suite.Suite
}

func (gis *githubIssuesSuite) AfterTest(suiteName, testName string) {
	depiq.SetColumnRenameFunction(strings.ToLower)
}

// Test for https://github.com/orn-id/depiq/issues/49
func (gis *githubIssuesSuite) TestIssue49() {
	dialect := depiq.Dialect("default")

	filters := depiq.Or()
	sql, args, err := dialect.From("table").Where(filters).ToSQL()
	gis.NoError(err)
	gis.Empty(args)
	gis.Equal(`SELECT * FROM "table"`, sql)

	sql, args, err = dialect.From("table").Where(depiq.Ex{}).ToSQL()
	gis.NoError(err)
	gis.Empty(args)
	gis.Equal(`SELECT * FROM "table"`, sql)

	sql, args, err = dialect.From("table").Where(depiq.ExOr{}).ToSQL()
	gis.NoError(err)
	gis.Empty(args)
	gis.Equal(`SELECT * FROM "table"`, sql)
}

// Test for https://github.com/orn-id/depiq/issues/115
func (gis *githubIssuesSuite) TestIssue115() {
	type TestStruct struct {
		Field string
	}
	depiq.SetColumnRenameFunction(func(col string) string {
		return ""
	})

	_, _, err := depiq.Insert("test").Rows(TestStruct{Field: "hello"}).ToSQL()
	gis.EqualError(err, `goqu: a empty identifier was encountered, please specify a "schema", "table" or "column"`)
}

// Test for https://github.com/orn-id/depiq/issues/118
func (gis *githubIssuesSuite) TestIssue118_withEmbeddedStructWithoutExportedFields() {
	// struct is in a custom package
	type SimpleRole struct {
		sync.RWMutex
		permissions []string // nolint:structcheck,unused //needed for test
	}

	// .....

	type Role struct {
		*SimpleRole

		ID        string    `json:"id" db:"id" goqu:"skipinsert"`
		Key       string    `json:"key" db:"key"`
		Name      string    `json:"name" db:"name"`
		CreatedAt time.Time `json:"-" db:"created_at" goqu:"skipinsert"`
	}

	rUser := &Role{
		Key:  `user`,
		Name: `User role`,
	}

	sql, arg, err := depiq.Insert(`rbac_roles`).
		Returning(depiq.C(`id`)).
		Rows(rUser).
		ToSQL()
	gis.NoError(err)
	gis.Empty(arg)
	gis.Equal(`INSERT INTO "rbac_roles" ("key", "name") VALUES ('user', 'User role') RETURNING "id"`, sql)

	sql, arg, err = depiq.Update(`rbac_roles`).
		Returning(depiq.C(`id`)).
		Set(rUser).
		ToSQL()
	gis.NoError(err)
	gis.Empty(arg)
	gis.Equal(
		`UPDATE "rbac_roles" SET "created_at"='0001-01-01T00:00:00Z',"id"='',"key"='user',"name"='User role' RETURNING "id"`,
		sql,
	)

	rUser = &Role{
		SimpleRole: &SimpleRole{},
		Key:        `user`,
		Name:       `User role`,
	}

	sql, arg, err = depiq.Insert(`rbac_roles`).
		Returning(depiq.C(`id`)).
		Rows(rUser).
		ToSQL()
	gis.NoError(err)
	gis.Empty(arg)
	gis.Equal(`INSERT INTO "rbac_roles" ("key", "name") VALUES ('user', 'User role') RETURNING "id"`, sql)

	sql, arg, err = depiq.Update(`rbac_roles`).
		Returning(depiq.C(`id`)).
		Set(rUser).
		ToSQL()
	gis.NoError(err)
	gis.Empty(arg)
	gis.Equal(
		`UPDATE "rbac_roles" SET `+
			`"created_at"='0001-01-01T00:00:00Z',"id"='',"key"='user',"name"='User role' RETURNING "id"`,
		sql,
	)
}

// Test for https://github.com/orn-id/depiq/issues/118
func (gis *githubIssuesSuite) TestIssue118_withNilEmbeddedStructWithExportedFields() {
	// struct is in a custom package
	type SimpleRole struct {
		sync.RWMutex
		permissions []string // nolint:structcheck,unused // needed for test
		IDStr       string
	}

	// .....

	type Role struct {
		*SimpleRole

		ID        string    `json:"id" db:"id" goqu:"skipinsert"`
		Key       string    `json:"key" db:"key"`
		Name      string    `json:"name" db:"name"`
		CreatedAt time.Time `json:"-" db:"created_at" goqu:"skipinsert"`
	}

	rUser := &Role{
		Key:  `user`,
		Name: `User role`,
	}
	sql, arg, err := depiq.Insert(`rbac_roles`).
		Returning(depiq.C(`id`)).
		Rows(rUser).
		ToSQL()
	gis.NoError(err)
	gis.Empty(arg)
	// it should not insert fields on nil embedded pointers
	gis.Equal(`INSERT INTO "rbac_roles" ("key", "name") VALUES ('user', 'User role') RETURNING "id"`, sql)

	sql, arg, err = depiq.Update(`rbac_roles`).
		Returning(depiq.C(`id`)).
		Set(rUser).
		ToSQL()
	gis.NoError(err)
	gis.Empty(arg)
	// it should not insert fields on nil embedded pointers
	gis.Equal(
		`UPDATE "rbac_roles" SET "created_at"='0001-01-01T00:00:00Z',"id"='',"key"='user',"name"='User role' RETURNING "id"`,
		sql,
	)

	rUser = &Role{
		SimpleRole: &SimpleRole{},
		Key:        `user`,
		Name:       `User role`,
	}
	sql, arg, err = depiq.Insert(`rbac_roles`).
		Returning(depiq.C(`id`)).
		Rows(rUser).
		ToSQL()
	gis.NoError(err)
	gis.Empty(arg)
	// it should not insert fields on nil embedded pointers
	gis.Equal(
		`INSERT INTO "rbac_roles" ("idstr", "key", "name") VALUES ('', 'user', 'User role') RETURNING "id"`,
		sql,
	)

	sql, arg, err = depiq.Update(`rbac_roles`).
		Returning(depiq.C(`id`)).
		Set(rUser).
		ToSQL()
	gis.NoError(err)
	gis.Empty(arg)
	// it should not insert fields on nil embedded pointers
	gis.Equal(
		`UPDATE "rbac_roles" SET `+
			`"created_at"='0001-01-01T00:00:00Z',"id"='',"idstr"='',"key"='user',"name"='User role' RETURNING "id"`,
		sql,
	)
}

// Test for https://github.com/orn-id/depiq/issues/118
func (gis *githubIssuesSuite) TestIssue140() {
	sql, arg, err := depiq.Insert(`test`).Returning().ToSQL()
	gis.NoError(err)
	gis.Empty(arg)
	gis.Equal(`INSERT INTO "test" DEFAULT VALUES`, sql)

	sql, arg, err = depiq.Update(`test`).Set(depiq.Record{"a": "b"}).Returning().ToSQL()
	gis.NoError(err)
	gis.Empty(arg)
	gis.Equal(
		`UPDATE "test" SET "a"='b'`,
		sql,
	)

	sql, arg, err = depiq.Delete(`test`).Returning().ToSQL()
	gis.NoError(err)
	gis.Empty(arg)
	gis.Equal(
		`DELETE FROM "test"`,
		sql,
	)

	sql, arg, err = depiq.Insert(`test`).Returning(nil).ToSQL()
	gis.NoError(err)
	gis.Empty(arg)
	gis.Equal(`INSERT INTO "test" DEFAULT VALUES`, sql)

	sql, arg, err = depiq.Update(`test`).Set(depiq.Record{"a": "b"}).Returning(nil).ToSQL()
	gis.NoError(err)
	gis.Empty(arg)
	gis.Equal(
		`UPDATE "test" SET "a"='b'`,
		sql,
	)

	sql, arg, err = depiq.Delete(`test`).Returning(nil).ToSQL()
	gis.NoError(err)
	gis.Empty(arg)
	gis.Equal(
		`DELETE FROM "test"`,
		sql,
	)
}

// Test for https://github.com/orn-id/depiq/issues/164
func (gis *githubIssuesSuite) TestIssue164() {
	insertDs := depiq.Insert("foo").Rows(depiq.Record{"user_id": 10}).Returning("id")

	ds := depiq.From("bar").
		With("ins", insertDs).
		Select("bar_name").
		Where(depiq.Ex{"bar.user_id": depiq.I("ins.user_id")})

	sql, args, err := ds.ToSQL()
	gis.NoError(err)
	gis.Empty(args)
	gis.Equal(
		`WITH ins AS (INSERT INTO "foo" ("user_id") VALUES (10) RETURNING "id") `+
			`SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "ins"."user_id")`,
		sql,
	)

	sql, args, err = ds.Prepared(true).ToSQL()
	gis.NoError(err)
	gis.Equal([]interface{}{int64(10)}, args)
	gis.Equal(
		`WITH ins AS (INSERT INTO "foo" ("user_id") VALUES (?) RETURNING "id")`+
			` SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "ins"."user_id")`,
		sql,
	)

	updateDs := depiq.Update("foo").Set(depiq.Record{"bar": "baz"}).Returning("id")

	ds = depiq.From("bar").
		With("upd", updateDs).
		Select("bar_name").
		Where(depiq.Ex{"bar.user_id": depiq.I("upd.user_id")})

	sql, args, err = ds.ToSQL()
	gis.NoError(err)
	gis.Empty(args)
	gis.Equal(
		`WITH upd AS (UPDATE "foo" SET "bar"='baz' RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "upd"."user_id")`,
		sql,
	)

	sql, args, err = ds.Prepared(true).ToSQL()
	gis.NoError(err)
	gis.Equal([]interface{}{"baz"}, args)
	gis.Equal(
		`WITH upd AS (UPDATE "foo" SET "bar"=? RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "upd"."user_id")`,
		sql,
	)

	deleteDs := depiq.Delete("foo").Where(depiq.Ex{"bar": "baz"}).Returning("id")

	ds = depiq.From("bar").
		With("del", deleteDs).
		Select("bar_name").
		Where(depiq.Ex{"bar.user_id": depiq.I("del.user_id")})

	sql, args, err = ds.ToSQL()
	gis.NoError(err)
	gis.Empty(args)
	gis.Equal(
		`WITH del AS (DELETE FROM "foo" WHERE ("bar" = 'baz') RETURNING "id")`+
			` SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "del"."user_id")`,
		sql,
	)

	sql, args, err = ds.Prepared(true).ToSQL()
	gis.NoError(err)
	gis.Equal([]interface{}{"baz"}, args)
	gis.Equal(
		`WITH del AS (DELETE FROM "foo" WHERE ("bar" = ?) RETURNING "id")`+
			` SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "del"."user_id")`,
		sql,
	)
}

// Test for https://github.com/orn-id/depiq/issues/177
func (gis *githubIssuesSuite) TestIssue177() {
	ds := depiq.Dialect("postgres").
		From("ins1").
		With("ins1",
			depiq.Dialect("postgres").
				Insert("account").
				Rows(depiq.Record{"email": "email@email.com", "status": "active", "uuid": "XXX-XXX-XXXX"}).
				Returning("*"),
		).
		With("ins2",
			depiq.Dialect("postgres").
				Insert("account_user").
				Cols("account_id", "user_id").
				FromQuery(depiq.Dialect("postgres").
					From("ins1").
					Select(
						"id",
						depiq.V(1001),
					),
				),
		).
		Select("*")
	sql, args, err := ds.ToSQL()
	gis.NoError(err)
	gis.Equal(`WITH ins1 AS (`+
		`INSERT INTO "account" ("email", "status", "uuid") VALUES ('email@email.com', 'active', 'XXX-XXX-XXXX') RETURNING *),`+
		` ins2 AS (INSERT INTO "account_user" ("account_id", "user_id") SELECT "id", 1001 FROM "ins1")`+
		` SELECT * FROM "ins1"`, sql)
	gis.Len(args, 0)

	sql, args, err = ds.Prepared(true).ToSQL()
	gis.NoError(err)
	gis.Equal(`WITH ins1 AS (INSERT INTO "account" ("email", "status", "uuid") VALUES ($1, $2, $3) RETURNING *), ins2`+
		` AS (INSERT INTO "account_user" ("account_id", "user_id") SELECT "id", $4 FROM "ins1") SELECT * FROM "ins1"`, sql)
	gis.Equal(args, []interface{}{"email@email.com", "active", "XXX-XXX-XXXX", int64(1001)})
}

// Test for https://github.com/orn-id/depiq/issues/183
func (gis *githubIssuesSuite) TestIssue184() {
	expectedErr := fmt.Errorf("an error")
	testCases := []struct {
		ds exp.AppendableExpression
	}{
		{ds: depiq.From("test").As("t").SetError(expectedErr)},
		{ds: depiq.Insert("test").Rows(depiq.Record{"foo": "bar"}).Returning("foo").SetError(expectedErr)},
		{ds: depiq.Update("test").Set(depiq.Record{"foo": "bar"}).Returning("foo").SetError(expectedErr)},
		{ds: depiq.Update("test").Set(depiq.Record{"foo": "bar"}).Returning("foo").SetError(expectedErr)},
		{ds: depiq.Delete("test").Returning("foo").SetError(expectedErr)},
	}

	for _, tc := range testCases {
		ds := depiq.From(tc.ds)
		sql, args, err := ds.ToSQL()
		gis.Equal(expectedErr, err)
		gis.Empty(sql)
		gis.Empty(args)

		sql, args, err = ds.Prepared(true).ToSQL()
		gis.Equal(expectedErr, err)
		gis.Empty(sql)
		gis.Empty(args)

		ds = depiq.From("test2").Where(depiq.Ex{"foo": tc.ds})

		sql, args, err = ds.ToSQL()
		gis.Equal(expectedErr, err)
		gis.Empty(sql)
		gis.Empty(args)

		sql, args, err = ds.Prepared(true).ToSQL()
		gis.Equal(expectedErr, err)
		gis.Empty(sql)
		gis.Empty(args)
	}
}

// Test for https://github.com/orn-id/depiq/issues/185
func (gis *githubIssuesSuite) TestIssue185() {
	mDB, sqlMock, err := sqlmock.New()
	gis.NoError(err)
	sqlMock.ExpectQuery(
		`SELECT \* FROM \(SELECT "id" FROM "table" ORDER BY "id" ASC\) AS "t1" UNION 
\(SELECT \* FROM \(SELECT "id" FROM "table" ORDER BY "id" ASC\) AS "t1"\)`,
	).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1\n2\n3\n4\n"))
	db := depiq.New("mock", mDB)

	ds := db.Select("id").From("table").Order(depiq.C("id").Asc()).
		Union(
			db.Select("id").From("table").Order(depiq.C("id").Asc()),
		)

	ctx := context.Background()
	var i []int
	gis.NoError(ds.ScanValsContext(ctx, &i))
	gis.Equal([]int{1, 2, 3, 4}, i)
}

// Test for https://github.com/orn-id/depiq/issues/203
func (gis *githubIssuesSuite) TestIssue203() {
	// Schema definitions.
	authSchema := depiq.S("company_auth")

	// Table definitions
	usersTable := authSchema.Table("users")

	u := usersTable.As("u")

	ds := depiq.From(u).Select(
		u.Col("id"),
		u.Col("name"),
		u.Col("created_at"),
		u.Col("updated_at"),
	)

	sql, args, err := ds.ToSQL()
	gis.NoError(err)
	gis.Equal(`SELECT "u"."id", "u"."name", "u"."created_at", "u"."updated_at" FROM "company_auth"."users" AS "u"`, sql)
	gis.Empty(args, []interface{}{})

	sql, args, err = ds.Prepared(true).ToSQL()
	gis.NoError(err)
	gis.Equal(`SELECT "u"."id", "u"."name", "u"."created_at", "u"."updated_at" FROM "company_auth"."users" AS "u"`, sql)
	gis.Empty(args, []interface{}{})
}

func (gis *githubIssuesSuite) TestIssue290() {
	type OcomModel struct {
		ID           uint      `json:"id" db:"id" goqu:"skipinsert"`
		CreatedDate  time.Time `json:"created_date" db:"created_date" goqu:"skipupdate"`
		ModifiedDate time.Time `json:"modified_date" db:"modified_date"`
	}

	type ActiveModel struct {
		OcomModel
		ActiveStartDate time.Time  `json:"active_start_date" db:"active_start_date"`
		ActiveEndDate   *time.Time `json:"active_end_date" db:"active_end_date"`
	}

	type CodeModel struct {
		ActiveModel

		Code        string `json:"code" db:"code"`
		Description string `json:"description" binding:"required" db:"description"`
	}

	type CodeExample struct {
		CodeModel
	}

	var item CodeExample
	item.Code = "Code"
	item.Description = "Description"
	item.ID = 1 // Value set HERE!
	item.CreatedDate = time.Date(
		2021, 1, 1, 1, 1, 1, 1, time.UTC)
	item.ModifiedDate = time.Date(
		2021, 2, 2, 2, 2, 2, 2, time.UTC) // The Value we Get!
	item.ActiveStartDate = time.Date(
		2021, 3, 3, 3, 3, 3, 3, time.UTC)

	updateQuery := depiq.From("example").Update().Set(item).Where(depiq.C("id").Eq(1))

	sql, params, err := updateQuery.ToSQL()

	gis.NoError(err)
	gis.Empty(params)
	gis.Equal(`UPDATE "example" SET "active_end_date"=NULL,"active_start_date"='2021-03-03T03:03:03.000000003Z',"code"='Code',"description"='Description',"id"=1,"modified_date"='2021-02-02T02:02:02.000000002Z' WHERE ("id" = 1)`, sql) //nolint:lll
}

func TestGithubIssuesSuite(t *testing.T) {
	suite.Run(t, new(githubIssuesSuite))
}
