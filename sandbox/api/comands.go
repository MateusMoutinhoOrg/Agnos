package api

// CommandArg is one positional argument declared inside a command's entries.yaml.
type CommandArg struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description,omitempty"`
	Type        string   `yaml:"type,omitempty"`
	Required    bool     `yaml:"required,omitempty"`
	Examples    []string `yaml:"examples,omitempty"`
}

// CommandFlag is one flag declared inside a command's entries.yaml.
type CommandFlag struct {
	Name        string   `yaml:"name"`
	Identifiers []string `yaml:"identifiers,omitempty"`
	Description string   `yaml:"description,omitempty"`
	Type        string   `yaml:"type,omitempty"`
	Required    bool     `yaml:"required,omitempty"`
	Default     string   `yaml:"default,omitempty"`
	Array       bool     `yaml:"array,omitempty"`
	Examples    []string `yaml:"examples,omitempty"`
}

// Command is the top-level structure of a command's entries.yaml.
type Command struct {
	Identifiers     []string      `yaml:"identifiers,omitempty"`
	Category        string        `yaml:"category,omitempty"`
	Help            string        `yaml:"help,omitempty"`
	LongDescription string        `yaml:"long-description,omitempty"`
	Examples        []string      `yaml:"examples,omitempty"`
	Args            []CommandArg  `yaml:"args,omitempty"`
	Flags           []CommandFlag `yaml:"flags,omitempty"`
}
