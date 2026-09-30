package commandprops

// CommandProps is what one command line carries from the commands that run on
// it to the ones after them: the dispatch builds one, empty, per command line,
// and hands the same one to every InternalPureHandler of the chain as its
// first argument. A middleware that parsed a --path or a --profile sets it
// here, and the command it runs in front of reads it:
//
//	type CommandProps struct {
//		Path string
//	}
//
// It lives under sandbox/internal, not in sandbox/api, so a field may name any
// type of the project, as long as that package imports no command.
//
// Written once by `{{.GeneratorName}} build` and then yours: declare the
// fields your middlewares hand on. No build rewrites this file once it is
// there.
type CommandProps struct {
}
