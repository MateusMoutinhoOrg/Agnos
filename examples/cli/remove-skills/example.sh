# The remove-skills example: take the agent skill away again
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos drop-skills --path test-dir -q

agnos remove-skills --path test-dir

# Both roots are gone with the skill: nothing else was in them.
test ! -e test-dir/.agents || exit 1
test ! -e test-dir/.claude || exit 1

# What result.yaml records: the declaration the pair wrote, which is the whole
# of what tells the build the mechanic is off from here.
mkdir -p assert-dir/AgnosConfig
cp test-dir/AgnosConfig/extensions.yaml assert-dir/AgnosConfig/extensions.yaml
