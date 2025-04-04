package file

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/idelchi/godyl/pkg/utils"
)

// File represents a file path as a string, providing methods for file operations.
type File string

// New creates a new File by joining the provided paths.
func New(paths ...string) File {
	return File(filepath.Clean(filepath.Join(paths...))) // .Normalized()
}

// Normalized converts the file path to use forward slashes.
func (f File) Normalized() File {
	return New(filepath.ToSlash(f.Path()))
}

// Create creates a new file.
func (f File) Create() error {
	file, err := os.Create(f.String())
	if err != nil {
		return fmt.Errorf("creating file %q: %w", f, err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("closing file %q: %w", f, err)
	}

	return nil
}

// OpenForWriting opens the file for writing and returns a pointer to the os.File object.
// If the file doesn't exist, it will be created.
// If it exists, it will be truncated.
// The user must close the file after use.
func (f File) OpenForWriting() (*os.File, error) {
	const perm = 0o600

	file, err := os.OpenFile(f.String(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return nil, fmt.Errorf("opening file %q for writing: %w", f, err)
	}

	return file, nil
}

// Open opens the file for reading and returns a pointer to the os.File object, or an error.
// The user must close the file after use.
func (f File) Open() (*os.File, error) {
	file, err := os.Open(f.String())
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", f, err)
	}

	return file, nil
}

// Remove deletes the file from the file system.
func (f File) Remove() error {
	if err := os.Remove(f.String()); err != nil {
		return fmt.Errorf("removing file %q: %w", f, err)
	}

	return nil
}

// Path returns the path of the File.
func (f File) Path() string {
	return f.String()
}

// String returns the string representation of the File.
func (f File) String() string {
	return string(f)
}

// Base returns the base name of the File.
func (f File) Base() string {
	return filepath.Base(f.String())
}

// Expanded expands the file path in case of ~ and returns the expanded path.
func (f File) Expanded() File {
	return New(utils.ExpandHome(f.Path()))
}

// Dir returns the directory of the file.
// If the file is a directory, it returns the path of the directory itself.
func (f File) Dir() string {
	if f.IsDir() {
		return f.Path()
	}

	return filepath.Dir(f.Path())
}

// IsExecutable checks if the file has executable permissions.
func (f File) IsExecutable() (bool, error) {
	info, err := os.Stat(f.String())
	if err != nil {
		return false, fmt.Errorf("getting file info: %w", err)
	}

	return info.Mode()&0o111 != 0, nil
}

// Chmod changes the file permissions to the specified fs.FileMode.
func (f *File) Chmod(mode fs.FileMode) error {
	return os.Chmod(f.String(), mode)
}

// Copy copies the contents of the current file to the specified destination File.
// It returns an error if the operation fails.
func (f File) Copy(other File) error {
	// Open the source file
	source, err := f.Open()
	if err != nil {
		return fmt.Errorf("opening source file: %w", err)
	}
	defer source.Close()

	// Create the destination file
	destination, err := other.OpenForWriting()
	if err != nil {
		return fmt.Errorf("opening destination file: %w", err)
	}
	defer destination.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("copying file: %w", err)
	}

	const perm = 0o755

	// Set permissions on the destination file (executable permission)
	if err := destination.Chmod(perm); err != nil {
		return fmt.Errorf("setting permissions: %w", err)
	}

	return nil
}

// Exists checks if the file exists in the file system.
func (f File) Exists() bool {
	_, err := os.Stat(f.String())

	return err == nil
}

// IsFile checks if the path is a regular file (not a directory or special file).
func (f File) IsFile() bool {
	info, err := os.Stat(f.String())
	if err != nil {
		return false // File does not exist or error accessing it
	}

	return info.Mode().IsRegular()
}

// IsDir checks if the path represents a directory.
func (f File) IsDir() bool {
	info, err := os.Stat(f.String())
	if err != nil {
		return false // File does not exist or error accessing it
	}

	return info.Mode().IsDir()
}

// Extension returns the file extension of the File, mapped to a predefined Extension constant.
func (f File) Extension() Extension {
	ext := filepath.Ext(f.Path())

	switch strings.ToLower(ext) {
	case ".exe":
		return EXE
	case ".gz":
		return GZ
	case ".zip":
		return ZIP
	case "":
		return None
	default:
		return Other
	}
}
