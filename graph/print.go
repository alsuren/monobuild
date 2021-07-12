package graph

import (
	"fmt"
	"sort"
	"strings"

	"github.com/alsuren/monobuild/set"
)

// Weak shows as a dashed line
var Weak = 1

// Strong shows as solid line
var Strong = 2

// Text returns graph as text suitable for output
func (g Graph) Text(selection []string, showType bool) string {
	var result string
	filter := set.New(selection)

	cs := make([]string, len(g.edges))
	for c := range g.edges {
		cs = append(cs, c)
	}
	sort.Strings(cs)

	for _, c := range cs {
		if !filter.Has(c) {
			continue
		}

		d := g.edges[c]

		names := make([]string, 0, len(d))
		for _, v := range d {
			if !filter.Has(v.Label) {
				continue
			}

			label := v.Label
			if showType && v.Colour == Strong {
				label = "!" + v.Label
			}

			names = append(names, label)
		}

		result += fmt.Sprintf("%s: %s\n", c, strings.Join(names, ", "))
	}

	return result
}

// Dot returns a simple text representation of the graph in the DOT language
func (g Graph) Dot(selection []string) string {
	result := fmt.Sprintln("digraph dependencies {")
	filter := set.New(selection)

	cs := make([]string, len(g.edges))
	for c := range g.edges {
		cs = append(cs, c)
	}
	sort.Strings(cs)

	for _, c := range cs {
		if !filter.Has(c) {
			continue
		}

		deps := g.edges[c]

		noDeps := true
		for _, d := range deps {
			if filter.Has(d.Label) {
				noDeps = false
				break
			}
		}

		if noDeps {
			result += fmt.Sprintf("  \"%s\"\n", c)
		}

		for _, d := range deps {
			if !filter.Has(d.Label) {
				continue
			}

			var format string
			if d.Colour == Weak { // FIXME clean up when --full is supported? Breaking change...
				format = " [style=dashed]"
			}

			result += fmt.Sprintf("  \"%s\" -> \"%s\"%s\n", c, d.Label, format)
		}
	}

	return result + "}\n"
}

// DotSchedule returns a text representation of the graph in the DOT language
// formatted as a schedule
func (g Graph) DotSchedule(selection []string) string {
	result := fmt.Sprintln("digraph schedule {\n  rankdir=\"LR\"\n  node [shape=box]")
	filter := set.New(selection)

	for c, deps := range g.edges {
		if !filter.Has(c) {
			continue
		}

		if len(deps) < 1 {
			result += fmt.Sprintf("  \"%s\"\n", c)
		}

		for _, d := range deps {
			if !filter.Has(d.Label) {
				continue
			}

			// reverse the graph during print
			result += fmt.Sprintf("  \"%s\" -> \"%s\"\n", d.Label, c)
		}
	}

	return result + "}\n"
}

// Text returns graph as text suitable for output
func (g Graph) GithubMatrix(selection []string) string {
	names := make([]string, 0)
	filter := set.New(selection)

	cs := make([]string, len(g.edges))
	for c := range g.edges {
		cs = append(cs, c)
	}
	sort.Strings(cs)

	for _, c := range cs {
		if !filter.Has(c) {
			continue
		}

		d := g.edges[c]

		for _, v := range d {
			if !filter.Has(v.Label) {
				continue
			}

			panic("TODO: work out how to deal with sub-dependencies. I'm not sure how meaningful it is to encode sub-dependencies in a build matrix, unless you force everything to run serially")
		}

		names = append(names, c)
	}

	return fmt.Sprintf("[%s]\n", strings.Join(names, ", "))
}
