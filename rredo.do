redo-ifchange redo.rs
rustc --dep-info $2.dep -o $3 redo.rs
read DEPS < $2.dep
redo-ifchange ${DEPS#*:}
