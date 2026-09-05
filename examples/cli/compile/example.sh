# The compile example: cross-compile a project's cmd/main into release/
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos cli-init --path TestDir -q

agnos compile --target linux86 --path TestDir
