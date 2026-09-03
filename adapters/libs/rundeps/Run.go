package rundeps

import (
	"bytes"
	"os/exec"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	rundeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/rundeps"
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
