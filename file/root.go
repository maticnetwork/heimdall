package file

import "path/filepath"

func Rootify(path, root string) string {
	if filepath.IsAbs(path) {
		return path
	}

	return filepath.Join(root, path)
}
