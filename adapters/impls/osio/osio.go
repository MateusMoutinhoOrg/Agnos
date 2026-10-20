package osio

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	iodeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/iodeps"
)

// Bind fills deps.Deps.IoDeps with the filesystem implementation built on
// the standard library's os and filepath packages.
func Bind(deps *deps.Deps) {
	deps.IoDeps = iodeps.Contract{
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
		Exists: func(path string) bool {
			_, err := os.Stat(path)
			return err == nil || !os.IsNotExist(err)
		},
		CreateDir: func(path string) {
			_ = os.MkdirAll(path, 0755)
		},
		RemoveDir: func(path string) {
			_ = os.RemoveAll(path)
		},
		Rename: func(from string, to string) error {
			if _, err := os.Lstat(to); err == nil {
				return fmt.Errorf("cannot move %s to %s: the destination already exists", from, to)
			}
			if err := os.MkdirAll(filepath.Dir(to), 0755); err != nil {
				return err
			}
			return os.Rename(from, to)
		},
		Copy: copyTree,
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
		Join: func(elements ...string) string {
			return filepath.Join(elements...)
		},
		Dir: func(path string) string {
			return filepath.Dir(path)
		},
		UserHomeDir: func() (string, error) {
			return os.UserHomeDir()
		},
	}
}

// copyTree copies the file or directory at from to to, every file of a
// directory at any depth, each one keeping its permissions — a store that
// keeps its files private keeps them private in the copy.
func copyTree(from string, to string) error {
	return filepath.WalkDir(from, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(from, path)
		if err != nil {
			return err
		}
		target := filepath.Join(to, relative)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		return copyFile(path, target, info.Mode().Perm())
	})
}

// copyFile copies one regular file, creating the directory it lands in.
func copyFile(from string, to string, mode fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(to), 0755); err != nil {
		return err
	}
	source, err := os.Open(from)
	if err != nil {
		return err
	}
	defer source.Close()
	target, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(target, source); err != nil {
		target.Close()
		return err
	}
	return target.Close()
}
