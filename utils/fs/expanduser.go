package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExpandUser expands the ~ in a path to the users home directory
// TODO: find a lib, or move to a lib
// TODO: os.user can't be used with cross compiling
func ExpandUser(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	parts := strings.Split(path, fmt.Sprintf("%c", filepath.Separator))
	username := strings.TrimPrefix(parts[0], "~")
	switch username {
	case "":
		parts[0] = os.Getenv("HOME")
		return filepath.Join(parts...), nil
	default:
		return path, fmt.Errorf("expanding ~user/ paths are not supported yet")
	}
}
