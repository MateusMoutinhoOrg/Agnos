# The {{ .Name }} example, run by `{{ .ProjectName }} exec-test`.
#
# Write it the way a reader would type it: an example is documentation first
# and a check second. It runs with this directory as the working directory,
# `{{ .ProjectName }}` on the PATH is this project's own cli built from source,
# and TestDir is the only place it may write. Everything it prints, the status
# it exits with and every file it leaves in TestDir become result.yaml.

mkdir -p TestDir
echo "the {{ .Name }} example"
