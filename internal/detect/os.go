//go:build windows || darwin || freebsd || openbsd || netbsd

package detect

import (
	"runtime"

	"github.com/idelchi/godyl/internal/detect/platform"
)

// Detect gathers information about the current platform, including the operating system, architecture,
// library, and file extension, and populates the Platform struct accordingly for Windows and macOS.
func (p *Platform) Detect() error {
	var os platform.OS
	var arch platform.Architecture
	var library platform.Library
	var distro platform.Distribution
	var extension platform.Extension

	// Determine the OS from runtime information
	if err := os.From(runtime.GOOS); err != nil {
		return err
	}

	// Set the default library based on the OS (distribution is irrelevant for Windows/macOS)
	library = library.Default(os, distro)

	// Determine the architecture from runtime information
	if err := arch.From(runtime.GOARCH, distro); err != nil {
		return err
	}

	// Populate the Platform struct with the detected values
	*p = Platform{
		OS:           os,
		Distribution: distro,
		Architecture: arch,
		Library:      library,
		Extension:    extension.Default(os),
	}

	return nil
}
