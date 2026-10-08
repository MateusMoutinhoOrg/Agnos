# The interview example: the guided session, answered from a pipe
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.
#
# Stdin is not a terminal here, so the ttyinterview adapter asks every question
# as a numbered list read line by line — which is what makes the session
# scriptable. The answers below are the ones a person would give on an empty
# folder: the suggested first step, a name, no to --force, a module path, the
# folder it offers, and yes to running it.

mkdir -p test-dir

printf '1\nTest\n1\nTest\n\n1\n' | agnos interview --path test-dir

# The same session again, answering nothing: it prints what the project needs
# next now that it exists. The first menu is the whole point of the feature —
# with no cli there is no Cli area, and `cli-init` is the suggested
# step instead.

agnos interview --path test-dir < /dev/null

# A third session, on the same project once it has a cli and one command of
# its own — declared here with the plain cli, which is the cheap way to reach
# the state the session is about.
#
# The answers: the Cli area, add-flag, the command `greet` — asked
# first, although the command line spells it after the name, because it is what
# the other questions are about — then a flag named `name` with no --key of its
# own, of type string, described, defaulting to "world", at no particular
# position, with no --enum, --pattern or --trigger, and yes to running it —
# then nothing else for `greet`, and out.
# `--required` is never asked — a flag carrying a default cannot also be
# required, so the session does not offer the combination add-flag would
# refuse, and the confirm screen says which questions went and why.

agnos cli-init --path test-dir -q
agnos add-command greet --summary "Greet someone" --category Core --path test-dir -q

# The area is row 6: this project's first menu is led by `· exit`, because every
# step left on it is an offer to install a whole layer and no menu of this
# session puts one of those under an enter pressed blind; the three steps and
# the Core area come before Cli.

printf '6\n3\n1\nname\n\n1\nwho to greet\nworld\n\n\n\n\n1\n3\n1\n' | agnos interview --path test-dir

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. The declaration `start` wrote and the one the third session added a
# flag to are the whole of what the sessions changed on disk.
mkdir -p assert-dir/AgnosConfig assert-dir/sandbox/internal/commands/core/greet
cp test-dir/AgnosConfig/project.yaml assert-dir/AgnosConfig/project.yaml
cp test-dir/AgnosConfig/extensions.yaml assert-dir/AgnosConfig/extensions.yaml
cp test-dir/sandbox/internal/commands/core/greet/command.yaml assert-dir/sandbox/internal/commands/core/greet/command.yaml
