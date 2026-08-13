package standard

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/iodeps"
)

// IoLibFactory returns the value that fills deps.Deps.IoLib: the implementation
// of the Io dependency using the standard library's os and filepath packages.
func IoLibFactory(s *StandardAdapter) iodeps.Lib {
	return iodeps.Lib{
		ReadFile: func(path string) ([]byte, error) {
			return os.ReadFile(path)
		},
		WriteFile: func(path string, content []byte) error {
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}
			return os.WriteFile(path, content, 0644)
		},
		IsDir: func(path string) bool {
			info, err := os.Stat(path)
			return err == nil && info.IsDir()
		},
		IsFile: func(path string) bool {
			info, err := os.Stat(path)
			return err == nil && !info.IsDir()
		},
		Exist: func(path string) bool {
			_, err := os.Stat(path)
			return err == nil || !os.IsNotExist(err)
		},
		CreateDir: func(path string) {
			_ = os.MkdirAll(path, 0755)
		},
		ListDirs: func(path string) []string {
			var dirs []string
			entries, err := os.ReadDir(path)
			if err == nil {
				for _, entry := range entries {
					if entry.IsDir() {
						dirs = append(dirs, filepath.Join(path, entry.Name()))
					}
				}
			}
			return dirs
		},
		ListFiles: func(path string) []string {
			var files []string
			entries, err := os.ReadDir(path)
			if err == nil {
				for _, entry := range entries {
					if !entry.IsDir() {
						files = append(files, filepath.Join(path, entry.Name()))
					}
				}
			}
			return files
		},
		ListAll: func(path string) []string {
			var all []string
			entries, err := os.ReadDir(path)
			if err == nil {
				for _, entry := range entries {
					all = append(all, filepath.Join(path, entry.Name()))
				}
			}
			return all
		},
		ListDirsRecursively: func(path string) []string {
			var dirs []string
			_ = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
				if err == nil && p != path && d.IsDir() {
					dirs = append(dirs, p)
				}
				return nil
			})
			return dirs
		},
		ListFilesRecursively: func(path string) []string {
			var files []string
			_ = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
				if err == nil && !d.IsDir() {
					files = append(files, p)
				}
				return nil
			})
			return files
		},
		ListAllRecursively: func(path string) []string {
			var all []string
			_ = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
				if err == nil && p != path {
					all = append(all, p)
				}
				return nil
			})
			return all
		},
	}
}
