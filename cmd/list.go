package cmd

import (
	"fmt"
	"sort"

	"strings"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/spf13/cobra"
)

type listOptions struct {
	all  bool
	tags []string
}

func (o listOptions) tagMatch(tags []string) bool {
	for _, otag := range o.tags {
		for _, tag := range tags {
			if tag == otag {
				return true
			}
		}
	}
	return false
}

func newListCommand(opts *dobiOptions) *cobra.Command {
	var listOpts listOptions
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(opts, listOpts)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(
		&listOpts.all, "all", "a", false,
		"List all resources, including those without descriptions")
	flags.StringSliceVarP(
		&listOpts.tags, "tags", "t", nil,
		"List tasks matching the tag")
	return cmd
}

func runList(opts *dobiOptions, listOpts listOptions) error {
	conf, err := config.Load(opts.filename)
	if err != nil {
		return err
	}

	resources := filterResources(conf, listOpts)
	descriptions := getDescriptions(resources)
	if len(descriptions) == 0 {
		logging.Log.Warn("No resources found. Try --all or --tags.")
		return nil
	}

	tags := getTags(conf.Resources)

	fmt.Print(format(descriptions, tags))
	return nil
}

func filterResources(conf *config.Config, listOpts listOptions) []namedResource {
	resources := []namedResource{}
	for _, name := range conf.Sorted() {
		res := conf.Resources[name]
		if include(res, listOpts) {
			resources = append(resources, namedResource{name: name, resource: res})
		}
	}
	return resources
}

type namedResource struct {
	name     string
	resource config.Resource
}

func (n namedResource) Describe() string {
	desc := n.resource.Describe()
	if desc == "" {
		return n.resource.String()
	}
	return desc
}

func include(res config.Resource, listOpts listOptions) bool {
	if listOpts.all || listOpts.tagMatch(res.CategoryTags()) {
		return true
	}
	return len(listOpts.tags) == 0 && res.Describe() != ""
}

func getDescriptions(resources []namedResource) []string {
	lines := []string{}
	for _, named := range resources {
		line := fmt.Sprintf("%-20s %s", named.name, named.Describe())
		lines = append(lines, line)
	}
	return lines
}

func getTags(resources map[string]config.Resource) []string {
	mapped := make(map[string]struct{})
	for _, res := range resources {
		for _, tag := range res.CategoryTags() {
			mapped[tag] = struct{}{}
		}
	}
	tags := []string{}
	for tag := range mapped {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

func format(descriptions []string, tags []string) string {
	resources := strings.Join(descriptions, "\n  ")

	msg := fmt.Sprintf("Resources:\n  %s\n", resources)
	if len(tags) > 0 {
		msg += fmt.Sprintf("\nTags:\n  %s\n", strings.Join(tags, ", "))
	}
	return msg
}
