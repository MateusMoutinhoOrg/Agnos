package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/extensionsconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AssetGroup is one renderable directory of assets/. Name is both the
// directory under assets/ and the name the group is rendered by; Requires is
// the set of extensions that must all be enabled for it to render.
//
// A mechanic's own group is named after its extension and requires only that
// one (assets/sandbox-cli renders when sandbox-cli is on). A doc group is
// named doc-<a>-<b> and requires doc plus each sandbox-<x> it names, because a
// layer's pages only exist when both the docs and that layer do.
type AssetGroup struct {
	Name     string
	Requires []string

	// Code marks a group carrying Go that the build's own collectors read back
	// off disk — the contracts of sandbox/api/ above all. Such a group is
	// rendered once into the transaction that turns its mechanic on, so the
	// build that follows collects against the finished tree instead of the one
	// from before the mechanic existed.
	Code bool
}

// AssetGroups is every group a build may render, in render order. assets/start
// is not here: it is written once by `start` and is not a mechanic.
func AssetGroups() []AssetGroup {
	return []AssetGroup{
		{ExtensionSandbox, []string{ExtensionSandbox}, true},
		{ExtensionSandboxDeps, []string{ExtensionSandboxDeps}, true},
		{ExtensionSandboxCli, []string{ExtensionSandboxCli}, true},
		{ExtensionSandboxServer, []string{ExtensionSandboxServer}, true},
		{ExtensionSandboxFront, []string{ExtensionSandboxFront}, true},
		{ExtensionSandboxDatabase, []string{ExtensionSandboxDatabase}, true},
		{ExtensionDoc, []string{ExtensionDoc}, false},
		{"doc-cli", []string{ExtensionDoc, ExtensionSandboxCli}, false},
		{"doc-server", []string{ExtensionDoc, ExtensionSandboxServer}, false},
		{"doc-front", []string{ExtensionDoc, ExtensionSandboxFront}, false},
		{"doc-database", []string{ExtensionDoc, ExtensionSandboxDatabase}, false},
		{"doc-example", []string{ExtensionDoc, ExtensionSandboxExample}, false},
		{"doc-example-cli", []string{ExtensionDoc, ExtensionSandboxExample, ExtensionSandboxCli}, false},
		{ExtensionReadme, []string{ExtensionReadme}, false},
	}
}

// GroupEnabled reports whether every extension a group requires is enabled.
func GroupEnabled(conf *extensionsconf.ExtensionsConf, group AssetGroup) bool {
	for _, required := range group.Requires {
		if !conf.IsEnabled(required) {
			return false
		}
	}
	return true
}

// RenderableGroups is the names of the groups this declaration turns on, in
// render order. It is the whole of what decides which assets a build writes.
func RenderableGroups(conf *extensionsconf.ExtensionsConf) []string {
	var groups []string
	for _, group := range AssetGroups() {
		if GroupEnabled(conf, group) {
			groups = append(groups, group.Name)
		}
	}
	return groups
}

// DocGroups is RenderableGroups narrowed to the groups that carry docs, which
// is what CollectGeneratedDocs walks to index the pages this build writes.
func DocGroups(sandbox *api.Sandbox, conf *extensionsconf.ExtensionsConf) []string {
	var groups []string
	for _, group := range AssetGroups() {
		if !GroupEnabled(conf, group) {
			continue
		}
		if group.Name == ExtensionDoc || sandbox.Deps.Stringsdeps.HasPrefix(group.Name, ExtensionDoc+"-") {
			groups = append(groups, group.Name)
		}
	}
	return groups
}

// GroupsRequiring is every group that names one extension among its
// requirements — the set of assets that mechanic owns, which is what an
// <x>-purge removes. A doc group belongs to the layer it documents as well as
// to doc, so purging a layer takes its pages with it.
func GroupsRequiring(name string) []string {
	var groups []string
	for _, group := range AssetGroups() {
		for _, required := range group.Requires {
			if required == name {
				groups = append(groups, group.Name)
				break
			}
		}
	}
	return groups
}

// ExtensionFiles is every group-relative path the groups one mechanic owns
// would install. It is what an <x>-purge removes: the asset groups only name
// the files they write, so this is the exact inverse of what a build renders
// for that mechanic.
func ExtensionFiles(sandbox *api.Sandbox, name string) ([]string, error) {
	var paths []string

	for _, group := range GroupsRequiring(name) {
		files, err := sandbox.Deps.Embeddeps.ListFilesRecursively(group)
		if err != nil {
			return nil, err
		}
		paths = append(paths, files...)
	}

	return paths, nil
}

// RenderExtensionCode renders the code group one mechanic owns into an already
// open transaction, ahead of the build that follows. Everything else the
// mechanic writes — its pages above all — is left to that build, which has the
// collector output those templates need.
//
// A mechanic with no code group renders nothing here.
func RenderExtensionCode(sandbox *api.Sandbox, io *smartio.SmartIO, name string) error {
	for _, group := range AssetGroups() {
		if group.Name != name || !group.Code {
			continue
		}

		module_conf, err := LoadModuleConf(sandbox, io)
		if err != nil {
			return err
		}

		return RenderGroup(sandbox, io, group.Name, map[string]interface{}{
			"Module": module_conf.Module,
			// The front group's pageio names the mount its links are built
			// against. The build that follows renders it again off the same
			// declaration, so the value is not new here — leaving it out is,
			// and an unfilled var renders a verbatim `%!q(<nil>)` into a Go
			// file the format pass then refuses.
			"StaticMount": CollectFrontMount(sandbox, io),
		})
	}

	return nil
}
