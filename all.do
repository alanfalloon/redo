exec >&2
redo-ifchange bin/redod bin/redocli
export PATH="$PWD/bin:$PATH"
redocli -v foo.bar
bin/redocli --quiet -- baz.quux
