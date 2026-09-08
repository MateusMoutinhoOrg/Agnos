# The {{ .Name }} example, run by `{{ .ProjectName }} exec-test`.
#
# Write it the way a reader would type it: an example is documentation first
# and a check second. It runs with this directory as the working directory,
# `{{ .ProjectName }}` on the PATH is this project's own cli built from source,
# and TestDir is the only place it may write.

mkdir -p TestDir
echo "the {{ .Name }} example" > TestDir/example.txt
echo "wrote TestDir/example.txt"

# What result.yaml records is AssertDir, not TestDir: copy into it the
# directories this example is about, keeping the place each one holds in the
# tree. TestDir stays as it is, for reading. Copy nothing and the run fails —
# an example that asserts nothing passes for the wrong reason. The lib side
# copies the same set, or the cli-vs-lib check breaks over the copy itself.
mkdir -p AssertDir
cp -R TestDir/. AssertDir/
