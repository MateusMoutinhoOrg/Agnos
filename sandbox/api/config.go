package api

// Config is what the project knows about itself: the values of
// AgnosConfig/project.yaml, rendered into the sandbox by the build
// that read them. It is a field of the Sandbox like any other contract, so a
// caller may replace it — a test that runs the cli under another name, say —
// and every reader of it follows.
type Config struct {
	// ProjectConfig is a part of the Config, declared in sandbox/api/projectconfig.go.
	// Embedded, so each of its fields is read as sandbox.Config.<Field>.
	// Every struct of a sandbox/api/<x>config.go file is one, so a
	// mechanic adds its own part beside the project's rather than editing it.
	ProjectConfig

	// ProjectName is the project's name, the `project-name` key of
	// AgnosConfig/project.yaml; lower-cased it is the name the cli answers
	// to.
	ProjectName string

	// Version is the release the project is at, the `version` key of
	// AgnosConfig/project.yaml.
	Version string
}
