package exec_tests

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/rundeps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/projectconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/resultconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// aliasDir is where the cli alias every example types is written, relative to
// the project. It holds one executable named after the project, so `agnos` in
// an example.sh runs this tree's ./cmd/main and never an installed binary.
// release/ is the project's git-ignored binary directory, which is exactly
// what this is.
const aliasDir = "release/exec-test"

// volatileFiles and volatileDirs are the TestDir paths left out of a result's
// tree. An example that reaches the Go runtime (`start` runs `build`, which
// runs `go mod tidy` and `go build`) writes a go.sum holding whatever the
// module proxy resolved that day and binaries under release/: neither is a
// property of the project, so neither can be a golden.
var volatileFiles = []string{"go.sum"}
var volatileDirs = []string{"release/"}

// exampleRun is one example about to be executed.
type exampleRun struct {
	Side string
	Name string
	Dir  string
}

// ExecTestInternal runs every planned example, compares what it produced with
// its golden result.yaml (or writes that golden), and reports how many
// examples failed. An example present on both sides is run through the cli
// first and then through the lib, and the two runs are cross-checked: the cli
// is a wrapper over the lib, so they must leave the same tree and exit the
// same way.
func ExecTestInternal(deps *deps.Deps, path string, only string, update bool) error {
	root, err := projectRoot(deps, path)
	if err != nil {
		return err
	}

	prefix, err := writeCliAlias(deps, path, root)
	if err != nil {
		return err
	}

	runs, err := planRuns(deps, path, only)
	if err != nil {
		return err
	}

	failed := []string{}
	produced := map[string]*resultconf.ResultConf{}

	for _, run := range runs {
		result, err := execExample(deps, path, root, prefix, run)
		if err != nil {
			return err
		}
		produced[run.Side+"/"+run.Name] = result

		divergences, err := checkGolden(deps, path, run, result, update)
		if err != nil {
			return err
		}
		if len(divergences) > 0 {
			failed = append(failed, run.Side+"/"+run.Name)
			report(deps, run.Side+"/"+run.Name, divergences)
		}
	}

	crossed := crossCheckedNames(runs)
	for _, name := range crossed {
		divergences := crossCheck(produced[utils.ExampleCliSide+"/"+name], produced[utils.ExampleLibSide+"/"+name])
		if len(divergences) > 0 {
			failed = append(failed, name)
			report(deps, name+" (cli vs lib)", divergences)
		}
	}

	deps.Std.Log("exec-test done: %d examples, %d cross-checks, %d failed \n", len(runs), len(crossed), len(failed))

	if len(failed) > 0 {
		return deps.Std.Errorf("exec-test: %d of %d checks failed: %s",
			len(failed), len(runs)+len(crossed), strings.Join(failed, ", "))
	}
	return nil
}

// planRuns lists the examples to run, in the order they are run: alphabetical
// by name, cli before lib. `only` narrows the plan to one name — both sides of
// it — and is a usage error when no side declares that name.
func planRuns(deps *deps.Deps, path string, only string) ([]exampleRun, error) {
	only = strings.TrimSpace(only)

	names := []string{}
	seen := map[string]bool{}
	has := map[string]bool{}

	for _, side := range utils.ExampleSides {
		for _, name := range listExamples(deps, path, side) {
			has[side+"/"+name] = true
			if !seen[name] {
				seen[name] = true
				names = append(names, name)
			}
		}
	}
	sort.Strings(names)

	if only != "" {
		if !seen[only] {
			return nil, deps.Std.Errorf("exec-test: no example named %q (declared: %s)",
				only, strings.Join(names, ", "))
		}
		names = []string{only}
	}

	runs := []exampleRun{}
	for _, name := range names {
		for _, side := range utils.ExampleSides {
			if has[side+"/"+name] {
				runs = append(runs, exampleRun{Side: side, Name: name, Dir: utils.ExampleDir(side, name)})
			}
		}
	}
	return runs, nil
}

// listExamples lists one side's example names, sorted. It reads disk directly:
// exec-test opens no SmartIO, so nothing here is filtered or buffered.
func listExamples(deps *deps.Deps, path string, side string) []string {
	dir := join(path, utils.ExampleSideDir(side))
	if !deps.Iodeps.IsDir(dir) {
		return nil
	}

	var names []string
	for _, entry := range deps.Iodeps.ListDirs(dir) {
		names = append(names, utils.LastSegment(entry))
	}
	sort.Strings(names)
	return names
}

// crossCheckedNames lists the example names the plan ran on both sides, in
// plan order: only those can be cross-checked.
func crossCheckedNames(runs []exampleRun) []string {
	sides := map[string][]string{}
	var names []string

	for _, run := range runs {
		if len(sides[run.Name]) == 0 {
			names = append(names, run.Name)
		}
		sides[run.Name] = append(sides[run.Name], run.Side)
	}

	var both []string
	for _, name := range names {
		if len(sides[name]) == len(utils.ExampleSides) {
			both = append(both, name)
		}
	}
	return both
}

// execExample runs one example from scratch: its TestDir is removed straight
// off disk (no buffer to persist), the example is executed with its own
// directory as the working directory, and what it produced is gathered into a
// result — the exit status and merged output of the run, plus every file the
// TestDir ended up holding.
func execExample(deps *deps.Deps, path string, root string, prefix []string, run exampleRun) (*resultconf.ResultConf, error) {
	dir := join(path, run.Dir)
	test_dir := dir + "/" + utils.ExampleTestDir

	deps.Iodeps.RemoveDir(test_dir)

	program, args := invocation(run.Side)
	result, err := deps.Rundeps.Run(rundeps.RunProps{
		Dir:        dir,
		Program:    program,
		Args:       args,
		PathPrefix: prefix,
	})
	if err != nil {
		return nil, deps.Std.Errorf("exec-test %s/%s: could not run %s: %w", run.Side, run.Name, program, err)
	}

	conf := resultconf.NewEmpty(deps)
	conf.CliOutput = normalize(result.Output, root+"/"+run.Dir)
	conf.ExitCode = result.ExitCode

	for _, entry := range treeOf(deps, test_dir) {
		conf.AddTreeEntry(entry.File, entry.Sha)
	}

	deps.Std.Log("exec-test %s/%s: exit %d, %d files \n", run.Side, run.Name, conf.ExitCode, len(conf.Tree))
	return conf, nil
}

// invocation is how one side's example is started: a shell for the cli side,
// the Go runtime for the lib side.
func invocation(side string) (string, []string) {
	if side == utils.ExampleCliSide {
		return "sh", []string{utils.ExampleCliFile}
	}
	return "go", []string{"run", utils.ExampleLibFile}
}

// treeOf is every file inside one example's TestDir, ordered by path relative
// to that directory and carrying the sha256 of its content. A file that
// cannot be read is recorded with an empty sha rather than aborting the suite:
// the comparison then reports it like any other divergence.
func treeOf(deps *deps.Deps, test_dir string) []resultconf.TreeEntry {
	var entries []resultconf.TreeEntry

	for _, file := range deps.Iodeps.ListFilesRecursively(test_dir) {
		name := strings.TrimPrefix(strings.TrimPrefix(file, test_dir), "/")
		if name == "" || isVolatile(name) {
			continue
		}

		content, err := deps.Iodeps.ReadFile(file)
		if err != nil {
			entries = append(entries, resultconf.TreeEntry{File: name})
			continue
		}

		sum := sha256.Sum256(content)
		entries = append(entries, resultconf.TreeEntry{File: name, Sha: hex.EncodeToString(sum[:])})
	}

	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].File < entries[j].File
	})
	return entries
}

// isVolatile reports whether a TestDir-relative path is one the tree leaves
// out (see volatileFiles / volatileDirs).
func isVolatile(name string) bool {
	for _, file := range volatileFiles {
		if name == file {
			return true
		}
	}
	for _, dir := range volatileDirs {
		if strings.HasPrefix(name, dir) {
			return true
		}
	}
	return false
}

// normalize is what makes a golden portable: the absolute directory the
// example ran in becomes <dir>, and carriage returns are dropped. Any other
// absolute path, timestamp or resolved version left in the output belongs to
// the machine that ran it, and the example carrying it is not a valid one.
func normalize(output string, dir string) string {
	output = strings.ReplaceAll(output, "\r\n", "\n")
	return strings.ReplaceAll(output, dir, "<dir>")
}

// checkGolden compares one produced result with the example's result.yaml, or
// writes that file when it does not exist yet or --update was passed. It
// returns one line per divergence, empty when the example passed.
func checkGolden(deps *deps.Deps, path string, run exampleRun, produced *resultconf.ResultConf, update bool) ([]string, error) {
	golden_path := join(path, run.Dir+"/"+utils.ExampleResultFile)

	if update || !deps.Iodeps.IsFile(golden_path) {
		deps.Std.Log("exec-test %s/%s: writing %s \n", run.Side, run.Name, run.Dir+"/"+utils.ExampleResultFile)
		return nil, deps.Iodeps.WriteFile(golden_path, []byte(produced.Render()))
	}

	content, err := deps.Iodeps.ReadFile(golden_path)
	if err != nil {
		return nil, deps.Std.Errorf("exec-test %s/%s: could not read %s: %w", run.Side, run.Name, golden_path, err)
	}

	golden, err := resultconf.New(deps, string(content))
	if err != nil {
		return nil, deps.Std.Errorf("exec-test %s/%s: %s: %w", run.Side, run.Name, golden_path, err)
	}

	divergences := diffExitCode(golden.ExitCode, produced.ExitCode)
	divergences = append(divergences, diffOutput(golden.CliOutput, produced.CliOutput)...)
	return append(divergences, diffTree(golden.Tree, produced.Tree)...), nil
}

// crossCheck holds the cli side of an example to the lib side: the same tree
// and the same exit status, which is the assertion that the cli is only a
// wrapper over the lib. cli-output is deliberately left out — each side has
// its own text (a cli error on one, a panic on the other) and is checked
// against its own golden.
func crossCheck(cli *resultconf.ResultConf, lib *resultconf.ResultConf) []string {
	if cli == nil || lib == nil {
		return nil
	}
	divergences := diffExitCode(cli.ExitCode, lib.ExitCode)
	return append(divergences, diffTree(cli.Tree, lib.Tree)...)
}

// diffExitCode reports a differing exit status.
func diffExitCode(expected int, got int) []string {
	if expected == got {
		return nil
	}
	return []string{"exit-code: expected " + itoa(expected) + ", got " + itoa(got)}
}

// diffOutput reports the differing lines of the two outputs, `-` for the
// expected line and `+` for the one produced.
func diffOutput(expected string, got string) []string {
	if expected == got {
		return nil
	}

	expected_lines := strings.Split(expected, "\n")
	got_lines := strings.Split(got, "\n")

	lines := []string{"cli-output:"}
	for index := 0; index < len(expected_lines) || index < len(got_lines); index++ {
		left := lineAt(expected_lines, index)
		right := lineAt(got_lines, index)
		if left == right {
			continue
		}
		if left != nil {
			lines = append(lines, "  - "+*left)
		}
		if right != nil {
			lines = append(lines, "  + "+*right)
		}
	}
	return lines
}

// lineAt is the line at index, or nil when the side ended before it.
func lineAt(lines []string, index int) *string {
	if index >= len(lines) {
		return nil
	}
	return &lines[index]
}

// diffTree reports every file the two trees do not agree on: `+` for a file
// only the produced tree has, `-` for one only the expected tree has, and `~`
// for one both have with a different sha. Saying only "failed" about two trees
// of hundreds of files is no help at all.
func diffTree(expected []resultconf.TreeEntry, got []resultconf.TreeEntry) []string {
	expected_shas := shasOf(expected)
	got_shas := shasOf(got)

	var lines []string
	for _, entry := range expected {
		sha, ok := got_shas[entry.File]
		if !ok {
			lines = append(lines, "  - "+entry.File)
			continue
		}
		if sha != entry.Sha {
			lines = append(lines, "  ~ "+entry.File)
		}
	}
	for _, entry := range got {
		if _, ok := expected_shas[entry.File]; !ok {
			lines = append(lines, "  + "+entry.File)
		}
	}

	if len(lines) == 0 {
		return nil
	}
	sort.Strings(lines)
	return append([]string{"tree:"}, lines...)
}

// shasOf indexes a tree by file path.
func shasOf(entries []resultconf.TreeEntry) map[string]string {
	shas := map[string]string{}
	for _, entry := range entries {
		shas[entry.File] = entry.Sha
	}
	return shas
}

// report prints one failed check and everything that diverged. It goes to the
// error channel, not the log one: --quiet silences progress, never the reason
// a run failed.
func report(deps *deps.Deps, label string, divergences []string) {
	deps.Std.Error("exec-test %s: FAILED\n", label)
	for _, line := range divergences {
		deps.Std.Error("  %s\n", line)
	}
}

// projectRoot is the absolute directory the project sits in. The sandbox
// cannot resolve a path itself, and every example runs with its own directory
// as the working directory, so the alias and the output normalization both
// need the answer a child `pwd` gives.
func projectRoot(deps *deps.Deps, path string) (string, error) {
	result, err := deps.Rundeps.Run(rundeps.RunProps{Dir: path, Program: "pwd"})
	if err != nil {
		return "", deps.Std.Errorf("exec-test: could not resolve %s: %w", path, err)
	}
	if result.ExitCode != 0 {
		return "", deps.Std.Errorf("exec-test: could not resolve %s:\n%s", path, result.Output)
	}
	return strings.TrimRight(result.Output, "\r\n"), nil
}

// writeCliAlias writes the executable an example types — named after the
// project, running this tree's ./cmd/main through `go run` — and returns the
// directory to put in front of the PATH of every run. An example is
// documentation first, so it types `agnos`, not a path to a binary; this is
// what makes that name resolve to the code in the repo rather than to an
// installed release. A project with no cli has nothing to alias and gets no
// prefix.
func writeCliAlias(deps *deps.Deps, path string, root string) ([]string, error) {
	if !deps.Iodeps.IsDir(join(path, "cmd/main")) {
		return nil, nil
	}

	conf, err := loadProjectConf(deps, path)
	if err != nil {
		return nil, err
	}

	alias := aliasDir + "/" + conf.Name
	script := "#!/bin/sh\nexec go run " + root + "/cmd/main \"$@\"\n"

	if err := deps.Iodeps.WriteFile(join(path, alias), []byte(script)); err != nil {
		return nil, err
	}

	// The filesystem contract writes content, not a mode, and a PATH entry
	// only answers to an executable file.
	result, err := deps.Rundeps.Run(rundeps.RunProps{Dir: path, Program: "chmod", Args: []string{"755", alias}})
	if err != nil {
		return nil, deps.Std.Errorf("exec-test: could not make %s executable: %w", alias, err)
	}
	if result.ExitCode != 0 {
		return nil, deps.Std.Errorf("exec-test: could not make %s executable:\n%s", alias, result.Output)
	}

	return []string{root + "/" + aliasDir}, nil
}

// loadProjectConf reads the project.yaml the alias is named after. exec-test
// opens no SmartIO, so it reads the file straight off disk.
func loadProjectConf(deps *deps.Deps, path string) (*projectconf.ProjectConf, error) {
	rel := config.ProjectName + "Config/project.yaml"

	content, err := deps.Iodeps.ReadFile(join(path, rel))
	if err != nil {
		return nil, deps.Std.Errorf("could not read %s: run `agnos start` first (%w)", rel, err)
	}
	return projectconf.New(deps, string(content))
}

// join is the project-relative path helper this action needs in place of the
// SmartIO boundary it does without: "" and "." mean the current directory, so
// they add no prefix.
func join(path string, rel string) string {
	path = strings.TrimSuffix(strings.TrimSpace(path), "/")
	if path == "" || path == "." {
		return rel
	}
	return path + "/" + rel
}

// itoa renders one exit status; the sandbox has no strconv-free alternative
// worth writing, and only small non-negative numbers reach it.
func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	sign := ""
	if value < 0 {
		sign = "-"
		value = -value
	}
	digits := ""
	for value > 0 {
		digits = string(rune('0'+value%10)) + digits
		value /= 10
	}
	return sign + digits
}
