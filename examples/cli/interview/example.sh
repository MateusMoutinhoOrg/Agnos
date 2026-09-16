# The interview example: the guided session, answered from a pipe
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.
#
# Stdin is not a terminal here, so the interviewer adapter asks every question
# as a numbered list read line by line — which is what makes the session
# scriptable. The answers below are the ones a person would give on an empty
# folder: the suggested first step, the folder it offers, a name, no to
# --force, a module path, and yes to running it.

mkdir -p TestDir

printf '1\n\nTest\n1\nTest\n1\n' | agnos interview --path TestDir

# The same session again, answering nothing: it prints what the project needs
# next now that it exists. The first menu is the whole point of the feature —
# with no cli there is no Cli System area, and `cli-init` is the suggested
# step instead.

agnos interview --path TestDir < /dev/null

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The declaration `start` wrote is the whole of what the session
# changed on disk.
mkdir -p AssertDir/AgnosConfig
cp TestDir/AgnosConfig/project.yaml AssertDir/AgnosConfig/project.yaml
cp TestDir/AgnosConfig/extensions.yaml AssertDir/AgnosConfig/extensions.yaml
