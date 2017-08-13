package config

import (
	"testing"

	"github.com/gotestyourself/gotestyourself/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSorted(t *testing.T) {
	config := NewConfig()
	config.Resources = map[string]Resource{
		"beta":  &ImageConfig{},
		"alpha": &ImageConfig{},
		"cabo":  &ImageConfig{},
	}
	sorted := config.Sorted()
	assert.Equal(t, []string{"alpha", "beta", "cabo"}, sorted)
}

func TestResourceResolveDoesNotMutate(t *testing.T) {
	resolver := &fakeResolver{}

	for name, fromConfigFunc := range resourceTypeRegistry {
		value := make(map[string]interface{})
		resource, err := fromConfigFunc(name, value)
		assert.Nil(t, err)
		resolved, err := resource.Resolve(resolver)
		assert.Nil(t, err)
		assert.True(t, resource != resolved,
			"Expected different pointers for %q: %p, %p",
			name, resource, resolved)
	}
}

type fakeResolver struct {
}

func (r *fakeResolver) Resolve(tmpl string) (string, error) {
	return tmpl, nil
}

func (r *fakeResolver) ResolveSlice(tmpls []string) ([]string, error) {
	return tmpls, nil
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

	config, err := Load(dir.Join("dobi.yaml"))
	require.NoError(t, err)
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
	}
	assert.Equal(t, expected, config)
}
