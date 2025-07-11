package tmp

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// ConfigFile returns the first existing “godyl” configuration file it finds.
//
// Search order:
//  1. $CONFIG_DIR/godyl.yaml
//  2. $CONFIG_DIR/godyl.yml
//  3. ./godyl.yml (project root)
//
// If none exist, it falls back to $CONFIG_DIR/godyl.yml, allowing callers
// to create or write it later.
func ConfigFile() file.File {
	base := ConfigDir().WithFile("godyl")

	// Check config directory for YAML/YML variants.
	for _, ext := range []string{"yaml", "yml"} {
		if f := base.WithExtension(ext); f.Exists() {
			return f
		}
	}

	// Fallback: project-local godyl.yml
	if local := file.New("godyl.yml"); local.Exists() {
		return local
	}

	// Default path when nothing exists yet.
	return base.WithExtension("yml")
}

// ConfigDir returns the config directory for Godyl.
func ConfigDir() folder.Folder {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return folder.New(".")
	}

	return folder.New(configDir, "godyl")
}

// CacheDir returns the cache directory for Godyl.
func CacheDir() folder.Folder {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return folder.New(os.TempDir(), "godyl")
	}

	return folder.New(cacheDir, "godyl")
}

// CacheFile returns the cache file from the specified folder.
func CacheFile(folder folder.Folder) file.File {
	return folder.WithFile("godyl.json")
}

// DownloadDir returns the download directory for Godyl.
func DownloadDir() folder.Folder {
	downloadDir, err := UserRuntimeDir()
	if err != nil {
		return folder.New(os.TempDir(), "godyl")
	}

	return folder.New(downloadDir, "godyl")
}

// GodylDir returns the temporary directory for Godyl.
// Optionally pass in subdirectories to create a path within the Godyl directory.
func GodylDir(paths ...string) folder.Folder {
	return DownloadDir().Join(paths...)
}

// GodylCreateRandomDir creates a random directory in the Godyl temporary directory.
func GodylCreateRandomDir() (folder.Folder, error) {
	// Create a random temporary directory for Godyl
	return folder.CreateRandomInDir(GodylDir().Path(), Prefix())
}

// GodylCreateRandomDirIn creates a random directory in the specified directory.
func GodylCreateRandomDirIn(dir folder.Folder) (folder.Folder, error) {
	// Create a random temporary directory for Godyl
	return folder.CreateRandomInDir(dir.Path(), Prefix())
}

// Prefix returns the prefix used for Godyl temporary directories.
func Prefix() string {
	return "godyl-*"
}

func UserRuntimeDir() (string, error) {
	var dir string

	switch runtime.GOOS {
	case "windows":
		// Windows typically uses %TEMP% for runtime/temporary files
		dir = os.Getenv("TEMP")
		if dir == "" {
			// Fall back to LocalAppData\Temp if TEMP isn't set
			appData := os.Getenv("LocalAppData")
			if appData == "" {
				return "", errors.New("%LocalAppData% is not defined")
			}

			dir = appData + "\\Temp"
		}

	case "darwin", "ios":
		// macOS typically uses /private/var/folders/XX/XXXXXXXX/T/ for temporary files
		// But for simplicity, we'll follow the tmpdir pattern:
		dir = os.Getenv("TMPDIR")
		if dir == "" {
			// Default macOS temporary directory
			dir = "/private/tmp"
		}

	case "plan9":
		dir = os.Getenv("home")
		if dir == "" {
			return "", errors.New("$home is not defined")
		}

		dir += "/tmp"

	default: // Unix
		// On Linux, XDG_RUNTIME_DIR is the standard location
		dir = os.Getenv("XDG_RUNTIME_DIR")
		if dir == "" {
			// If XDG_RUNTIME_DIR is not set, fall back to /tmp with user-specific suffix
			user := os.Getenv("USER")
			if user == "" {
				// If USER isn't available, we can only use the generic tmp
				dir = "/tmp"
			} else {
				dir = "/tmp/" + user
			}
		} else if !filepath.IsAbs(dir) {
			return "", errors.New("path in $XDG_RUNTIME_DIR is relative")
		}
	}

	return dir, nil
}
