# The front-init example: add the html front layer to a project that has none
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q

agnos front-init --path TestDir

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. sandbox/internal/routes/static proves the server layer came with it —
# a page is answered over http — and assets/frontend holds the skeleton the
# helpers link against, which //go:embed would not keep if it were empty.
mkdir -p AssertDir/sandbox/internal/pageio
cp -R TestDir/sandbox/internal/pageio/. AssertDir/sandbox/internal/pageio/
mkdir -p AssertDir/sandbox/internal/routes/static
cp -R TestDir/sandbox/internal/routes/static/. AssertDir/sandbox/internal/routes/static/
mkdir -p AssertDir/assets/frontend
cp -R TestDir/assets/frontend/. AssertDir/assets/frontend/
