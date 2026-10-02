package commandprops

// Project is the part of CommandProps this project declares itself, embedded
// so each field is read as props.<Field>. A middleware that parsed a --path or
// a --profile sets it here, and the command it runs in front of reads it:
//
//	type Project struct {
//		Path string
//	}
//
// It lives under sandbox/internal, not in sandbox/api, so a field may name any
// type of the project, as long as that package imports no command.
//
// Written once by `{{.GeneratorName}} build` and the project's from then on.
type Project struct {
}
