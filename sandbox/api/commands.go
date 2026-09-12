package api

// Command is one command the project declares, as the sandbox sees it.
type Command struct {
}

// Commands is the command surface of the sandbox: every command declared in
// sandbox/internal/commands/, in declaration order.
type Commands struct {
	List []Command
}
