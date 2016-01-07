package config

var (
	volumeKeys  = []string{"path", "mount"}
	imageKeys   = []string{"image"}
	commandKeys = []string{"use"}
)

type stringKeyMap struct {
	value    map[string]interface{}
	resource Resource
}

type rawMap struct {
	value map[string]stringKeyMap
}

func newRawMap() *rawMap {
	return &rawMap{
		value: make(map[string]stringKeyMap),
	}
}

// UnmarshalYAML unmarshals a raw config resource
func (m *stringKeyMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	m.value = make(map[string]interface{})
	err := unmarshal(m.value)
	if err != nil {
		return err
	}

	var conf Resource
	switch {
	case m.hasKeys(volumeKeys):
		conf = NewVolumeConfig()
	case m.hasKeys(imageKeys):
		conf = NewImageConfig()
	case m.hasKeys(commandKeys):
		conf = &CommandConfig{}
	default:
		// TODO: error on unknown resource type
	}

	if err = conf.Validate(); err != nil {
		return err
	}

	// TODO: error on unexpected fields
	err = unmarshal(conf)
	m.resource = conf
	return err
}

func (m *stringKeyMap) hasKeys(keys []string) bool {
	for _, key := range keys {
		if _, ok := m.value[key]; ok {
			return true
		}
	}
	return false
}
