//go:build !darwin

package checks

// All returns no checks on unsupported platforms (Windows/Linux support coming in v0.3).
func All() []Check {
	return nil
}
