# The add-cli-example example: create an example under examples/cli/
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos cli-init --path TestDir -q

agnos add-cli-example greet --path TestDir
