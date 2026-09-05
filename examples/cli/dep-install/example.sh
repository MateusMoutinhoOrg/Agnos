# The dep-install example: install one dep of the built-in list
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos deps-init --path TestDir -q

agnos dep-install iodeps --path TestDir
