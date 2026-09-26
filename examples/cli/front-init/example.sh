# The front-init example: add the html front layer to a project that has none
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q

agnos front-init --path TestDir

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. sandbox/internal/routeslist/frontend proves the server layer came
# with it — the front is answered over http — and assets/frontend holds the
# index.html "/" answers, which also keeps //go:embed from dropping the tree.
mkdir -p AssertDir/sandbox/internal/generated/frontio
cp -R TestDir/sandbox/internal/generated/frontio/. AssertDir/sandbox/internal/generated/frontio/
mkdir -p AssertDir/sandbox/internal/routeslist/frontend
cp -R TestDir/sandbox/internal/routeslist/frontend/. AssertDir/sandbox/internal/routeslist/frontend/
mkdir -p AssertDir/assets/frontend
cp -R TestDir/assets/frontend/. AssertDir/assets/frontend/

# The declaration the pair wrote: this is the whole of what tells the build the
# mechanic is on or off from here.
mkdir -p AssertDir/AgnosConfig
cp TestDir/AgnosConfig/extensions.yaml AssertDir/AgnosConfig/extensions.yaml
