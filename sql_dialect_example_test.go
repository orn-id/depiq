package depiq_test

import (
	"fmt"

	"github.com/orn-id/depiq/v9"
)

func ExampleRegisterDialect() {
	opts := depiq.DefaultDialectOptions()
	opts.QuoteRune = '`'
	depiq.RegisterDialect("custom-dialect", opts)

	dialect := depiq.Dialect("custom-dialect")

	ds := dialect.From("test")

	sql, args, _ := ds.ToSQL()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM `test` []
}
