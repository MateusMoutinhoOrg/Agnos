package commandprops

// Project is the part of CommandProps agnos declares itself, embedded so each
// field is read as props.<Field>: what the project middleware in front of
// every command hands on to the command after it.
type Project struct {
	// Path is the project directory every action works on, read off --path
	// by the project middleware in front of every command.
	Path string
}
