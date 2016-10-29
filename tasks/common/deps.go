package common

// NoDependencies returns an empty list, used with Tasks that have no
// dependencies.
func NoDependencies() []string {
	return []string{}
}
