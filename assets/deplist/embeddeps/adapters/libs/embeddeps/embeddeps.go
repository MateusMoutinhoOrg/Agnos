package embeddeps

// The standard adapter's embedded-asset implementation: the factory filling
// deps.Deps.Embeddeps out of the assets compiled into the binary. It is the
// one adapter that reaches for a package of the project outside the sandbox —
// assets — which is what the conversion helpers below exist for.
//
// Everything here is outside the sandbox, which is what makes the `embed`
// directive and the io/fs walk legal: the sandbox only ever sees the three
// function fields of embeddeps.Lib.

import (
	"bytes"
	"io/fs"
	"path"
	"text/template"

	"{{.Module}}/assets"
	"{{.Module}}/sandbox/deps"
	embeddeps "{{.Module}}/sandbox/deps/embeddeps"
)

// assetPath resolves one path the library asked for against the root of the
// asset tree, so "report.tmpl" is "report.tmpl" and "" is the root itself.
// path.Join cleans the result, which is what an embedded filesystem requires:
// slash separators, no "." element, and no leading slash.
func assetPath(requested string) string {
	return path.Join(".", requested)
}

// assetFiles lists the assets under one directory of the embedded filesystem.
// Only files are reported, in the lexical order fs.WalkDir walks in, each as
// a slash-separated path relative to the directory that was asked for.
// descend chooses between one level and the whole subtree, which is the only
// difference between ListFiles and ListFilesRecursively.
func assetFiles(root string, descend bool) ([]string, error) {
	listed := []string{}
	err := fs.WalkDir(assets.Files, root, func(current string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if current == root {
			return nil
		}
		if entry.IsDir() {
			if descend {
				return nil
			}
			return fs.SkipDir
		}
		listed = append(listed, relativeTo(root, current))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return listed, nil
}

// relativeTo reports where current sits inside root. It is written by hand
// rather than taken from path/filepath because an embedded filesystem is
// always slash-separated, whatever separator the host operating system uses,
// and because fs.WalkDir only ever hands back paths already under root.
func relativeTo(root string, current string) string {
	if root == "." {
		return current
	}
	return current[len(root)+1:]
}

// Bind fills deps.Deps.Embeddeps with the project's assets, compiled
// into the binary by the assets package and served from the root of that tree.
// Each field of the assigned struct is a closure reading the embedded
// filesystem at call time.
func Bind(deps *deps.Deps) {
	deps.Embeddeps = embeddeps.Lib{
		ReadFile: func(requested string) ([]byte, error) {
			return assets.Files.ReadFile(assetPath(requested))
		},
		ListFiles: func(requested string) ([]string, error) {
			return assetFiles(assetPath(requested), false)
		},
		ListFilesRecursively: func(requested string) ([]string, error) {
			return assetFiles(assetPath(requested), true)
		},
		RenderTemplate: func(requested string, vars interface{}) ([]byte, error) {
			content, err := assets.Files.ReadFile(assetPath(requested))
			if err != nil {
				return nil, err
			}
			t, err := template.New(path.Base(requested)).Parse(string(content))
			if err != nil {
				return nil, err
			}
			var buf bytes.Buffer
			if err := t.Execute(&buf, vars); err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		},
	}
}
