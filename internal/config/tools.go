package config

import (
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/pkg/file"
)

type Tool struct {
	// Output path for the downloaded tools
	Output string

	// Tags to filter tools by
	Tags []string

	// Source from which to install the tools
	Source sources.Type `validate:"oneof=github url go command"`

	// Strategy to use for updating tools
	Strategy tools.Strategy `validate:"oneof=none upgrade force"`

	// Tokens for authentication
	Tokens Tokens `mapstructure:",squash"`

	// Operating system to install the tools for
	OS string

	// Architecture to install the tools for
	Arch string

	// Path to tools configuration file
	Tools file.File // Positional argument

	// Number of parallel downloads (>= 0)
	Parallel int `validate:"gte=0"`

	// Skip SSL verification
	NoVerifySSL bool `mapstructure:"no-verify-ssl"`

	// Hints to use for tool resolution
	Hints []string
}

// Tokens holds the configuration options for authentication tokens.
type Tokens struct {
	// GitHub token for authentication
	GitHub string `mapstructure:"github-token" mask:"fixed"`
}
