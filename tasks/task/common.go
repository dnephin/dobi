package task

// NoDependencies returns an empty list, used with Tasks that have no
// dependencies.
func NoDependencies() []Name {
	return []Name{}
}
