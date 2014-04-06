# Redo-in-rust TODO
This is the list of things left to implement on the rust implementation of redo.

## Parity with the Python Implementation
The broad features that are missing, in the order of importance:
- redo
- redo-ifchange
- redo-ifcreate
- redo-always
- redo-ood, redo-targets, redo-sources, redo-dofile
- redo-stamp
- Keep-going (-k)
- Shuffle
- Parallel builds (-j)
- redo-exec -- This has a terrible name, it should be called
  redo-background or something, and since its whole purpose is
  spinning off long-running processes without having them hold-up the
  build, it should probably background the commands itself without
  having to background the redo-exec call in the shell.
- Fancy proctitles (less urgent since rust compiles to native)

Some stuff I don't know if I'll bother reimplementing:
- redo-delegate -- This whole thing about the temporary output
  directory is a mistake. Avery's original plan of putting the
  temporary output in the same directory as the output file was the
  correct option. Then there is no need for redo-delegate and the
  strange rules about sticking output files in the same directory as
  the temporary. If you have one action that generates multiple
  targets, and you want dependencies on all of them, just do what they
  have been doing in Make for decades: choose one to be the main
  target that actually triggers the action, and have all the other
  outputs depend on that one.
- redo-log -- The idea seems ok, I guess, but its a lot of effort for
  little gain.

## Above and Beyond the Python Implementation

- Add a something like Makes -w option that prints the directory that
  we are in. This would make emacs work much better with the output.

Since this is a rust implementation, and rust is still a moving
target, it is critical to get this project added to the rust-nightly
test runs in Travis CI.

Also, since speed is one of the reasons to prefer rust, a little
benchmark we can run would be nice.
