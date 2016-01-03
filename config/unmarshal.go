package config

var (
	volumeKeys = []string{"path", "mount"}
	imageKeys  = []string{"image"}
	taskKeys   = []string{"use"}
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

	var conf interface{}
	switch {
	case m.hasKeys(volumeKeys):
		conf = &VolumeConfig{}
	case m.hasKeys(imageKeys):
		conf = &ImageConfig{}
	case m.hasKeys(taskKeys):
		conf = &TaskConfig{}
	default:
		// TODO: error on unknown resource type
	}

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
