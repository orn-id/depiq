package depiq_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/orn-id/depiq"
	"github.com/orn-id/depiq/exp"
	"github.com/orn-id/depiq/internal/errors"
	"github.com/orn-id/depiq/internal/sb"
	"github.com/orn-id/depiq/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type (
	truncateTestCase struct {
		ds      *depiq.TruncateDataset
		clauses exp.TruncateClauses
	}
	truncateDatasetSuite struct {
		suite.Suite
	}
)

func (tds *truncateDatasetSuite) assertCases(cases ...truncateTestCase) {
	for _, s := range cases {
		tds.Equal(s.clauses, s.ds.GetClauses())
	}
}

func (tds *truncateDatasetSuite) TestClone() {
	ds := depiq.Truncate("test")
	tds.Equal(ds, ds.Clone())
}

func (tds *truncateDatasetSuite) TestExpression() {
	ds := depiq.Truncate("test")
	tds.Equal(ds, ds.Expression())
}

func (tds *truncateDatasetSuite) TestDialect() {
	ds := depiq.Truncate("test")
	tds.NotNil(ds.Dialect())
}

func (tds *truncateDatasetSuite) TestWithDialect() {
	ds := depiq.Truncate("test")
	md := new(mocks.SQLDialect)
	ds = ds.SetDialect(md)

	dialect := depiq.GetDialect("default")
	dialectDs := ds.WithDialect("default")
	tds.Equal(md, ds.Dialect())
	tds.Equal(dialect, dialectDs.Dialect())
}

func (tds *truncateDatasetSuite) TestPrepared() {
	ds := depiq.Truncate("test")
	preparedDs := ds.Prepared(true)
	tds.True(preparedDs.IsPrepared())
	tds.False(ds.IsPrepared())
	// should apply the prepared to any datasets created from the root
	tds.True(preparedDs.Restrict().IsPrepared())

	defer depiq.SetDefaultPrepared(false)
	depiq.SetDefaultPrepared(true)

	// should be prepared by default
	ds = depiq.Truncate("test")
	tds.True(ds.IsPrepared())
}

func (tds *truncateDatasetSuite) TestGetClauses() {
	ds := depiq.Truncate("test")
	ce := exp.NewTruncateClauses().SetTable(exp.NewColumnListExpression(depiq.I("test")))
	tds.Equal(ce, ds.GetClauses())
}

func (tds *truncateDatasetSuite) TestTable() {
	bd := depiq.Truncate("test")
	tds.assertCases(
		truncateTestCase{
			ds: bd.Table("test2"),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test2")),
		},
		truncateTestCase{
			ds: bd.Table("test1", "test2"),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test1", "test2")),
		},
		truncateTestCase{
			ds: bd,
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")),
		},
	)
}

func (tds *truncateDatasetSuite) TestCascade() {
	bd := depiq.Truncate("test")
	tds.assertCases(
		truncateTestCase{
			ds: bd.Cascade(),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Cascade: true}),
		},
		truncateTestCase{
			ds: bd.Restrict().Cascade(),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Cascade: true, Restrict: true}),
		},
		truncateTestCase{
			ds: bd,
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")),
		},
	)
}

func (tds *truncateDatasetSuite) TestNoCascade() {
	bd := depiq.Truncate("test").Cascade()
	tds.assertCases(
		truncateTestCase{
			ds: bd.NoCascade(),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{}),
		},
		truncateTestCase{
			ds: bd.Restrict().NoCascade(),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Cascade: false, Restrict: true}),
		},
		truncateTestCase{
			ds: bd,
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Cascade: true}),
		},
	)
}

func (tds *truncateDatasetSuite) TestRestrict() {
	bd := depiq.Truncate("test")
	tds.assertCases(
		truncateTestCase{
			ds: bd.Restrict(),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Restrict: true}),
		},
		truncateTestCase{
			ds: bd.Cascade().Restrict(),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Cascade: true, Restrict: true}),
		},
		truncateTestCase{
			ds: bd,
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")),
		},
	)
}

func (tds *truncateDatasetSuite) TestNoRestrict() {
	bd := depiq.Truncate("test").Restrict()
	tds.assertCases(
		truncateTestCase{
			ds: bd.NoRestrict(),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{}),
		},
		truncateTestCase{
			ds: bd.Cascade().NoRestrict(),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Cascade: true, Restrict: false}),
		},
		truncateTestCase{
			ds: bd,
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Restrict: true}),
		},
	)
}

func (tds *truncateDatasetSuite) TestIdentity() {
	bd := depiq.Truncate("test")
	tds.assertCases(
		truncateTestCase{
			ds: bd.Identity("RESTART"),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Identity: "RESTART"}),
		},
		truncateTestCase{
			ds: bd.Identity("CONTINUE"),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Identity: "CONTINUE"}),
		},
		truncateTestCase{
			ds: bd.Cascade().Restrict().Identity("CONTINUE"),
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")).
				SetOptions(exp.TruncateOptions{Cascade: true, Restrict: true, Identity: "CONTINUE"}),
		},
		truncateTestCase{
			ds: bd,
			clauses: exp.NewTruncateClauses().
				SetTable(exp.NewColumnListExpression("test")),
		},
	)
}

func (tds *truncateDatasetSuite) TestToSQL() {
	md := new(mocks.SQLDialect)
	ds := depiq.Truncate("test").SetDialect(md)
	c := ds.GetClauses()
	sqlB := sb.NewSQLBuilder(false)
	md.On("ToTruncateSQL", sqlB, c).Return(nil).Once()

	sql, args, err := ds.ToSQL()
	tds.NoError(err)
	tds.Empty(sql)
	tds.Empty(args)
	md.AssertExpectations(tds.T())
}

func (tds *truncateDatasetSuite) TestToSQL__withPrepared() {
	md := new(mocks.SQLDialect)
	ds := depiq.Truncate("test").Prepared(true).SetDialect(md)
	c := ds.GetClauses()
	sqlB := sb.NewSQLBuilder(true)
	md.On("ToTruncateSQL", sqlB, c).Return(nil).Once()

	sql, args, err := ds.ToSQL()
	tds.Empty(sql)
	tds.Empty(args)
	tds.Nil(err)
	md.AssertExpectations(tds.T())
}

func (tds *truncateDatasetSuite) TestToSQL_withError() {
	md := new(mocks.SQLDialect)
	ds := depiq.Truncate("test").SetDialect(md)
	c := ds.GetClauses()
	ee := errors.New("expected error")
	sqlB := sb.NewSQLBuilder(false)
	md.On("ToTruncateSQL", sqlB, c).Run(func(args mock.Arguments) {
		args.Get(0).(sb.SQLBuilder).SetError(ee)
	}).Once()

	sql, args, err := ds.ToSQL()
	tds.Empty(sql)
	tds.Empty(args)
	tds.Equal(ee, err)
	md.AssertExpectations(tds.T())
}

func (tds *truncateDatasetSuite) TestExecutor() {
	mDB, _, err := sqlmock.New()
	tds.NoError(err)

	ds := depiq.New("mock", mDB).Truncate("table1", "table2")

	tsql, args, err := ds.Executor().ToSQL()
	tds.NoError(err)
	tds.Empty(args)
	tds.Equal(`TRUNCATE "table1", "table2"`, tsql)

	tsql, args, err = ds.Prepared(true).Executor().ToSQL()
	tds.NoError(err)
	tds.Empty(args)
	tds.Equal(`TRUNCATE "table1", "table2"`, tsql)

	defer depiq.SetDefaultPrepared(false)
	depiq.SetDefaultPrepared(true)

	tsql, args, err = ds.Executor().ToSQL()
	tds.NoError(err)
	tds.Empty(args)
	tds.Equal(`TRUNCATE "table1", "table2"`, tsql)
}

func (tds *truncateDatasetSuite) TestSetError() {
	err1 := errors.New("error #1")
	err2 := errors.New("error #2")
	err3 := errors.New("error #3")

	// Verify initial error set/get works properly
	md := new(mocks.SQLDialect)
	ds := depiq.Truncate("test").SetDialect(md)
	ds = ds.SetError(err1)
	tds.Equal(err1, ds.Error())
	sql, args, err := ds.ToSQL()
	tds.Empty(sql)
	tds.Empty(args)
	tds.Equal(err1, err)

	// Repeated SetError calls on Dataset should not overwrite the original error
	ds = ds.SetError(err2)
	tds.Equal(err1, ds.Error())
	sql, args, err = ds.ToSQL()
	tds.Empty(sql)
	tds.Empty(args)
	tds.Equal(err1, err)

	// Builder functions should not lose the error
	ds = ds.Cascade()
	tds.Equal(err1, ds.Error())
	sql, args, err = ds.ToSQL()
	tds.Empty(sql)
	tds.Empty(args)
	tds.Equal(err1, err)

	// Deeper errors inside SQL generation should still return original error
	c := ds.GetClauses()
	sqlB := sb.NewSQLBuilder(false)
	md.On("ToTruncateSQL", sqlB, c).Run(func(args mock.Arguments) {
		args.Get(0).(sb.SQLBuilder).SetError(err3)
	}).Once()

	sql, args, err = ds.ToSQL()
	tds.Empty(sql)
	tds.Empty(args)
	tds.Equal(err1, err)
}

func TestTruncateDataset(t *testing.T) {
	suite.Run(t, new(truncateDatasetSuite))
}
