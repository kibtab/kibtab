// Package ports holds the interfaces of the core.
//
// An adapter implements each port. The core owns the interfaces.
// Release v0.2.0 fills the package.
package ports

import "time"

// Clock gives the current time to the core.
type Clock interface {
	Now() time.Time
}
