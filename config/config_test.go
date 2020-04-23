package config

import (
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
	"gotest.tools/v3/fs"
)

func TestSorted(t *testing.T) {
	config := NewConfig()
	config.Resources = map[string]Resource{
		"beta":  &ImageConfig{},
		"alpha": &ImageConfig{},
		"cabo":  &ImageConfig{},
	}
	sorted := config.Sorted()
	assert.Check(t, is.DeepEqual([]string{"alpha", "beta", "cabo"}, sorted))
}

func TestResourceResolveDoesNotMutate(t *testing.T) {
	resolver := newFakeResolver(nil)

	for name, fromConfigFunc := range resourceTypeRegistry {
		value := make(map[string]interface{})
		resource, err := fromConfigFunc("resourcename", value)
		assert.NilError(t, err)
		resolved, err := resource.Resolve(resolver)
		assert.NilError(t, err)
		assert.Check(t, resource != resolved,
			"Expected different pointers for %q: %p, %p",
			name, resource, resolved)
	}
}

type fakeResolver struct {
	mapping map[string]string
}

func (r *fakeResolver) Resolve(tmpl string) (string, error) {
	value, ok := r.mapping[tmpl]
	if ok {
		return value, nil
	}
	return tmpl, nil
}

func (r *fakeResolver) ResolveSlice(tmpls []string) ([]string, error) {
	values := []string{}
	for _, key := range tmpls {
		value, _ := r.Resolve(key)
		values = append(values, value)
	}
	return values, nil
}

func newFakeResolver(mapping map[string]string) *fakeResolver {
	if mapping == nil {
		mapping = make(map[string]string)
	}
	return &fakeResolver{mapping: mapping}
}

// FIXME: not a full config
func TestLoadFullFromYaml(t *testing.T) {
	dir := fs.NewDir(t, "load-full-yaml",
		fs.WithFile("dobi.yaml", `
meta:
    project: fulltest
    default: one
    exec-id: exec_id

alias=one:
    tasks: []
alias=two:
    tasks: []
alias=three:
    tasks: []

alias=aliasresource:
    tasks: [one, two, three]
    annotations:
        description: This is an alias resource
        tags: [lots, things]
`))
	defer dir.Remove()

	yamlPath := dir.Join("dobi.yaml")
	config, err := Load(yamlPath)
	assert.NilError(t, err)
	expected := &Config{
		Meta: &MetaConfig{
			Project: "fulltest",
			Default: "one",
			ExecID:  "exec_id",
		},
		Resources: map[string]Resource{
			"aliasresource": &AliasConfig{
				Tasks: []string{"one", "two", "three"},
				Annotations: Annotations{
					Annotations: AnnotationFields{
						Description: "This is an alias resource",
						Tags:        []string{"lots", "things"},
					},
				},
			},
			"one":   &AliasConfig{Tasks: []string{}},
			"two":   &AliasConfig{Tasks: []string{}},
			"three": &AliasConfig{Tasks: []string{}},
		},
		WorkingDir: dir.Path(),
		FilePath:   yamlPath,
	}
	assert.Check(t, is.DeepEqual(expected, config, cmpConfigOpt))
}
