package parsables

import (
	"path"
	"strings"
)

type IgnorableItens struct {
	paths []string

	AddPath    func(path string)
	IsIgnorable func(path string) bool
}

func NewIgnorableItens(content []string) (*IgnorableItens, error) {

	ignorable := &IgnorableItens{
		paths: make([]string, 0),
	}

	for _, p := range content {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			ignorable.paths = append(ignorable.paths, trimmed)
		}
	}

	ignorable.AddPath = func(p string) {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			ignorable.paths = append(ignorable.paths, trimmed)
		}
	}

	ignorable.IsIgnorable = func(p string) bool {
		for _, pattern := range ignorable.paths {
			// exact glob match (handles *, ?, [...])
			matched, err := path.Match(pattern, p)
			if err == nil && matched {
				return true
			}

			// directory wildcard: "dir/*" should also match "dir/sub/file"
			if strings.HasSuffix(pattern, "/*") {
				prefix := strings.TrimSuffix(pattern, "/*")
				if strings.HasPrefix(p, prefix+"/") {
					return true
				}
			}
		}
		return false
	}

	return ignorable, nil
}
