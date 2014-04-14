redo-ifchange redo.test
env -i RUST_LOG=redo ./redo.test
touch "$3"
