package depiq_test

import (
	"testing"

	"github.com/orn-id/depiq/v9"
	"github.com/orn-id/depiq/v9/exp"
	"github.com/stretchr/testify/suite"
)

type (
	goquExpressionsSuite struct {
		suite.Suite
	}
)

func (ges *goquExpressionsSuite) TestCast() {
	ges.Equal(exp.NewCastExpression(depiq.C("test"), "string"), depiq.Cast(depiq.C("test"), "string"))
}

func (ges *goquExpressionsSuite) TestDoNothing() {
	ges.Equal(exp.NewDoNothingConflictExpression(), depiq.DoNothing())
}

func (ges *goquExpressionsSuite) TestDoUpdate() {
	ges.Equal(exp.NewDoUpdateConflictExpression("test", depiq.Record{"a": "b"}), depiq.DoUpdate("test", depiq.Record{"a": "b"}))
}

func (ges *goquExpressionsSuite) TestOr() {
	e1 := depiq.C("a").Eq("b")
	e2 := depiq.C("b").Eq(2)
	ges.Equal(exp.NewExpressionList(exp.OrType, e1, e2), depiq.Or(e1, e2))
}

func (ges *goquExpressionsSuite) TestAnd() {
	e1 := depiq.C("a").Eq("b")
	e2 := depiq.C("b").Eq(2)
	ges.Equal(exp.NewExpressionList(exp.AndType, e1, e2), depiq.And(e1, e2))
}

func (ges *goquExpressionsSuite) TestFunc() {
	ges.Equal(exp.NewSQLFunctionExpression("count", depiq.L("*")), depiq.Func("count", depiq.L("*")))
}

func (ges *goquExpressionsSuite) TestDISTINCT() {
	ges.Equal(exp.NewSQLFunctionExpression("DISTINCT", depiq.I("col")), depiq.DISTINCT("col"))
}

func (ges *goquExpressionsSuite) TestCOUNT() {
	ges.Equal(exp.NewSQLFunctionExpression("COUNT", depiq.I("col")), depiq.COUNT("col"))
}

func (ges *goquExpressionsSuite) TestMIN() {
	ges.Equal(exp.NewSQLFunctionExpression("MIN", depiq.I("col")), depiq.MIN("col"))
}

func (ges *goquExpressionsSuite) TestMAX() {
	ges.Equal(exp.NewSQLFunctionExpression("MAX", depiq.I("col")), depiq.MAX("col"))
}

func (ges *goquExpressionsSuite) TestAVG() {
	ges.Equal(exp.NewSQLFunctionExpression("AVG", depiq.I("col")), depiq.AVG("col"))
}

func (ges *goquExpressionsSuite) TestFIRST() {
	ges.Equal(exp.NewSQLFunctionExpression("FIRST", depiq.I("col")), depiq.FIRST("col"))
}

func (ges *goquExpressionsSuite) TestLAST() {
	ges.Equal(exp.NewSQLFunctionExpression("LAST", depiq.I("col")), depiq.LAST("col"))
}

func (ges *goquExpressionsSuite) TestSUM() {
	ges.Equal(exp.NewSQLFunctionExpression("SUM", depiq.I("col")), depiq.SUM("col"))
}

func (ges *goquExpressionsSuite) TestCOALESCE() {
	ges.Equal(exp.NewSQLFunctionExpression("COALESCE", depiq.I("col"), nil), depiq.COALESCE(depiq.I("col"), nil))
}

func (ges *goquExpressionsSuite) TestROW_NUMBER() {
	ges.Equal(exp.NewSQLFunctionExpression("ROW_NUMBER"), depiq.ROW_NUMBER())
}

func (ges *goquExpressionsSuite) TestRANK() {
	ges.Equal(exp.NewSQLFunctionExpression("RANK"), depiq.RANK())
}

func (ges *goquExpressionsSuite) TestDENSE_RANK() {
	ges.Equal(exp.NewSQLFunctionExpression("DENSE_RANK"), depiq.DENSE_RANK())
}

func (ges *goquExpressionsSuite) TestPERCENT_RANK() {
	ges.Equal(exp.NewSQLFunctionExpression("PERCENT_RANK"), depiq.PERCENT_RANK())
}

func (ges *goquExpressionsSuite) TestCUME_DIST() {
	ges.Equal(exp.NewSQLFunctionExpression("CUME_DIST"), depiq.CUME_DIST())
}

func (ges *goquExpressionsSuite) TestNTILE() {
	ges.Equal(exp.NewSQLFunctionExpression("NTILE", 1), depiq.NTILE(1))
}

func (ges *goquExpressionsSuite) TestFIRST_VALUE() {
	ges.Equal(exp.NewSQLFunctionExpression("FIRST_VALUE", depiq.I("col")), depiq.FIRST_VALUE("col"))
}

func (ges *goquExpressionsSuite) TestLAST_VALUE() {
	ges.Equal(exp.NewSQLFunctionExpression("LAST_VALUE", depiq.I("col")), depiq.LAST_VALUE("col"))
}

func (ges *goquExpressionsSuite) TestNTH_VALUE() {
	ges.Equal(exp.NewSQLFunctionExpression("NTH_VALUE", depiq.I("col"), 1), depiq.NTH_VALUE("col", 1))
	ges.Equal(exp.NewSQLFunctionExpression("NTH_VALUE", depiq.I("col"), 1), depiq.NTH_VALUE(depiq.C("col"), 1))
}

func (ges *goquExpressionsSuite) TestI() {
	ges.Equal(exp.NewIdentifierExpression("s", "t", "c"), depiq.I("s.t.c"))
}

func (ges *goquExpressionsSuite) TestC() {
	ges.Equal(exp.NewIdentifierExpression("", "", "c"), depiq.C("c"))
}

func (ges *goquExpressionsSuite) TestS() {
	ges.Equal(exp.NewIdentifierExpression("s", "", ""), depiq.S("s"))
}

func (ges *goquExpressionsSuite) TestT() {
	ges.Equal(exp.NewIdentifierExpression("", "t", ""), depiq.T("t"))
}

func (ges *goquExpressionsSuite) TestW() {
	ges.Equal(exp.NewWindowExpression(nil, nil, nil, nil), depiq.W())
	ges.Equal(exp.NewWindowExpression(depiq.I("a"), nil, nil, nil), depiq.W("a"))
	ges.Equal(exp.NewWindowExpression(depiq.I("a"), depiq.I("b"), nil, nil), depiq.W("a", "b"))
	ges.Equal(exp.NewWindowExpression(depiq.I("a"), depiq.I("b"), nil, nil), depiq.W("a", "b", "c"))
}

func (ges *goquExpressionsSuite) TestOn() {
	ges.Equal(exp.NewJoinOnCondition(depiq.Ex{"a": "b"}), depiq.On(depiq.Ex{"a": "b"}))
}

func (ges *goquExpressionsSuite) TestUsing() {
	ges.Equal(exp.NewJoinUsingCondition("a", "b"), depiq.Using("a", "b"))
}

func (ges *goquExpressionsSuite) TestL() {
	ges.Equal(exp.NewLiteralExpression("? + ?", 1, 2), depiq.L("? + ?", 1, 2))
}

func (ges *goquExpressionsSuite) TestLiteral() {
	ges.Equal(exp.NewLiteralExpression("? + ?", 1, 2), depiq.Literal("? + ?", 1, 2))
}

func (ges *goquExpressionsSuite) TestV() {
	ges.Equal(exp.NewLiteralExpression("?", "a"), depiq.V("a"))
}

func (ges *goquExpressionsSuite) TestRange() {
	ges.Equal(exp.NewRangeVal("a", "b"), depiq.Range("a", "b"))
}

func (ges *goquExpressionsSuite) TestStar() {
	ges.Equal(exp.NewLiteralExpression("*"), depiq.Star())
}

func (ges *goquExpressionsSuite) TestDefault() {
	ges.Equal(exp.Default(), depiq.Default())
}

func (ges *goquExpressionsSuite) TestLateral() {
	ds := depiq.From("test")
	ges.Equal(exp.NewLateralExpression(ds), depiq.Lateral(ds))
}

func (ges *goquExpressionsSuite) TestAny() {
	ds := depiq.From("test").Select("id")
	ges.Equal(exp.NewSQLFunctionExpression("ANY ", ds), depiq.Any(ds))
}

func (ges *goquExpressionsSuite) TestAll() {
	ds := depiq.From("test").Select("id")
	ges.Equal(exp.NewSQLFunctionExpression("ALL ", ds), depiq.All(ds))
}

func TestGoquExpressions(t *testing.T) {
	suite.Run(t, new(goquExpressionsSuite))
}
