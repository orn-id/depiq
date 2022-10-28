package depiq_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/orn-id/depiq/v9"
	"github.com/stretchr/testify/suite"
)

type (
	dialectWrapperSuite struct {
		suite.Suite
	}
)

func (dws *dialectWrapperSuite) SetupSuite() {
	testDialect := depiq.DefaultDialectOptions()
	// override to some value to ensure correct dialect is set
	depiq.RegisterDialect("test", testDialect)
}

func (dws *dialectWrapperSuite) TearDownSuite() {
	depiq.DeregisterDialect("test")
}

func (dws *dialectWrapperSuite) TestFrom() {
	dw := depiq.Dialect("test")
	dws.Equal(depiq.From("table").WithDialect("test"), dw.From("table"))
}

func (dws *dialectWrapperSuite) TestSelect() {
	dw := depiq.Dialect("test")
	dws.Equal(depiq.Select("col").WithDialect("test"), dw.Select("col"))
}

func (dws *dialectWrapperSuite) TestInsert() {
	dw := depiq.Dialect("test")
	dws.Equal(depiq.Insert("table").WithDialect("test"), dw.Insert("table"))
}

func (dws *dialectWrapperSuite) TestDelete() {
	dw := depiq.Dialect("test")
	dws.Equal(depiq.Delete("table").WithDialect("test"), dw.Delete("table"))
}

func (dws *dialectWrapperSuite) TestTruncate() {
	dw := depiq.Dialect("test")
	dws.Equal(depiq.Truncate("table").WithDialect("test"), dw.Truncate("table"))
}

func (dws *dialectWrapperSuite) TestDB() {
	mDB, _, err := sqlmock.New()
	dws.Require().NoError(err)
	dw := depiq.Dialect("test")
	dws.Equal(depiq.New("test", mDB), dw.DB(mDB))
}

func TestDialectWrapper(t *testing.T) {
	suite.Run(t, new(dialectWrapperSuite))
}
