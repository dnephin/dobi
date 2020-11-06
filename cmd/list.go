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
	all    bool
	groups bool
	tags   []string
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
	flags.BoolVarP(
		&listOpts.groups, "groups", "g", false,
		"List resources sorted by their matching tags")
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

	tags := getTags(conf.Resources)
	var descriptions []string
	if listOpts.groups {
		resources := filterResourcesTags(conf, listOpts)
		descriptions = getDescriptionsByTag(resources)
	} else {
		resources := filterResources(conf, listOpts)
		descriptions = getDescriptions(resources)
	}

	if len(descriptions) == 0 {
		logging.Log.Warn("No resources found. Try --all or --tags.")
		return nil
	}
	fmt.Print(format(descriptions, tags))
	return nil
}

func filterResourcesTags(conf *config.Config, listOpts listOptions) []resourceGroup {
	tags := []resourceGroup{}
	if listOpts.all {
		tags = append(tags, resourceGroup{tag: "none"})
	}
	for _, name := range conf.Sorted() {
		res := conf.Resources[name]
		if len(res.CategoryTags()) > 0 {
			for _, tagname := range res.CategoryTags() {
				currentGroupIndex := 0
				if i, found := findGroup(tags, tagname); found {
					currentGroupIndex = i
				} else {
					tags = append(tags, resourceGroup{
						tag: tagname,
					})
					currentGroupIndex = len(tags) - 1
				}
				tags[currentGroupIndex].resources = append(tags[currentGroupIndex].resources, namedResource{name: name, resource: res})
			}
		} else {
			if listOpts.all {
				tags[0].resources = append(tags[0].resources, namedResource{name: name, resource: res})
			}
		}
	}

	return tags
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

type resourceGroup struct {
	tag       string
	resources []namedResource
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

func getDescriptionsByTag(resources []resourceGroup) []string {
	lines := []string{}
	for _, tag := range resources {
		descriptions := getDescriptions(tag.resources)
		lines = append(lines, formatTags(tag.tag, descriptions))
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

func findGroup(slice []resourceGroup, tag string) (int, bool) {
	for i, item := range slice {
		if item.tag == tag {
			return i, true
		}
	}
	return -1, false
}

func format(descriptions []string, tags []string) string {
	resources := strings.Join(descriptions, "\n  ")

	msg := fmt.Sprintf("Resources:\n  %s\n", resources)
	if len(tags) > 0 {
		msg += fmt.Sprintf("\nTags:\n  %s\n", strings.Join(tags, ", "))
	}
	return msg
}

func formatTags(tag string, descriptions []string) string {
	msg := fmt.Sprintf("Tag: %s\n", tag)
	resources := strings.Join(descriptions, "\n  ")
	msg += fmt.Sprintf("  %s\n", resources)
	return msg
}
