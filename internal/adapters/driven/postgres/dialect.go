package postgres

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/kibtab/kibtab/internal/core/ports"
)

// Dialect quotes an identifier and maps a value to a text type for PostgreSQL.
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
