// Package ports holds the interfaces of the core.
//
// An adapter implements each port. The core owns the interfaces.
// Release v0.2.0 fills the package.
package ports

// Dialect quotes an identifier and maps a value to a type.
type Dialect interface {
	// Quote returns the quoted form of an identifier.
	Quote(name string) string

	// MapValue returns the engine form of a cell value for one field.
	MapValue(field string, value string) (string, error)
}
