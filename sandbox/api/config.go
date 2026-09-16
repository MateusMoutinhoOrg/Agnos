package api

// Config is what the project knows about itself: the values of
// <ProjectName>Config/project.yaml, rendered into the sandbox by the build
// that read them. It is a field of the Sandbox like any other contract, so a
// caller may replace it — a test that runs the cli under another name, say —
// and every reader of it follows.
type Config struct {
	// ProjectName is the project's name, title-cased. It prefixes the
	// <ProjectName>Config/ directory that holds every declaration, and
	// lower-cased it is the name the cli answers to.
	ProjectName string

	// Version is the release the project is at, the `version` key of
	// <ProjectName>Config/project.yaml.
	Version string
}
