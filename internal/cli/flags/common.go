package flags

import (
	"fmt"
	"strings"

	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Viperable is an interface for types that can hold a viper instance
type Viperable interface {
	SetViper(v *viper.Viper)
	GetViper() *viper.Viper
}

func PrefixToYAML(prefix string, root string) string {
	prefix = strings.TrimPrefix(prefix, root)
	prefix = strings.ReplaceAll(prefix, "_", ".")
	prefix = strings.TrimPrefix(prefix, ".")

	return prefix
}

// Bind connects cobra flags to viper and unmarshals the configuration into the provided struct.
// It sets up environment variable handling with the given prefix and handles flag binding.
// Omit the prefix to use the command hierarchy as the prefix.
func Bind(cmd *cobra.Command, cfg Viperable, prefix ...string) error {
	// Set up Viper with our environment prefix
	envPrefix := prefixFromCmdOrPrefixes(cmd, prefix...)

	// Reuse the same instance if already set
	if cfg.GetViper() == nil {
		cfg.SetViper(viper.New())
	}

	viper := cfg.GetViper()

	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	configFile := cmd.Root().Context().Value("config-file")
	isSet := cmd.Root().Context().Value("config-file-set")

	isConfigError := func(err error) bool {
		return err != nil && isSet != nil && isSet.(bool)
	}

	if configFile != nil {
		config := file.File(configFile.(string))
		content, err := Trim(config, PrefixToYAML(envPrefix, cmd.Root().Name()))
		if isConfigError(err) {
			return fmt.Errorf("trimming config file: %w", err)
		} else if err == nil {
			viper.SetConfigType("yaml")
			if err := viper.ReadConfig(content); isConfigError(err) {
				return fmt.Errorf("reading config file: %w", err)
			}
		}
	}

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return fmt.Errorf("binding flags: %w", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("unmarshalling config for %q: %w", cmd.Name(), err)
	}

	return nil
}

// prefixFromCmdOrPrefixes builds an environment variable prefix string either from
// a command's hierarchy or from explicitly provided prefix parts.
// When prefixes are provided, they take precedence over the command hierarchy.
func prefixFromCmdOrPrefixes(cmd *cobra.Command, prefixes ...string) string {
	if len(prefixes) > 0 {
		// Use explicitly provided prefixes if available
		return strings.Join(prefixes, "_")
	}

	// Otherwise build prefix from command hierarchy
	var commandPathParts []string
	currentCmd := cmd

	// Traverse up the command tree to build the path
	for currentCmd != nil {
		commandPathParts = append([]string{currentCmd.Name()}, commandPathParts...)
		currentCmd = currentCmd.Parent()
	}

	return strings.Join(commandPathParts, "_")
}
