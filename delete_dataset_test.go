package depiq_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/orn-id/depiq/v9"
	"github.com/orn-id/depiq/v9/exp"
	"github.com/orn-id/depiq/v9/internal/errors"
	"github.com/orn-id/depiq/v9/internal/sb"
	"github.com/orn-id/depiq/v9/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type (
	deleteTestCase struct {
		ds      *depiq.DeleteDataset
		clauses exp.DeleteClauses
	}
	deleteDatasetSuite struct {
		suite.Suite
	}
)

func (dds *deleteDatasetSuite) assertCases(cases ...deleteTestCase) {
	for _, s := range cases {
		dds.Equal(s.clauses, s.ds.GetClauses())
	}
}

func (dds *deleteDatasetSuite) SetupSuite() {
	noReturn := depiq.DefaultDialectOptions()
	noReturn.SupportsReturn = false
	depiq.RegisterDialect("no-return", noReturn)

	limitOnDelete := depiq.DefaultDialectOptions()
	limitOnDelete.SupportsLimitOnDelete = true
	depiq.RegisterDialect("limit-on-delete", limitOnDelete)

	orderOnDelete := depiq.DefaultDialectOptions()
	orderOnDelete.SupportsOrderByOnDelete = true
	depiq.RegisterDialect("order-on-delete", orderOnDelete)
}

func (dds *deleteDatasetSuite) TearDownSuite() {
	depiq.DeregisterDialect("no-return")
	depiq.DeregisterDialect("limit-on-delete")
	depiq.DeregisterDialect("order-on-delete")
}

func (dds *deleteDatasetSuite) TestDelete() {
	ds := depiq.Delete("test")
	dds.IsType(&depiq.DeleteDataset{}, ds)
	dds.Implements((*exp.Expression)(nil), ds)
	dds.Implements((*exp.AppendableExpression)(nil), ds)
}

func (dds *deleteDatasetSuite) TestClone() {
	ds := depiq.Delete("test")
	dds.Equal(ds.Clone(), ds)
}

func (dds *deleteDatasetSuite) TestExpression() {
	ds := depiq.Delete("test")
	dds.Equal(ds.Expression(), ds)
}

func (dds *deleteDatasetSuite) TestDialect() {
	ds := depiq.Delete("test")
	dds.NotNil(ds.Dialect())
}

func (dds *deleteDatasetSuite) TestWithDialect() {
	ds := depiq.Delete("test")
	md := new(mocks.SQLDialect)
	ds = ds.SetDialect(md)

	dialect := depiq.GetDialect("default")
	dialectDs := ds.WithDialect("default")
	dds.Equal(md, ds.Dialect())
	dds.Equal(dialect, dialectDs.Dialect())
}

func (dds *deleteDatasetSuite) TestPrepared() {
	ds := depiq.Delete("test")
	preparedDs := ds.Prepared(true)
	dds.True(preparedDs.IsPrepared())
	dds.False(ds.IsPrepared())
	// should apply the prepared to any datasets created from the root
	dds.True(preparedDs.Where(depiq.Ex{"a": 1}).IsPrepared())

	defer depiq.SetDefaultPrepared(false)
	depiq.SetDefaultPrepared(true)

	// should be prepared by default
	ds = depiq.Delete("test")
	dds.True(ds.IsPrepared())
}

func (dds *deleteDatasetSuite) TestGetClauses() {
	ds := depiq.Delete("test")
	ce := exp.NewDeleteClauses().SetFrom(depiq.I("test"))
	dds.Equal(ce, ds.GetClauses())
}

func (dds *deleteDatasetSuite) TestWith() {
	from := depiq.From("cte")
	bd := depiq.Delete("items")
	dds.assertCases(
		deleteTestCase{
			ds: bd.With("test-cte", from),
			clauses: exp.NewDeleteClauses().SetFrom(depiq.C("items")).
				CommonTablesAppend(exp.NewCommonTableExpression(false, "test-cte", from)),
		},
		deleteTestCase{
			ds:      bd,
			clauses: exp.NewDeleteClauses().SetFrom(depiq.C("items")),
		},
	)
}

func (dds *deleteDatasetSuite) TestWithRecursive() {
	from := depiq.From("cte")
	bd := depiq.Delete("items")
	dds.assertCases(
		deleteTestCase{
			ds: bd.WithRecursive("test-cte", from),
			clauses: exp.NewDeleteClauses().SetFrom(depiq.C("items")).
				CommonTablesAppend(exp.NewCommonTableExpression(true, "test-cte", from)),
		},
		deleteTestCase{
			ds:      bd,
			clauses: exp.NewDeleteClauses().SetFrom(depiq.C("items")),
		},
	)
}

func (dds *deleteDatasetSuite) TestFrom_withIdentifier() {
	bd := depiq.Delete("items")
	dds.assertCases(
		deleteTestCase{
			ds:      bd.From("items2"),
			clauses: exp.NewDeleteClauses().SetFrom(depiq.C("items2")),
		},
		deleteTestCase{
			ds:      bd.From(depiq.C("items2")),
			clauses: exp.NewDeleteClauses().SetFrom(depiq.C("items2")),
		},
		deleteTestCase{
			ds:      bd.From(depiq.T("items2")),
			clauses: exp.NewDeleteClauses().SetFrom(depiq.T("items2")),
		},
		deleteTestCase{
			ds:      bd.From("schema.table"),
			clauses: exp.NewDeleteClauses().SetFrom(depiq.I("schema.table")),
		},
		deleteTestCase{
			ds:      bd,
			clauses: exp.NewDeleteClauses().SetFrom(depiq.C("items")),
		},
	)

	dds.PanicsWithValue(depiq.ErrBadFromArgument, func() {
		depiq.Delete("test").From(true)
	})
}

func (dds *deleteDatasetSuite) TestWhere() {
	bd := depiq.Delete("items")
	dds.assertCases(
		deleteTestCase{
			ds: bd.Where(depiq.Ex{"a": 1}),
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("items")).
				WhereAppend(depiq.Ex{"a": 1}),
		},
		deleteTestCase{
			ds: bd.Where(depiq.Ex{"a": 1}).Where(depiq.C("b").Eq("c")),
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("items")).
				WhereAppend(depiq.Ex{"a": 1}).
				WhereAppend(depiq.C("b").Eq("c")),
		},
		deleteTestCase{
			ds:      bd,
			clauses: exp.NewDeleteClauses().SetFrom(depiq.C("items")),
		},
	)
}

func (dds *deleteDatasetSuite) TestClearWhere() {
	bd := depiq.Delete("items").Where(depiq.Ex{"a": 1})
	dds.assertCases(
		deleteTestCase{
			ds: bd.ClearWhere(),
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("items")),
		},
		deleteTestCase{
			ds: bd,
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("items")).
				WhereAppend(depiq.Ex{"a": 1}),
		},
	)
}

func (dds *deleteDatasetSuite) TestOrder() {
	bd := depiq.Delete("items")
	dds.assertCases(
		deleteTestCase{
			ds: bd.Order(depiq.C("a").Asc()),
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("items")).
				SetOrder(depiq.C("a").Asc()),
		},
		deleteTestCase{
			ds: bd.Order(depiq.C("a").Asc()).Order(depiq.C("b").Desc()),
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("items")).
				SetOrder(depiq.C("b").Desc()),
		},
		deleteTestCase{
			ds: bd.Order(depiq.C("a").Asc(), depiq.C("b").Desc()),
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("items")).
				SetOrder(depiq.C("a").Asc(), depiq.C("b").Desc()),
		},
		deleteTestCase{
			ds:      bd,
			clauses: exp.NewDeleteClauses().SetFrom(depiq.C("items")),
		},
	)
}

func (dds *deleteDatasetSuite) TestOrderAppend() {
	bd := depiq.Delete("items").Order(depiq.C("a").Asc())
	dds.assertCases(
		deleteTestCase{
			ds: bd.OrderAppend(depiq.C("b").Desc()),
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("items")).
				SetOrder(depiq.C("a").Asc(), depiq.C("b").Desc()),
		},
		deleteTestCase{
			ds: bd,
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("items")).
				SetOrder(depiq.C("a").Asc()),
		},
	)
}

func (dds *deleteDatasetSuite) TestOrderPrepend() {
	bd := depiq.Delete("items").Order(depiq.C("a").Asc())
	dds.assertCases(
		deleteTestCase{
			ds: bd.OrderPrepend(depiq.C("b").Desc()),
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("items")).
				SetOrder(depiq.C("b").Desc(), depiq.C("a").Asc()),
		},
		deleteTestCase{
			ds: bd,
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("items")).
				SetOrder(depiq.C("a").Asc()),
		},
	)
}

func (dds *deleteDatasetSuite) TestClearOrder() {
	bd := depiq.Delete("items").Order(depiq.C("a").Asc())
	dds.assertCases(
		deleteTestCase{
			ds:      bd.ClearOrder(),
			clauses: exp.NewDeleteClauses().SetFrom(depiq.C("items")),
		},
		deleteTestCase{
			ds: bd,
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("items")).
				SetOrder(depiq.C("a").Asc()),
		},
	)
}

func (dds *deleteDatasetSuite) TestLimit() {
	bd := depiq.Delete("test")
	dds.assertCases(
		deleteTestCase{
			ds: bd.Limit(10),
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("test")).
				SetLimit(uint(10)),
		},
		deleteTestCase{
			ds:      bd.Limit(0),
			clauses: exp.NewDeleteClauses().SetFrom(depiq.C("test")),
		},
		deleteTestCase{
			ds: bd.Limit(10).Limit(2),
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("test")).
				SetLimit(uint(2)),
		},
		deleteTestCase{
			ds:      bd.Limit(10).Limit(0),
			clauses: exp.NewDeleteClauses().SetFrom(depiq.C("test")),
		},
		deleteTestCase{
			ds:      bd,
			clauses: exp.NewDeleteClauses().SetFrom(depiq.C("test")),
		},
	)
}

func (dds *deleteDatasetSuite) TestLimitAll() {
	bd := depiq.Delete("test")
	dds.assertCases(
		deleteTestCase{
			ds: bd.LimitAll(),
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("test")).
				SetLimit(depiq.L("ALL")),
		},
		deleteTestCase{
			ds: bd.Limit(10).LimitAll(),
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("test")).
				SetLimit(depiq.L("ALL")),
		},
		deleteTestCase{
			ds:      bd,
			clauses: exp.NewDeleteClauses().SetFrom(depiq.C("test")),
		},
	)
}

func (dds *deleteDatasetSuite) TestClearLimit() {
	bd := depiq.Delete("test").Limit(10)
	dds.assertCases(
		deleteTestCase{
			ds:      bd.ClearLimit(),
			clauses: exp.NewDeleteClauses().SetFrom(depiq.C("test")),
		},
		deleteTestCase{
			ds:      bd,
			clauses: exp.NewDeleteClauses().SetFrom(depiq.C("test")).SetLimit(uint(10)),
		},
	)
}

func (dds *deleteDatasetSuite) TestReturning() {
	bd := depiq.Delete("items")
	dds.assertCases(
		deleteTestCase{
			ds: bd.Returning("a"),
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("items")).
				SetReturning(exp.NewColumnListExpression("a")),
		},
		deleteTestCase{
			ds: bd.Returning(),
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("items")).
				SetReturning(exp.NewColumnListExpression()),
		},
		deleteTestCase{
			ds: bd.Returning(nil),
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("items")).
				SetReturning(exp.NewColumnListExpression()),
		},
		deleteTestCase{
			ds: bd.Returning("a").Returning("b"),
			clauses: exp.NewDeleteClauses().
				SetFrom(depiq.C("items")).
				SetReturning(exp.NewColumnListExpression("b")),
		},
		deleteTestCase{
			ds:      bd,
			clauses: exp.NewDeleteClauses().SetFrom(depiq.C("items")),
		},
	)
}

func (dds *deleteDatasetSuite) TestReturnsColumns() {
	ds := depiq.Delete("test")
	dds.False(ds.ReturnsColumns())
	dds.True(ds.Returning("foo", "bar").ReturnsColumns())
}

func (dds *deleteDatasetSuite) TestToSQL() {
	md := new(mocks.SQLDialect)
	ds := depiq.Delete("test").SetDialect(md)
	c := ds.GetClauses()
	sqlB := sb.NewSQLBuilder(false)
	md.On("ToDeleteSQL", sqlB, c).Return(nil).Once()

	sql, args, err := ds.ToSQL()
	dds.Empty(sql)
	dds.Empty(args)
	dds.Nil(err)
	md.AssertExpectations(dds.T())
}

func (dds *deleteDatasetSuite) TestToSQL_Prepared() {
	md := new(mocks.SQLDialect)
	ds := depiq.Delete("test").Prepared(true).SetDialect(md)
	c := ds.GetClauses()
	sqlB := sb.NewSQLBuilder(true)
	md.On("ToDeleteSQL", sqlB, c).Return(nil).Once()

	sql, args, err := ds.ToSQL()
	dds.Empty(sql)
	dds.Empty(args)
	dds.Nil(err)
	md.AssertExpectations(dds.T())
}

func (dds *deleteDatasetSuite) TestToSQL_WithError() {
	md := new(mocks.SQLDialect)
	ds := depiq.Delete("test").SetDialect(md)
	c := ds.GetClauses()
	ee := errors.New("expected error")
	sqlB := sb.NewSQLBuilder(false)
	md.On("ToDeleteSQL", sqlB, c).Run(func(args mock.Arguments) {
		args.Get(0).(sb.SQLBuilder).SetError(ee)
	}).Once()

	sql, args, err := ds.ToSQL()
	dds.Empty(sql)
	dds.Empty(args)
	dds.Equal(ee, err)
	md.AssertExpectations(dds.T())
}

func (dds *deleteDatasetSuite) TestExecutor() {
	mDB, _, err := sqlmock.New()
	dds.NoError(err)

	ds := depiq.New("mock", mDB).Delete("items").Where(depiq.Ex{"id": depiq.Op{"gt": 10}})

	dsql, args, err := ds.Executor().ToSQL()
	dds.NoError(err)
	dds.Empty(args)
	dds.Equal(`DELETE FROM "items" WHERE ("id" > 10)`, dsql)

	dsql, args, err = ds.Prepared(true).Executor().ToSQL()
	dds.NoError(err)
	dds.Equal([]interface{}{int64(10)}, args)
	dds.Equal(`DELETE FROM "items" WHERE ("id" > ?)`, dsql)

	defer depiq.SetDefaultPrepared(false)
	depiq.SetDefaultPrepared(true)

	dsql, args, err = ds.Executor().ToSQL()
	dds.NoError(err)
	dds.Equal([]interface{}{int64(10)}, args)
	dds.Equal(`DELETE FROM "items" WHERE ("id" > ?)`, dsql)
}

func (dds *deleteDatasetSuite) TestSetError() {
	err1 := errors.New("error #1")
	err2 := errors.New("error #2")
	err3 := errors.New("error #3")

	// Verify initial error set/get works properly
	md := new(mocks.SQLDialect)
	ds := depiq.Delete("test").SetDialect(md)
	ds = ds.SetError(err1)
	dds.Equal(err1, ds.Error())
	sql, args, err := ds.ToSQL()
	dds.Empty(sql)
	dds.Empty(args)
	dds.Equal(err1, err)

	// Repeated SetError calls on Dataset should not overwrite the original error
	ds = ds.SetError(err2)
	dds.Equal(err1, ds.Error())
	sql, args, err = ds.ToSQL()
	dds.Empty(sql)
	dds.Empty(args)
	dds.Equal(err1, err)

	// Builder functions should not lose the error
	ds = ds.ClearLimit()
	dds.Equal(err1, ds.Error())
	sql, args, err = ds.ToSQL()
	dds.Empty(sql)
	dds.Empty(args)
	dds.Equal(err1, err)

	// Deeper errors inside SQL generation should still return original error
	c := ds.GetClauses()
	sqlB := sb.NewSQLBuilder(false)
	md.On("ToDeleteSQL", sqlB, c).Run(func(args mock.Arguments) {
		args.Get(0).(sb.SQLBuilder).SetError(err3)
	}).Once()

	sql, args, err = ds.ToSQL()
	dds.Empty(sql)
	dds.Empty(args)
	dds.Equal(err1, err)
}

func TestDeleteDataset(t *testing.T) {
	suite.Run(t, new(deleteDatasetSuite))
}
