# The drop-skills example: write the agent skill into both roots
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q

agnos drop-skills --path test-dir

# What result.yaml records: the skill under each root — the same bytes twice,
# so whichever tool opens the project reads the same pages — and the
# declaration that keeps every later build writing it.
mkdir -p assert-dir/.agents/skills assert-dir/.claude/skills
cp -R test-dir/.agents/skills/. assert-dir/.agents/skills/
cp -R test-dir/.claude/skills/. assert-dir/.claude/skills/
mkdir -p assert-dir/AgnosConfig
cp test-dir/AgnosConfig/extensions.yaml assert-dir/AgnosConfig/extensions.yaml
