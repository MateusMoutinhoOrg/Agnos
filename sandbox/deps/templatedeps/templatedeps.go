package templatedeps

// This package is the sandbox's *copy* of the api a text-template engine
// exposes — the same mechanic as argvdeps, embeddeps, iodeps, rundeps, std,
// stringsdeps and sortdeps, for the same reason: the sandbox may import
// nothing but the sandbox, so `text/template` and the `bytes` buffer it
// renders into may not appear inside it. The contract is restated here, and
// the adapter — which lives outside the sandbox — is what fills it.
//
// Rendering is the whole of the contract: agnos parses no template it does
// not immediately execute, so parse and execute are one call and no parsed
// template ever crosses the boundary.

// Lib is the template engine injected whole as the Deps.Templatedeps field.
type Lib struct {
	// Render parses one template source and executes it over the given vars,
	// returning the result. The error reports a source that does not parse or
	// an execution that failed — a native function returning an error
	// included.
	Render func(props RenderProps) (string, error)
}

// RenderProps is the whole input of one render.
type RenderProps struct {
	// Name is the template name, used in error messages only.
	Name string
	// Source is the template text to parse and execute.
	Source string
	// Vars is the value the template renders over, reached as `.`.
	Vars any
	// Funcs are the native functions the template may call, by the name each
	// is registered under. A value must be a function the engine accepts:
	// one return value, or one return value and an error.
	Funcs map[string]any
}
