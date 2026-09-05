# The dep-list example: list the deps available to a project
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos deps-init --path TestDir -q
agnos dep-install iodeps --path TestDir -q

agnos dep-list --path TestDir
