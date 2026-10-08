# The {{ .ExampleName }} example, run by `{{ .GeneratorName }} run-examples`.
#
# Write it the way a reader would type it: an example is documentation first
# and a check second. It runs with this directory as the working directory,
# `{{ .ProjectName }}` on the PATH is this project's own cli built from source,
# and test-dir is the only place it may write.

mkdir -p test-dir
echo "the {{ .ExampleName }} example" > test-dir/example.txt
echo "wrote test-dir/example.txt"

# What result.yaml records is assert-dir, not test-dir: copy into it the
# directories this example is about, keeping the place each one holds in the
# tree. test-dir stays as it is, for reading. Copy nothing and the run fails —
# an example that asserts nothing passes for the wrong reason. The lib side
# copies the same set, or the cli-vs-lib check breaks over the copy itself.
mkdir -p assert-dir
cp -R test-dir/. assert-dir/
