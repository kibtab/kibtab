package postgres

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/kibtab/kibtab/internal/core/ports"
)

// Dialect quotes an identifier, maps a value to a text type for
// PostgreSQL, and pages a query.
type Dialect struct{}

// Compile-time proof that Dialect satisfies the Dialect port.
var _ ports.Dialect = Dialect{}

// Quote returns the quoted form of an identifier. Dot-separated names are
// quoted as a schema-qualified identifier.
func (Dialect) Quote(name string) string {
	return pgx.Identifier(strings.Split(name, ".")).Sanitize()
}

// MapValue returns the PostgreSQL text form of a cell value for one field.
func (Dialect) MapValue(field string, value string) (string, error) {
	if field == "" {
		return "", errors.New("postgres: field name is empty")
	}
	if strings.ContainsRune(value, 0) {
		return "", errors.New("postgres: value contains a NUL byte")
	}
	return value, nil
}

// Page returns the PostgreSQL paging clause and its argument. PostgreSQL
// names each placeholder with a dollar sign and the placeholder number.
func (Dialect) Page(limit int, first int) (string, []any) {
	return fmt.Sprintf(" LIMIT $%d", first), []any{limit}
}
