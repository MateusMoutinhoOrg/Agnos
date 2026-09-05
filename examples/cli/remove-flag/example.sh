# The remove-flag example: delete one declared flag
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos cli-init --path TestDir -q
agnos add-command greet --help "Greet someone" --category "Core" --path TestDir -q
agnos add-flag name --command greet --type string --default world --description "who to greet" --path TestDir -q

agnos remove-flag name --command greet --path TestDir
