package config

// describebale can be used to provide part of the Resource interface
type describable struct {
	// Description Description of the resource. Adding a description to a
	// resource makes it visible from ``dobi list``.
	Description string
}

func (d *describable) Describe() string {
	return d.Description
}
