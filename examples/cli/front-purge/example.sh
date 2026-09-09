# The front-purge example: remove the front layer, keeping assets/frontend
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos front-init --path TestDir -q
agnos add-page home --trigger / --title "Home" --path TestDir -q

agnos front-purge --path TestDir

# What result.yaml records: what the purge left. sandbox/internal/routes holds
# the health route alone — pageio, the static route and the page's route are
# gone, because their handlers import a package that no longer exists — while
# assets/frontend is untouched, so front-init + add-page put the routes back
# over the same content.
mkdir -p AssertDir/sandbox/internal/routes
cp -R TestDir/sandbox/internal/routes/. AssertDir/sandbox/internal/routes/
mkdir -p AssertDir/assets/frontend
cp -R TestDir/assets/frontend/. AssertDir/assets/frontend/
