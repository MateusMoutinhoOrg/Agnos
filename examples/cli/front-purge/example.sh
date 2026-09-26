# The front-purge example: remove the front layer, keeping assets/frontend
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos front-init --path TestDir -q
agnos add-page about --title "About" --path TestDir -q

agnos front-purge --path TestDir

# What result.yaml records: what the purge left. sandbox/internal/routeslist holds
# the health route alone — frontio and the frontend route are gone, because the
# route's handler imports a package that no longer exists — while
# assets/frontend is untouched, so front-init puts the route back over the
# same content.
mkdir -p AssertDir/sandbox/internal/routeslist
cp -R TestDir/sandbox/internal/routeslist/. AssertDir/sandbox/internal/routeslist/
mkdir -p AssertDir/assets/frontend
cp -R TestDir/assets/frontend/. AssertDir/assets/frontend/

# The declaration the pair wrote: this is the whole of what tells the build the
# mechanic is on or off from here.
mkdir -p AssertDir/AgnosConfig
cp TestDir/AgnosConfig/extensions.yaml AssertDir/AgnosConfig/extensions.yaml
