package smartio

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/utils/parsables/ignorableconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/utils/parsables/pathreplacerconf"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"
)

type SmartIO struct {
	Ignore       *ignorableconf.IgnorableConf
	Replacers    *pathreplacerconf.PathReplacerConf
	Transactions map[string][]byte

	ReadFile             func(path string) ([]byte, error)
	WriteFile            func(path string, content []byte) error
	WriteFileOverwrite   func(path string, content []byte) error
	Persist              func() error
	IsDir                func(path string) bool
	IsFile               func(path string) bool
	Exist                func(path string) bool
	CreateDir            func(path string)
	ListDirs             func(path string) []string
	ListFiles            func(path string) []string
	ListAll              func(path string) []string
	ListDirsRecursively  func(path string) []string
	ListFilesRecursively func(path string) []string
	ListAllRecursively   func(path string) []string
}

func joinPath(base string, name string) string {
	if strings.HasSuffix(base, "/") || strings.HasSuffix(base, "\\") {
		return base + name
	}
	// default to slash, as iodeps might be slash-based or adapter handles it
	return base + "/" + name
}

func addSmartIOMethods(sandbox *sandbox.SandBox, io *SmartIO) {

	// Helper to apply replacer and check ignore for input paths
	processInputPath := func(path string) (string, error) {
		p := io.Replacers.Format(path)
		if io.Ignore.IsIgnorable(p) {
			return p, sandbox.Deps.Errorf("path %q is ignorable", p)
		}
		return p, nil
	}

	io.ReadFile = func(path string) ([]byte, error) {
		p, err := processInputPath(path)
		if err != nil {
			return nil, err
		}
		return sandbox.Deps.IoLib.ReadFile(p)
	}

	io.WriteFile = func(path string, content []byte) error {
		p, err := processInputPath(path)
		if err != nil {
			return err
		}
		if sandbox.Deps.IoLib.Exist(p) {
			return sandbox.Deps.Errorf("file %q already exists", p)
		}
		io.Transactions[p] = content
		return nil
	}

	io.WriteFileOverwrite = func(path string, content []byte) error {
		p, err := processInputPath(path)
		if err != nil {
			return err
		}
		io.Transactions[p] = content
		return nil
	}

	io.Persist = func() error {
		for p, content := range io.Transactions {
			err := sandbox.Deps.IoLib.WriteFile(p, content)
			if err != nil {
				return err
			}
		}
		// Clear transactions after successful persist? The spec doesn't say, but it's common.
		io.Transactions = make(map[string][]byte)
		return nil
	}

	io.IsDir = func(path string) bool {
		p, err := processInputPath(path)
		if err != nil {
			return false
		}
		return sandbox.Deps.IoLib.IsDir(p)
	}

	io.IsFile = func(path string) bool {
		p, err := processInputPath(path)
		if err != nil {
			return false
		}
		return sandbox.Deps.IoLib.IsFile(p)
	}

	io.Exist = func(path string) bool {
		p, err := processInputPath(path)
		if err != nil {
			return false
		}
		return sandbox.Deps.IoLib.Exist(p)
	}

	io.CreateDir = func(path string) {
		p, err := processInputPath(path)
		if err != nil {
			return
		}
		sandbox.Deps.IoLib.CreateDir(p)
	}

	filterIgnored := func(paths []string) []string {
		var result []string
		for _, p := range paths {
			if !io.Ignore.IsIgnorable(p) {
				result = append(result, p)
			}
		}
		return result
	}

	io.ListDirs = func(path string) []string {
		p, err := processInputPath(path)
		if err != nil {
			return nil
		}
		return filterIgnored(sandbox.Deps.IoLib.ListDirs(p))
	}

	io.ListFiles = func(path string) []string {
		p, err := processInputPath(path)
		if err != nil {
			return nil
		}
		return filterIgnored(sandbox.Deps.IoLib.ListFiles(p))
	}

	io.ListAll = func(path string) []string {
		p, err := processInputPath(path)
		if err != nil {
			return nil
		}
		return filterIgnored(sandbox.Deps.IoLib.ListAll(p))
	}

	io.ListDirsRecursively = func(path string) []string {
		p, err := processInputPath(path)
		if err != nil {
			return nil
		}
		return filterIgnored(sandbox.Deps.IoLib.ListDirsRecursively(p))
	}

	io.ListFilesRecursively = func(path string) []string {
		p, err := processInputPath(path)
		if err != nil {
			return nil
		}
		return filterIgnored(sandbox.Deps.IoLib.ListFilesRecursively(p))
	}

	io.ListAllRecursively = func(path string) []string {
		p, err := processInputPath(path)
		if err != nil {
			return nil
		}
		return filterIgnored(sandbox.Deps.IoLib.ListAllRecursively(p))
	}
}

func NewSmartIO(sandbox *sandbox.SandBox, path string) *SmartIO {
	io := &SmartIO{
		Transactions: make(map[string][]byte),
	}

	configDir := joinPath(path, sandbox.Config.ProjectName+"Config")

	ignorePath := joinPath(configDir, "ignore.yaml")
	if sandbox.Deps.IoLib.Exist(ignorePath) && sandbox.Deps.IoLib.IsFile(ignorePath) {
		content, err := sandbox.Deps.IoLib.ReadFile(ignorePath)
		if err == nil {
			conf, err := ignorableconf.NewIgnorableConf(sandbox, string(content))
			if err == nil {
				io.Ignore = conf
			} else {
				io.Ignore = ignorableconf.NewIgnorableConfEmpty(sandbox)
			}
		} else {
			io.Ignore = ignorableconf.NewIgnorableConfEmpty(sandbox)
		}
	} else {
		io.Ignore = ignorableconf.NewIgnorableConfEmpty(sandbox)
	}

	replacersPath := joinPath(configDir, "paths.yaml")
	if sandbox.Deps.IoLib.Exist(replacersPath) && sandbox.Deps.IoLib.IsFile(replacersPath) {
		content, err := sandbox.Deps.IoLib.ReadFile(replacersPath)
		if err == nil {
			conf, err := pathreplacerconf.NewPathReplacerConf(sandbox, string(content))
			if err == nil {
				io.Replacers = conf
			} else {
				io.Replacers = pathreplacerconf.NewPathReplacerConfEmpty(sandbox)
			}
		} else {
			io.Replacers = pathreplacerconf.NewPathReplacerConfEmpty(sandbox)
		}
	} else {
		io.Replacers = pathreplacerconf.NewPathReplacerConfEmpty(sandbox)
	}

	addSmartIOMethods(sandbox, io)
	return io
}
