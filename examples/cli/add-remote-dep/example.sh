# The add-remote-dep example: install another agnos repo as a dep
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.
#
# TestDir/remote is the repo being installed — every agnos repo is installable
# by construction, so there is nothing to declare or turn on in it. TestDir/app
# is the consumer. A published repo is pinned with `add-dep <module>@<version>`;
# here the two live side by side, wired with a `replace`, which is how a pair of
# repos is developed together.

agnos start --path TestDir/remote --project-name Remote --module example/remote -q

mkdir -p TestDir/remote/sandbox/binds TestDir/remote/sandbox/internal/greeter

cat > TestDir/remote/sandbox/api/greeter.go <<'GO'
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

cat > TestDir/remote/sandbox/binds/greeter.go <<'GO'
package binds

import (
	api "example/remote/sandbox/api"
	greeter "example/remote/sandbox/internal/greeter"
)

func GreeterBind(sandbox *api.Sandbox) {
	sandbox.Greeter.Greet = func(props api.GreetProps) api.Greeting {
		return greeter.Greet(props)
	}
}
GO

cat > TestDir/remote/sandbox/internal/greeter/greeter.go <<'GO'
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

agnos build --path TestDir/remote -q

agnos start --path TestDir/app --project-name App --module example/app -q
agnos deps-init --path TestDir/app -q
printf '\nrequire example/remote v0.0.1\n\nreplace example/remote => ../remote\n' >> TestDir/app/go.mod

agnos add-dep example/remote --as remote --path TestDir/app

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The lib side copies the same set.
mkdir -p AssertDir/sandbox/deps
cp -R TestDir/app/sandbox/deps/. AssertDir/sandbox/deps/
mkdir -p AssertDir/adapters
cp -R TestDir/app/adapters/. AssertDir/adapters/
