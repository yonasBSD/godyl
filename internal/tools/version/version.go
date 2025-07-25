// Package version provides functionality for managing tool version information.
package version

import (
	"github.com/goccy/go-yaml/ast"

	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Version defines the target version of a tool, as well as how it can be parsed.
type Version struct {
	Commands *Commands `mapstructure:"commands" yaml:"commands"`
	Patterns *Patterns `mapstructure:"patterns" yaml:"patterns"`
	Version  string    `mapstructure:"version"  yaml:"version"  single:"true"`
}

type (
	Patterns = unmarshal.SingleOrSliceType[string]
	Commands = unmarshal.SingleOrSliceType[string]
)

func (v *Version) UnmarshalYAML(node ast.Node) error {
	type raw Version

	return unmarshal.SingleStringOrStruct(node, (*raw)(v))
}
