# The set-dep example: re-copy a remote dep after its contract moved
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.
#
# The setup is the one in add-remote-dep: test-dir/remote is the repo being
# installed and test-dir/app the consumer, wired with a `replace` so the pair can
# be developed side by side. Here the remote contract then gains a field, and
# `set-dep` is what carries it across — the copy and the generated shim together.

agnos start --path test-dir/remote --project-name Remote --module example/remote -q

cat > test-dir/remote/sandbox/api/greeter.go <<'GO'
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

agnos build --path test-dir/remote -q

agnos start --path test-dir/app --project-name App --module example/app -q
agnos deps-init --path test-dir/app -q
printf '\nrequire example/remote v0.0.1\n\nreplace example/remote => ../remote\n' >> test-dir/app/go.mod
agnos add-dep example/remote --as remote --path test-dir/app -q

# The remote contract moves: one more field on the props it takes.
cat > test-dir/remote/sandbox/api/greeter.go <<'GO'
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

agnos build --path test-dir/remote -q

agnos set-dep remote --version v0.0.1 --path test-dir/app

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. The lib side copies the same set.
mkdir -p assert-dir/sandbox/deps
cp -R test-dir/app/sandbox/deps/. assert-dir/sandbox/deps/
mkdir -p assert-dir/adapters
cp -R test-dir/app/adapters/. assert-dir/adapters/
