package config

// Config holds the top level configuration for godyl.
// It is split into sub-structs for each command.

type Config struct {
	// Root level configuration, mapping configurations on the root `godyl` command
	Root Root

	// Tool level configuration, mapping configurations on the `install`, `download`
	Tool Tool

	// Update level configuration, mapping configurations on the `update` command
	Update Update

	// Dump level configuration, mapping configurations on the `dump` command
	Dump Dump

	// Cache level configuration, mapping configurations on the `cache` command
	Cache Cache

	// Status level configuration, mapping configurations on the `status` command
	Status Status

	// Viper instance
	viperable `json:"-" mapstructure:"-" yaml:"-"`
}
