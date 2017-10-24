package context

// Settings are flags that can be set by a user to change the behaviour of some
// tasks
type Settings struct {
	Quiet     bool
	BindMount bool
}

// NewSettings returns a new Settings
func NewSettings(quiet bool, bindMount bool) Settings {
	return Settings{Quiet: quiet, BindMount: bindMount}
}
