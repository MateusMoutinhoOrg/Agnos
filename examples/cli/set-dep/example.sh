# The set-dep example: re-copy a remote dep after its contract moved
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.
#
# The setup is the one in add-remote-dep: TestDir/remote is the repo being
# installed and TestDir/app the consumer, wired with a `replace` so the pair can
# be developed side by side. Here the remote contract then gains a field, and
# `set-dep` is what carries it across — the copy and the generated shim together.

agnos start --path TestDir/remote --project-name Remote --module example/remote -q

cat > TestDir/remote/sandbox/api/greeter.go <<'GO'
package api

// Greeter is what this repo publishes.
type Greeter struct {
	// Greet returns the greeting for one name.
	Greet func(props GreetProps) string
}

// GreetProps names who is greeted.
type GreetProps struct {
	Name string
}
GO

agnos build --path TestDir/remote -q

agnos start --path TestDir/app --project-name App --module example/app -q
agnos deps-init --path TestDir/app -q
printf '\nrequire example/remote v0.0.1\n\nreplace example/remote => ../remote\n' >> TestDir/app/go.mod
agnos add-dep example/remote --as remote --path TestDir/app -q

# The remote contract moves: one more field on the props it takes.
cat > TestDir/remote/sandbox/api/greeter.go <<'GO'
package api

// Greeter is what this repo publishes.
type Greeter struct {
	// Greet returns the greeting for one name.
	Greet func(props GreetProps) string
}

// GreetProps names who is greeted and how loudly.
type GreetProps struct {
	Name string
	Loud bool
}
GO

agnos build --path TestDir/remote -q

agnos set-dep remote --version v0.0.1 --path TestDir/app

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The lib side copies the same set.
mkdir -p AssertDir/sandbox/deps
cp -R TestDir/app/sandbox/deps/. AssertDir/sandbox/deps/
mkdir -p AssertDir/adapters
cp -R TestDir/app/adapters/. AssertDir/adapters/
