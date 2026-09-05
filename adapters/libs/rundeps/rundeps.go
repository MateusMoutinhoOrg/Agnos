package rundeps

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	rundeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/rundeps"
)

// Bind fills deps.Deps.Rundeps, providing the capability to run one external
// program to completion over the standard library's os/exec.
func Bind(deps *deps.Deps) {
	deps.Rundeps = rundeps.Lib{
		Run: run,
	}
}

// run fills rundeps.Lib.Run, executing the program in props.Dir and merging
// its standard output and standard error into one buffer. A non-zero exit
// status comes back in Result.ExitCode; only a program that could not be
// started at all is reported as an error.
func run(props rundeps.RunProps) (rundeps.Result, error) {
	command := exec.Command(props.Program, props.Args...)
	command.Dir = props.Dir

	if len(props.PathPrefix) > 0 {
		path := prefixedPath(props.PathPrefix)
		command.Env = append(os.Environ(), "PATH="+path)
		command.Env = append(command.Env, props.Env...)
		lookInPath(command, path, props.Program)
	} else if len(props.Env) > 0 {
		command.Env = append(os.Environ(), props.Env...)
	}

	var output bytes.Buffer
	command.Stdout = &output
	command.Stderr = &output

	err := command.Run()
	result := rundeps.Result{Output: output.String()}

	if err == nil {
		return result, nil
	}

	exitError, ok := err.(*exec.ExitError)
	if !ok {
		return result, err
	}

	result.ExitCode = exitError.ExitCode()
	return result, nil
}

// prefixedPath is the PATH value a PathPrefix asks for: the requested
// directories, in order, ahead of the one this process inherited. Reading the
// current PATH is why the join happens here and not in the sandbox.
func prefixedPath(prefix []string) string {
	return strings.Join(append(append([]string{}, prefix...), os.Getenv("PATH")), string(os.PathListSeparator))
}

// lookInPath re-resolves the program against the prefixed PATH, so a name the
// prefix provides is the one that runs. exec.Command already looked it up in
// this process's PATH; a hit here replaces that answer, and a miss leaves it —
// together with the lookup error exec.Command recorded when there was none.
func lookInPath(command *exec.Cmd, path string, program string) {
	if strings.ContainsRune(program, os.PathSeparator) {
		return
	}
	for _, dir := range filepath.SplitList(path) {
		candidate := filepath.Join(dir, program)
		info, err := os.Stat(candidate)
		if err != nil || info.IsDir() || info.Mode()&0111 == 0 {
			continue
		}
		command.Path = candidate
		command.Err = nil
		return
	}
}
