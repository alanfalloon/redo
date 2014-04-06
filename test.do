redo-ifchange _all
redo bin/test t/all
[ -n "$DO_BUILT" ] || echo "Don't forget to test 'minimal/do test'" >&2
