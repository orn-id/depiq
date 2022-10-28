package postgres

import (
	"github.com/orn-id/depiq"
)

func DialectOptions() *depiq.SQLDialectOptions {
	do := depiq.DefaultDialectOptions()
	do.PlaceHolderFragment = []byte("$")
	do.IncludePlaceholderNum = true
	return do
}

func init() {
	depiq.RegisterDialect("postgres", DialectOptions())
}
