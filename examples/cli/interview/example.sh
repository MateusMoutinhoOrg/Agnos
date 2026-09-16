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

# A third session, on the same project once it has a cli and one command of
# its own — declared here with the plain cli, which is the cheap way to reach
# the state the session is about.
#
# The answers: the Cli System area, add-flag, a flag named `name` with no
# --identifier of its own, on `greet`, of type string, described, with no
# --example, defaulting to "world", not an array, at no particular position,
# and yes to running it. `--required` is never asked — a flag carrying a
# default cannot also be required, so the session does not offer the
# combination add-flag would refuse.

agnos cli-init --path TestDir -q
agnos add-command greet --help "Greet someone" --category Core --path TestDir -q

printf '3\n3\nname\n\n1\n1\nwho to greet\n\nworld\n1\n\n1\n' | agnos interview --path TestDir

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The declaration `start` wrote and the one the third session added a
# flag to are the whole of what the sessions changed on disk.
mkdir -p AssertDir/AgnosConfig AssertDir/sandbox/internal/commands/greet
cp TestDir/AgnosConfig/project.yaml AssertDir/AgnosConfig/project.yaml
cp TestDir/AgnosConfig/extensions.yaml AssertDir/AgnosConfig/extensions.yaml
cp TestDir/sandbox/internal/commands/greet/entries.yaml AssertDir/sandbox/internal/commands/greet/entries.yaml
