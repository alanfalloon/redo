# Redo-in-rust TODO

This is the list of things left to implement on the rust
implementation of redo.

## Parity with the Python Implementation

The broad features that are missing, in the order of importance:

- `redo`
- `redo-ifchange`
- `redo-ifcreate`
- `redo-always`
- `redo-stamp`
- `redo-ood`, `redo-targets`, `redo-sources`, `redo-dofile`
- `--version` support and a release
- Keep-going (`-k`)
- Shuffle
- Parallel builds (`-j`)
- `redo-exec` -- This has a terrible name, it should be called
  `redo-background` or something, and since its whole purpose is
  spinning off long-running processes without having them hold-up the
  build, it should probably background the commands itself without
  having to background the `redo-exec` call in the shell.
- Fancy proctitles (less urgent since rust compiles to native)

Some stuff I don't know if I'll bother re-implementing:

- `redo-delegate` -- This whole thing about the temporary output
  directory is a mistake (notwithstanding
  [this thread](https://groups.google.com/d/topic/redo-list/RwYdXXp1riA/discussion)). Avery's
  original design of putting the temporary output in the same
  directory as the output file was the correct option. If the tool
  overwrites the supplementary targets non-atomically, it really isn't
  that big a deal because the main target will still not be replaced
  on failure, so a rebuild will rerun the script. Worst case you can
  add a command to preserve the files by moving them out of the way
  and restoring them on failure. (perhaps `redo-precious` like the
  Make `.PRECIOUS` directive). Maybe even declaring them as extra
  output afterwards (handling some of the worse cases like `javac`)
  automatically gives you the move-out-of-the-way behaviour on later
  runs.

  Then there is no need for `redo-delegate` and the strange rules
  about sticking output files in the same directory as the
  temporary. If you have one action that generates multiple targets,
  and you want dependencies on all of them, just do what they have
  been doing in Make for decades: choose one to be the main target
  that actually triggers the action, and have all the other outputs
  depend on that one.

  In fact, even with the odd temporary directory behaviour, why do we
  need `redo-delegate`? What does it do that `redo-ifchange` doesn't?

- `redo-log` -- The idea seems ok, I guess, but its a lot of effort
  for little gain.

## Above and Beyond the Python Implementation

- Since this is a rust implementation, and rust is still a moving
  target, it is critical to get this project added to the rust-nightly
  test runs in Travis CI.

- Also, since speed is one of the reasons to prefer rust, a little
  benchmark we can run would be nice.

- Add a something like the `make -w` option that prints the directory
  that we are in. This would make emacs work much better with the
  output. Its probably best to not print anything unless there is
  output on stderr, which implies we buffer it. Perhaps commands could
  opt-in by redirecting stderr to `$REDO_ERRBUF` or something.

- Support for searching up the tree for do files that aren't
  default.*.do. This is to support the fancy configuration building
  that Avery mentioned in the README (`debug/libfoo.a` and
  `release/libfoo.a` both made from `libfoo.a.do` instead of needing
  to use `default.a.do` (or perhaps some special naming to enable this
  feature like `@libfoo.a.do`).

- Support for the client/server model of building that Avery mentioned
  in item 4 of
  [this post](https://groups.google.com/d/topic/redo-list/doWmTj32UXc/discussion). Essentially,
  the first `redo` process becomes a build daemon, and sub-processes
  get a UNIX socket from the environment that they use to ask the root
  server to build targets.

- `redo-search` -- give it a list of names, and it returns the first
  one that exists or can be built. Earlier ones are implicitly
  `redo-ifcreate`, and for the one that is found it is implicitly
  `redo-ifchange`. This is to handle the common case of `default.o.do`
  that needs to support multiple source formats.

- Remove the proliferation of `.redo/` directories throughout the
  tree. The `apenwarr/master` branch was nicer since it had only the
  one meta-data tree at the root of the project, but as people rightly
  point out, the root is ad-hoc (set the first time you run `redo`).

  [`tup`](http://gittup.org/tup/) has a decent solution: the first time
  you run, it scans for a `Tupfile.ini` to mark the root of the
  project, otherwise you must run `tup init` manually before you can
  run.

  I would make it even easier, and assume that the `.git/` directory
  marks the root of the project. (I have no problems being opinionated
  about version control.) For projects outside a repo (release
  tarballs) I would assume the path of the first target built in the
  clean tree is at the root (since `redo all` or `redo install` seems
  like the most likely use-case in that scenario)

  Although, perhaps `tup` has it right. There is room for a
  configuration file in `redo`, if for nothing more than to decide
  what to do with stdout.

- Add an (optional?) authoritative source list. This was the main
  reason why Tim made [gup](https://github.com/gfxmonk/gup) instead of
  hacking on redo. The sources and targets are implicit in redo, if
  you build from a clean tree, the metadata will be correct, but if
  that metadata gets damaged, redo has no way of knowing what are
  sources and what are targets; existing files are a dead-end on a
  fresh build.

  I ran in to this when switching from make to redo on a project, the
  CI server was doing incremental builds on pull-requests, but I had
  reused the same flag-file to stamp the build, when the redo
  implementation was swapped in, `redo-ifchange build.flag` (inside
  `all.do`) did nothing because the `build.flag` file already existed.

  In my case, it would have been enough for redo to assume that only
  files under revision control were sources, and to clobber everything
  else if there is no metadata. But, failure is probably a reasonable
  default behaviour if there are files that are `redo-ifchange` that
  aren't in source control and we don't have `.redo/` to tell us what
  its dependencies are.

  And we can fall back to something explicit like the `Gupfile` if the
  default doesn't meet your needs. (This might be another good
  use-case for a config file -- to decide what to do if `.redo/` is
  missing).
