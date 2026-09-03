package rundeps

// This package is the sandbox's *copy* of the api a process runner exposes —
// the same mechanic as argvdeps, dbdeps, embeddeps, iodeps, requestdeps and
// std, for the same reason: starting a child process is an OS-bound effect,
// so `os/exec` may not appear inside the sandbox. The contract is restated
// here, and the adapter — which lives outside the sandbox — is what fills it.

// Lib is the process runner injected whole as the Deps.Rundeps field. It is
// what the build action reaches for when it has to hand the rendered project
// to a real toolchain (`go mod tidy`, `go build ./...`) and report whether
// that toolchain accepted it.
type Lib struct {
	// Run executes one program to completion and returns what it wrote.
	// A non-zero exit status is reported in Result.ExitCode, not as the
	// error: the error is reserved for a program that could not be started
	// at all (missing binary, unreadable directory).
	Run func(props RunProps) (Result, error)
}

// RunProps describes one program invocation.
type RunProps struct {
	// Dir is the working directory the program runs in. "" means the
	// current directory.
	Dir string
	// Program is the executable to run, looked up on PATH.
	Program string
	// Args are the arguments handed to the program, excluding its own name.
	Args []string
}

// Result is what one finished invocation produced.
type Result struct {
	// Output is the program's standard output and standard error, merged
	// in the order they were written.
	Output string
	// ExitCode is the program's exit status; 0 means success.
	ExitCode int
}
