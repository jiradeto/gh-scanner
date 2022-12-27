package fileutils

import (
	"os"
	"path/filepath"
	"strings"
)

func ListFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

// get actual file name without preceeding root directory path
func GetBasePath(fullPath string, rootPath string) string {
	return strings.Replace(fullPath, rootPath+"/", "", -1)
}
