pub fn build(v: &::vars::Vars, targets: Vec<~str>) -> () {
    let _ = dofile::all_possible_do_files(&Path::new("foo"), &::std::os::getcwd());
    fail!("build {} {}", targets, v.DEPTH)
}

#[test]
fn test_builder() {
    build(&::vars::v(), Vec::new());
}

mod dofile {
    use std::path::Path;

    #[deriving(Eq, Clone)]
    struct DoPath {
        dodir: Path, // The directory containing the do-file
        basedir: Path, // The dir of the target relative to the do-file
        dofile: DoFile // The basename, target & extension of the do-file
    }

    #[deriving(Eq, Clone)]
    struct DoFile {
        dofile: Path,
        base: Path,
        ext: ~[u8]
    }

    static DOT: u8 = '.' as u8;
    static DEFAULT: &'static [u8] = bytes!("default");
    static DO: &'static [u8] = bytes!("do");

    type AllDoPathIter<'a> = ::std::iter::Chain<::std::option::Item<DoPath>, DefaultDoPathIter>;
    pub fn all_possible_do_files<'a>(filename: &Path, cwd: &Path) -> AllDoPathIter<'a> {
        let d: Path = Path::new(filename.dirname());
        let b: Path = Path::new(filename.filename().unwrap());
        Some(DoPath {
            dodir: d.clone(),
            basedir: Path::new(""),
            dofile: DoFile {
                dofile: Path::new([b.as_vec(), DO].connect_vec(&DOT)),
                base: b.clone(),
                ext: ~[]
            }
        }).move_iter().chain(DefaultDoPathIter::new(&d, &b, cwd))
    }


    struct DefaultDoPathIter {
        dodir: Path,
        basedir: Path,
        defaults: DefaultDoFileIter,
        up: Vec<~[u8]>,
        basename: Path
    }
    impl DefaultDoPathIter {
        fn new(dodir: &Path, basename: &Path, cwd: &Path) -> DefaultDoPathIter {
            let abs: Path = cwd.join(dodir);
            let comps: Vec<~[u8]> = abs.components().map(|x| x.to_owned()).collect();
            DefaultDoPathIter {
                dodir: dodir.clone(),
                basedir: Path::new(""),
                defaults: DefaultDoFileIter::new(basename),
                up: comps,
                basename: basename.clone()
            }
        }
    }
    impl Iterator<DoPath> for DefaultDoPathIter {
        fn next(&mut self) -> Option<DoPath> {
            self.defaults.next()
                .map(|d| {
                    DoPath {
                        dodir: self.dodir.clone(),
                        basedir: self.basedir.clone(),
                        dofile: d
                    }
                })
                .or_else(|| {
                    // look up one dir for the next defaults
                    self.up.pop().map(|comp| {
                        self.dodir.push("..");
                        self.basedir = Path::new(comp).join(&self.basedir);
                        self.defaults = DefaultDoFileIter::new(&self.basename);
                        DoPath {
                            dodir: self.dodir.clone(),
                            basedir: self.basedir.clone(),
                            dofile: self.defaults.next().unwrap()
                        }
                    })
                })
        }
    }


    struct DefaultDoFileIter {
        ext: Vec<~[u8]>,
        base: Vec<~[u8]>
    }
    impl DefaultDoFileIter {
        fn new(filename: &Path) -> DefaultDoFileIter {
            let fbytes: &[u8] = filename.as_vec();
            assert!(!fbytes.iter().any(::std::path::is_sep_byte));
            let parts: Vec<~[u8]> = fbytes.split(|b| *b == DOT).map(|p| p.to_owned()).collect();
            DefaultDoFileIter {ext: parts, base: Vec::new()}
        }
    }
    impl Iterator<DoFile> for DefaultDoFileIter {
        fn next(&mut self) -> Option<DoFile> {
            fn join_dots(v: &Vec<~[u8]>) -> ~[u8] {
                v.as_slice().connect_vec(&DOT)
            }

            fn mk_do_name(v: &Vec<~[u8]>) -> ~[u8] {
                let a4: Vec<~[u8]> =
                    Some(DEFAULT.to_owned()).move_iter()
                    .chain(v.clone().move_iter())
                    .chain(Some(DO.to_owned()).move_iter())
                    .collect();
                join_dots(&a4)
            }
            self.ext.shift().map(|x| {
                self.base.push(x);
                DoFile{
                    dofile: Path::new(mk_do_name(&self.ext)),
                    base: Path::new(join_dots(&self.base)),
                    ext: join_dots(&self.ext)
                }
            })
        }
    }



    #[cfg(test)]
    impl ::std::fmt::Show for DoPath {
        fn fmt(&self, f: &mut ::std::fmt::Formatter) -> ::std::fmt::Result {
            write!(f.buf, "DoPath\\{\n  dodir:\"{}\",\n  basedir:\"{}\",\n  dofile:{}\\}",
                   self.dodir.display(),
                   self.basedir.display(),
                   self.dofile)
        }
    }

    #[cfg(test)]
    impl ::std::fmt::Show for DoFile {
        fn fmt(&self, f: &mut ::std::fmt::Formatter) -> ::std::fmt::Result {
            write!(f.buf, "DoFile\\{\n    dofile:\"{}\",\n    base:\"{}\",\n    ext:\"{}\"\\}",
                   self.dofile.display(),
                   self.base.display(),
                   ::std::str::from_utf8(self.ext).unwrap())
        }
    }

    #[cfg(test)]
    mod tests {
        use builder::dofile::{DoPath, DoFile, all_possible_do_files, DefaultDoFileIter};
        fn p(x: &str) -> Path { return Path::new(x) }
        fn b(s: &str) -> ~[u8] { s.as_bytes().to_owned() }

        #[cfg(never)]
        fn dump<T : ::std::fmt::Show>(filename: &str, x: &T) {
            use std::io::fs;
            let mut f : fs::File =
                fs::File::create(&p(filename))
                .unwrap_or_handle(|e| fail!("create: {}", e));
            let () = write!(&mut f, "{}", x)
                .unwrap_or_handle(|e| fail!("write: {}", e));
        }

        #[test]
        fn do_file_paths() {
            let actual =
                all_possible_do_files(&p("foo/x.y"),
                                      &p("/bar")).collect();
            let expected = vec!(
                DoPath {
                    dodir: p("foo"),
                    basedir: p(""),
                    dofile: DoFile {
                        dofile: p("x.y.do"),
                        base: p("x.y"),
                        ext: b("")
                    }
                },
                DoPath {
                    dodir: p("foo"),
                    basedir: p(""),
                    dofile: DoFile {
                        dofile: p("default.y.do"),
                        base: p("x"),
                        ext: b("y")
                    }
                },
                DoPath {
                    dodir: p("foo"),
                    basedir: p(""),
                    dofile: DoFile {
                        dofile: p("default.do"),
                        base: p("x.y"),
                        ext: b("")
                    }
                },
                DoPath {
                    dodir: p("foo/.."),
                    basedir: p("foo"),
                    dofile: DoFile {
                        dofile: p("default.y.do"),
                        base: p("x"),
                        ext: b("y")
                    }
                },
                DoPath {
                    dodir: p("foo/.."),
                    basedir: p("foo"),
                    dofile: DoFile {
                        dofile: p("default.do"),
                        base: p("x.y"),
                        ext: b("")
                    }
                },
                DoPath {
                    dodir: p("foo/../.."),
                    basedir: p("bar/foo"),
                    dofile: DoFile {
                        dofile: p("default.y.do"),
                        base: p("x"),
                        ext: b("y")
                    }
                },
                DoPath {
                    dodir: p("foo/../.."),
                    basedir: p("bar/foo"),
                    dofile: DoFile {
                        dofile: p("default.do"),
                        base: p("x.y"),
                        ext: b("")
                    }
                });
            assert_eq!(expected, actual);
        }

        #[test]
        fn defaults_iter() {
            assert_eq!(
                vec!(
                    DoFile{
                        dofile: p("default.foo.bar.c.do"),
                        base: p("file"),
                        ext: b("foo.bar.c")
                    },
                    DoFile{
                        dofile: p("default.bar.c.do"),
                        base: p("file.foo"),
                        ext: b("bar.c")
                    },
                    DoFile{
                        dofile: p("default.c.do"),
                        base: p("file.foo.bar"),
                        ext: b("c")
                    },
                    DoFile{
                        dofile: p("default.do"),
                        base: p("file.foo.bar.c"),
                        ext: ~[]
                    }
                    ),
                DefaultDoFileIter::new(&p("file.foo.bar.c")).collect());
        }
    }
}
