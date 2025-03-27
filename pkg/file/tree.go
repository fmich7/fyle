package file

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// printDirIndent prints indent + entry name for files.
func printIndent(entry string, depth int) string {
	indent := strings.Repeat("|    ", depth)
	return fmt.Sprintf("%s|── %s\n", indent, entry)
}

// printDirIndent prints indent + entry name for dirs.
func printDirIndent(entry string, depth int) string {
	indent := strings.Repeat("|    ", depth)
	return fmt.Sprintf("%s|── %s/\n", indent, entry)
}

// printDir prints recursively directory to mimic tree command.
func printDir(fs afero.Fs, path, rootPath string, depth int, buf *strings.Builder) error {
	entries, err := afero.ReadDir(fs, path)
	if err != nil {
		return fmt.Errorf("reading specified path %v", err)
	}

	// print curr folder
	folder, err := filepath.Rel(rootPath, path)
	if err != nil {
		return err
	}

	_, err = buf.WriteString(printDirIndent(folder, depth))
	if err != nil {
		return err
	}

	depth++
	for _, entry := range entries {
		// recursive print dirs
		if entry.IsDir() {
			newDir := filepath.Join(path, entry.Name())
			err := printDir(fs, newDir, path, depth, buf)
			if err != nil {
				return err
			}
		} else {
			// print entry
			_, err := buf.WriteString(printIndent(entry.Name(), depth))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetDirTree is a recursive function to list directory content as a tree.
func GetDirTree(fs afero.Fs, path string) (string, error) {
	buf := new(strings.Builder)
	rootPath := path
	err := printDir(fs, path, rootPath, 0, buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
