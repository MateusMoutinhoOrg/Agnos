# The deps-purge example: remove the dependency layer and every installed dep
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos deps-init --path TestDir -q

agnos deps-purge --path TestDir
