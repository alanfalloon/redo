redo-ifchange ../redo.rs
DEPTMP="$(mktemp "$2.XXXXX")"
trap "rm -f '$DEPTMP'" EXIT
rustc $RUSTFLAGS --dep-info "$DEPTMP" -o "$3" $(readlink -f ../redo.rs)
read DEPS < "$DEPTMP"
redo-ifchange ${DEPS#*:}
