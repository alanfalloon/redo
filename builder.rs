pub fn build(v: &::vars::Vars, targets: Vec<~str>) -> () {
    fail!("build {} {}", targets, v.DEPTH)
}

#[test]
fn test_builder() {
    build(&::vars::v(), Vec::new());
}

mod dofile {
    /* Finding do-files:
     *
     * for a file foo/bar/baz.a.b.c the following .do-files should be
     * looked for in this order:
     *  - foo/bar/baz.a.b.c.do
     *  - foo/bar/default.a.b.c.do
     *  - foo/bar/default.b.c.do
     *  - foo/bar/default.c.do
     *  - foo/bar/default.do
     *  - foo/default.a.b.c.do
     *  - foo/default.b.c.do
     *  - foo/default.c.do
     *  - foo/default.do
     *  - default.a.b.c.do
     *  - default.b.c.do
     *  - default.c.do
     *  - default.do
     *  - ../default.a.b.c.do
     *  - ../default.b.c.do
     *  - ../default.c.do
     *  - ../default.do
     *  - and so on...
     */
    #[test]
    fn do_file_search() {
        let mut dofile_iter = DoFileIter::new("foo/bar/baz.a.b.c");
        assert_eq!(Some(~"foo/bar/baz.a.b.c.do"), dofile_iter.next());
        assert_eq!(Some(~"foo/bar/default.a.b.c.do"), dofile_iter.next());
        assert_eq!(Some(~"foo/bar/default.b.c.do"), dofile_iter.next());
        assert_eq!(Some(~"foo/bar/default.c.do"), dofile_iter.next());
        assert_eq!(Some(~"foo/bar/default.do"), dofile_iter.next());
        assert_eq!(Some(~"foo/default.a.b.c.do"), dofile_iter.next());
        assert_eq!(Some(~"foo/default.b.c.do"), dofile_iter.next());
        assert_eq!(Some(~"foo/default.c.do"), dofile_iter.next());
        assert_eq!(Some(~"foo/default.do"), dofile_iter.next());
        assert_eq!(Some(~"default.a.b.c.do"), dofile_iter.next());
        assert_eq!(Some(~"default.b.c.do"), dofile_iter.next());
        assert_eq!(Some(~"default.c.do"), dofile_iter.next());
        assert_eq!(Some(~"default.do"), dofile_iter.next());
        assert_eq!(Some(~"../default.a.b.c.do"), dofile_iter.next());
        assert_eq!(Some(~"../default.b.c.do"), dofile_iter.next());
        assert_eq!(Some(~"../default.c.do"), dofile_iter.next());
        assert_eq!(Some(~"../default.do"), dofile_iter.next());
    }

    struct DoFileIter {
        file: Option<~str>,
        extensions: Vec<~str>,
        dir: ~str,
        dir_iter: DirIter,
        do_iter: DefaultDoIter
    }

    impl Iterator<~str> for DoFileIter {
        fn next(&mut self) -> Option<~str> {
            // If we haven't yielded the file, do so.
            let file = self.file.clone();
            self.file = None;
            match file {
                None => (),
                Some(f) => {
                    return Some(f + ".do");
                }
            }
            // Try the next default
            match self.do_iter.next() {
                Some(f) => {
                    if self.dir.len() > 0 {
                        Some(self.dir.clone() + "/" + f)
                    } else {
                        Some(f)
                    }
                }
                // We are out of defaults at this level, go down one
                // and freshen the default-do iter.
                None => match self.dir_iter.next() {
                    None => None, // No more dirs, we are done
                    Some(d) => {
                        self.dir = d;
                        self.do_iter = DefaultDoIter::new(self.extensions.clone());
                        self.next() // recurse to get the value
                    }
                }
            }
        }
    }

    impl DoFileIter {
        fn new(file: &str) -> DoFileIter {
            let exts : Vec<~str> = file.split('.').skip(1).map(|x| x.to_owned()).collect();
            let mut diter = DirIter::new(file);

            DoFileIter {
                file: Some(file.to_owned()),
                do_iter: DefaultDoIter::new(exts.clone()),
                extensions: exts,
                dir: diter.next().unwrap(),
                dir_iter: diter,
            }
        }
    }

    struct DirIter {
        dir_elts: Vec<~str>,
        down_depth: uint
    }

    impl Iterator<~str> for DirIter {
        fn next(&mut self) -> Option<~str> {
            let d = self.dir_elts.connect("/");
            if self.dir_elts.pop().is_none() {
                // d is empty, and we are iterating on [../]*down_depth
                let downs = Vec::from_fn(self.down_depth, |_| ~"..");
                self.down_depth += 1;
                Some(downs.connect("/"))
            } else {
                Some(d)
            }
        }
    }

    impl DirIter {
        fn new(file: &str) -> DirIter {
            let mut delts : Vec<~str> = file.split('/').map(|x| x.to_owned()).collect();
            delts.pop(); // The file part
            DirIter {
                down_depth: 0,
                dir_elts: delts
            }
        }
    }

    #[test]
    fn dir_iter() {
        assert_eq!(vec!(~"foo/bar",
                        ~"foo",
                        ~"",
                        ~"..",
                        ~"../..",
                        ),
                   DirIter::new("foo/bar/baz.c").take(5).collect());
    }


    struct DefaultDoIter {
        exts: Vec<~str>
    }

    impl Iterator<~str> for DefaultDoIter {
        fn next(&mut self) -> Option<~str> {
            let e = self.exts.connect(".");
            if self.exts.shift().is_none() {
                None
            } else {
                Some("default." + e)
            }
        }
    }

    impl DefaultDoIter {
        fn new(exts: Vec<~str>) -> DefaultDoIter {
            let mut e = exts;
            e.push(~"do");
            DefaultDoIter { exts: e }
        }
    }

    #[test]
    fn default_do_iter() {
        assert_eq!(vec!(~"default.foo.bar.c.do",
                        ~"default.bar.c.do",
                        ~"default.c.do",
                        ~"default.do"),
                   DefaultDoIter::new(vec!(~"foo", ~"bar", ~"c")).collect());
    }
}
