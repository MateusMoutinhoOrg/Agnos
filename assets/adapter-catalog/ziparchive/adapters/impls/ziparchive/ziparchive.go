package ziparchive

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"

	archivedeps "{{.Module}}/sandbox/deps/archivedeps"

	"{{.Module}}/sandbox/deps"
)

// pack fills archivedeps.Contract.Zip: every file deflated, in the order
// given, with no modification time.
func pack(files []archivedeps.File) ([]byte, error) {
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	for _, file := range files {
		entry, err := writer.CreateHeader(&zip.FileHeader{Name: file.Path, Method: zip.Deflate})
		if err != nil {
			return nil, err
		}
		if _, err := entry.Write(file.Content); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// unpack fills archivedeps.Contract.Unzip: every file of archive, refusing
// one whose files add up to more than maxBytes uncompressed. The sizes the
// archive declares are never trusted: each entry is read through a limit of
// what is left, so a lying header fails as soon as it is passed.
func unpack(archive []byte, maxBytes int64) ([]archivedeps.File, error) {
	reader, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	if err != nil {
		return nil, err
	}
	files := []archivedeps.File{}
	left := maxBytes
	for _, entry := range reader.File {
		if entry.FileInfo().IsDir() {
			continue
		}
		opened, err := entry.Open()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", entry.Name, err)
		}
		content, err := io.ReadAll(io.LimitReader(opened, left+1))
		opened.Close()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", entry.Name, err)
		}
		left -= int64(len(content))
		if left < 0 {
			return nil, fmt.Errorf("the archive unpacks to more than %d bytes", maxBytes)
		}
		files = append(files, archivedeps.File{Path: entry.Name, Content: content})
	}
	return files, nil
}

// Bind fills deps.Deps.ArchiveDeps with the standard library's archive/zip.
func Bind(deps *deps.Deps) {
	deps.ArchiveDeps = archivedeps.Contract{
		Zip: func(files []archivedeps.File) ([]byte, error) {
			return pack(files)
		},
		Unzip: func(archive []byte, maxBytes int64) ([]archivedeps.File, error) {
			return unpack(archive, maxBytes)
		},
	}
}
