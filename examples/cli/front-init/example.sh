# The front-init example: add the html front layer to a project that has none
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q

agnos front-init --path test-dir

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. sandbox/deps/OpinionatedAgnosFront is the contract of the file layer
# the route serves through, sandbox/internal/routes/front proves the
# server layer came with it — the front is answered over http — and
# assets/front holds the index.html "/" answers, which also keeps
# //go:embed from dropping the tree.
mkdir -p assert-dir/sandbox/deps/OpinionatedAgnosFront
cp -R test-dir/sandbox/deps/OpinionatedAgnosFront/. assert-dir/sandbox/deps/OpinionatedAgnosFront/
mkdir -p assert-dir/sandbox/internal/routes/front
cp -R test-dir/sandbox/internal/routes/front/. assert-dir/sandbox/internal/routes/front/
mkdir -p assert-dir/assets/front
cp -R test-dir/assets/front/. assert-dir/assets/front/

# The declaration the pair wrote: this is the whole of what tells the build the
# mechanic is on or off from here.
mkdir -p assert-dir/AgnosConfig
cp test-dir/AgnosConfig/extensions.yaml assert-dir/AgnosConfig/extensions.yaml
