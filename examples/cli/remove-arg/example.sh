# The remove-arg example: delete one declared positional argument
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos cli-init --path TestDir -q
agnos add-command greet --help "Greet someone" --category "Core" --path TestDir -q
agnos add-arg who --command greet --type string --required --description "the person to greet" --path TestDir -q

agnos remove-arg who --command greet --path TestDir
