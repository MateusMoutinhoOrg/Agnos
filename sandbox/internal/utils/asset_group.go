package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/extensionsconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// AssetGroup is one renderable directory of assets/. Name is both the
// directory under assets/ and the name the group is rendered by; Requires is
// the set of extensions that must all be enabled for it to render.
//
// A mechanic's own group is named after its extension and requires only that
// one (assets/cli renders when cli is on). A doc group is
// named doc-<a>-<b> and requires doc plus each <x> it names, because a
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
// is not here: it is written once by `start` and is not a mechanic. Neither is
// front: its code is the OpinionatedAgnosFront lib, installed by front-init, so
// it renders its pages and nothing else. The database layer's code is the
// OpinionatedAgnosDatabase lib too; what it renders is its part of api.Config
// (database) and, with a cli to read it from, the --database middleware
// (database-cli).
func AssetGroups() []AssetGroup {
	return []AssetGroup{
		{ExtensionSandbox, []string{ExtensionSandbox}, true},
		{ExtensionDeps, []string{ExtensionDeps}, true},
		{ExtensionCli, []string{ExtensionCli}, true},
		{ExtensionServer, []string{ExtensionServer}, true},
		{ExtensionDatabase, []string{ExtensionDatabase}, true},
		{"database-cli", []string{ExtensionDatabase, ExtensionCli}, true},
		{ExtensionDoc, []string{ExtensionDoc}, false},
		{"doc-cli", []string{ExtensionDoc, ExtensionCli}, false},
		{"doc-server", []string{ExtensionDoc, ExtensionServer}, false},
		{"doc-front", []string{ExtensionDoc, ExtensionFront}, false},
		{"doc-database", []string{ExtensionDoc, ExtensionDatabase}, false},
		{"doc-backoffice", []string{ExtensionDoc, ExtensionBackoffice}, false},
		{"doc-example", []string{ExtensionDoc, ExtensionExample}, false},
		{"doc-example-cli", []string{ExtensionDoc, ExtensionExample, ExtensionCli}, false},
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
		if group.Name == ExtensionDoc || sandbox.Deps.StringsDeps.HasPrefix(group.Name, ExtensionDoc+"-") {
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
// for that mechanic. A generated file is named twice — generated.<name>, and
// <name> as a build before GeneratedPrefix wrote it — so a purge leaves
// neither behind.
func ExtensionFiles(sandbox *api.Sandbox, name string) ([]string, error) {
	var paths []string

	for _, group := range GroupsRequiring(name) {
		files, err := sandbox.Deps.EmbedDeps.ListFilesRecursively(group)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			paths = append(paths, file)
			if IsGeneratedFile(sandbox, file) {
				base := baseName(sandbox, file)
				paths = append(paths, sandbox.Deps.StringsDeps.TrimSuffix(file, base)+SourceName(sandbox, base))
			}
		}
	}

	return paths, nil
}

// RenderExtensionCode renders, into an already open transaction and ahead of
// the build that follows, every code group turning one mechanic on turns on:
// each one that requires it and whose requirements the declaration now holds
// all on — its own, and one shared with another layer already on (database-cli
// renders on database-init with the cli on, and on cli-init with the database
// on). Everything else the mechanic writes — its pages above all — is left to
// that build, which has the collector output those templates need.
//
// A mechanic with no code group renders nothing here.
func RenderExtensionCode(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string) error {
	conf, err := LoadExtensionsConf(sandbox, io)
	if err != nil {
		return err
	}
	NormalizeExtensions(conf)

	for _, group := range AssetGroups() {
		if !group.Code || !contains(group.Requires, name) || !GroupEnabled(conf, group) {
			continue
		}

		module_conf, err := LoadModuleConf(sandbox, io)
		if err != nil {
			return err
		}

		if err := RenderGroup(sandbox, io, group.Name, map[string]interface{}{
			"Module": module_conf.Module,
		}); err != nil {
			return err
		}
	}

	return nil
}
