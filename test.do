redo-ifchange _all
env -i - PATH="$PWD/bin:$PATH" bin/redo t/all >&2
[ -n "$DO_BUILT" ] || echo "Don't forget to test 'minimal/do test'" >&2
