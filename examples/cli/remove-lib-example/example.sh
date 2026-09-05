# The remove-lib-example example: delete an example of examples/lib/
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos add-lib-example greet --path TestDir -q

agnos remove-lib-example greet --path TestDir
