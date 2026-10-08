package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/depconf"
)

// DepCatalogGroup is the embedded catalog of installable contracts: one directory
// per dep, holding the dep.yaml that declares it beside the sandbox/deps/<dep>/
// tree it installs.
const DepCatalogGroup = "dep-catalog"

// DepConfFile is the declaration at the root of one catalog dep directory. It
// is the only asset there that is not installed: it describes the package
// rather than belonging to it.
const DepConfFile = "dep.yaml"

// ContractsDir holds one contract package per dep, the closed side of the pair
// an adapter fills.
const ContractsDir = "sandbox/deps"

// LoadCatalogDepConf reads assets/dep-catalog/<dep>/dep.yaml out of the embedded
// catalog. A dep with no declaration is not a dep, so the error is the one
// callers report for an unknown name.
func LoadCatalogDepConf(sandbox *api.Sandbox, dep string) (*depconf.DepConf, error) {
	content, err := sandbox.Deps.EmbedDeps.ReadFile(DepCatalogGroup + "/" + dep + "/" + DepConfFile)
	if err != nil {
		return nil, sandbox.Deps.StdDeps.Errorf("unknown dep %q", dep)
	}

	return depconf.New(sandbox, string(content))
}

// CatalogDeps returns the name of every dep in the embedded catalog, in
// listing order.
func CatalogDeps(sandbox *api.Sandbox) ([]string, error) {
	return catalogEntries(sandbox, DepCatalogGroup)
}

// catalogEntries returns the immediate sub-directory names of an embedded
// catalog group, in listing order — one name per unit, because each unit of
// the catalog is exactly one directory.
func catalogEntries(sandbox *api.Sandbox, group string) ([]string, error) {
	files, err := sandbox.Deps.EmbedDeps.ListFilesRecursively(group)
	if err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	entries := []string{}
	for _, file := range files {
		name := sandbox.Deps.StringsDeps.Split(file, "/")[0]
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		entries = append(entries, name)
	}

	return entries, nil
}

// ContractType is the root type every contract under sandbox/deps/ declares,
// the one its Deps field holds; RemoteContractType is the one of a remote dep,
// whose copy keeps the name of the api it was copied from.
const (
	ContractType       = "Contract"
	RemoteContractType = "Sandbox"
)

// DepField is the Deps field a contract directory fills, the same spelling
// sandbox/deps/deps.go is generated with: the directory title-cased, with a
// trailing "deps" spelled as its own word — argvdeps fills ArgvDeps,
// OpinionatedAgnosCli fills OpinionatedAgnosCli.
func DepField(sandbox *api.Sandbox, dep string) string {
	if len(dep) == 0 {
		return ""
	}
	strs := sandbox.Deps.StringsDeps
	if len(dep) > len("deps") && strs.HasSuffix(dep, "deps") {
		stem := strs.TrimSuffix(dep, "deps")
		return strs.ToUpper(stem[:1]) + stem[1:] + "Deps"
	}
	return strs.ToUpper(dep[:1]) + dep[1:]
}

// ValidateDepName rejects a name that could not be a directory under
// sandbox/deps/ and a Go package clause at the same time — the same check
// ValidateBindingName makes, for the same reason: the name is propagated
// straight into `package <name>` of the copied contract.
func ValidateDepName(sandbox *api.Sandbox, dep string) error {
	if dep == "" {
		return sandbox.Deps.StdDeps.Errorf("a dep needs a name")
	}
	if dep[0] < 'a' || dep[0] > 'z' {
		return sandbox.Deps.StdDeps.Errorf("invalid dep name %q: a dep name must start with a lowercase letter", dep)
	}
	for _, letter := range dep {
		if (letter >= 'a' && letter <= 'z') || (letter >= '0' && letter <= '9') {
			continue
		}
		return sandbox.Deps.StdDeps.Errorf("invalid dep name %q: only lowercase letters and digits are allowed (it becomes the directory %s/%s and a Go package name)",
			dep, ContractsDir, dep)
	}
	return nil
}
