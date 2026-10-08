# The front-purge example: remove the front layer, keeping assets/front
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos front-init --path test-dir -q
agnos add-page about --title "About" --path test-dir -q

agnos front-purge --path test-dir

# What result.yaml records: what the purge left. sandbox/internal/routes holds
# the health route alone — the front route is gone with the layer it
# belongs to — while
# assets/front is untouched, so front-init puts the route back over the
# same content.
mkdir -p assert-dir/sandbox/internal/routes
cp -R test-dir/sandbox/internal/routes/. assert-dir/sandbox/internal/routes/
mkdir -p assert-dir/assets/front
cp -R test-dir/assets/front/. assert-dir/assets/front/

# The declaration the pair wrote: this is the whole of what tells the build the
# mechanic is on or off from here.
mkdir -p assert-dir/AgnosConfig
cp test-dir/AgnosConfig/extensions.yaml assert-dir/AgnosConfig/extensions.yaml
