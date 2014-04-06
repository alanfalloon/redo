redo-ifchange smoke bin/test
bin/redo t/all
[ "$(cat t/result)" = "ok" ]
[ -n "$DO_BUILT" ] || echo "Don't forget to test 'minimal/do test'" >&2
