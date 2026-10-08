// Package services holds the use cases of the core.
//
// A service calls the ports only. Release v0.2.0 fills the package.
package services

// CompareVersions returns the newer version of two positive versions.
// It returns zero when both versions are zero or negative.
func CompareVersions(a, b int64) int64 {
	if a <= 0 || b <= 0 {
		return 0
	}
	if a > b {
		return a
	}
	return b
}
