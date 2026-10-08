# The add-remote-dep example: install another agnos repo as a dep
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.
#
# test-dir/remote is the repo being installed — every agnos repo is installable
# by construction, so there is nothing to declare or turn on in it. test-dir/app
# is the consumer. A published repo is pinned with `add-dep <module>@<version>`;
# here the two live side by side, wired with a `replace`, which is how a pair of
# repos is developed together.

agnos start --path test-dir/remote --project-name Remote --module example/remote -q

mkdir -p test-dir/remote/sandbox/internal/greeter

cat > test-dir/remote/sandbox/api/greeter.go <<'GO'
package api

// Greeter is what this repo publishes.
type Greeter struct {
	// Greet returns the greeting for one name.
	Greet func(props GreetProps) Greeting
}

// GreetProps names who is greeted and how loudly.
type GreetProps struct {
	Name string
	Loud bool
}

// Greeting is what came back.
type Greeting struct {
	Text  string
	Words int
}
GO

cat > test-dir/remote/sandbox/internal/greeter/new.go <<'GO'
package greeter

import (
	api "example/remote/sandbox/api"
)

func NewGreeter(sandbox *api.Sandbox) api.Greeter {
	greeter := api.Greeter{}

	greeter.Greet = func(props api.GreetProps) api.Greeting {
		return Greet(props)
	}

	return greeter
}
GO

cat > test-dir/remote/sandbox/internal/greeter/greeter.go <<'GO'
package greeter

import (
	api "example/remote/sandbox/api"
)

// Greet builds the greeting for one name.
func Greet(props api.GreetProps) api.Greeting {
	text := "hello " + props.Name
	if props.Loud {
		text = text + "!"
	}
	return api.Greeting{Text: text, Words: 2}
}
GO

# The Greeter is a field of the repo's own part of the Sandbox, which is what
# makes it part of the api another repo installs.
cat > test-dir/remote/sandbox/api/projectsandbox.go <<'GO'
package api

// ProjectSandbox is the part of the Sandbox this repo declares: the Greeter it
// publishes.
type ProjectSandbox struct {
	// Greeter is what this repo publishes.
	Greeter Greeter
}
GO

agnos build --path test-dir/remote -q

agnos start --path test-dir/app --project-name App --module example/app -q
agnos deps-init --path test-dir/app -q
printf '\nrequire example/remote v0.0.1\n\nreplace example/remote => ../remote\n' >> test-dir/app/go.mod

agnos add-dep example/remote --as remote --path test-dir/app

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. The lib side copies the same set.
mkdir -p assert-dir/sandbox/deps
cp -R test-dir/app/sandbox/deps/. assert-dir/sandbox/deps/
mkdir -p assert-dir/adapters
cp -R test-dir/app/adapters/. assert-dir/adapters/
