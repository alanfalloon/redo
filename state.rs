/**
 * Make targets relative to the original $PWD of the .do script
 *
 * The actions in the .do script can chdir() before calling
 * redo-ifchanged and friends. Adjust the paths so that they are
 * correct from the original $PWD.
 */
pub fn un_chdir_targets(v: &::vars::Vars, cwd: &Path, targets: Vec<~str>) -> (Vec<~str>, Path) {
    let ocwd = Path::new(v.STARTDIR.clone()).join(v.PWD.clone());
    if cwd == &ocwd {
        return (targets, ocwd);
    }
    let rel = cwd.path_relative_from(&ocwd).unwrap();
    (targets.iter().map(|tgt| rel.join(tgt.as_slice()).as_str().unwrap().to_owned()).collect(),
     ocwd)
}

#[test]
fn test_un_chdir_targets() {
    let mut v : ::vars::Vars = ::vars::v();
    let cwd = Path::new("/foo/bar/baz");
    v.STARTDIR = ~"/foo/bar";
    let tgt = vec!(~"foo", ~"../bar", ~"buzz/boz", ~"/foo/bar/baz/x");

    v.PWD = ~"baz";
    match un_chdir_targets(&v, &cwd, tgt.clone()) {
        (_tgt, _cwd) => {
            assert_eq!(_tgt, vec!(~"foo", ~"../bar", ~"buzz/boz", ~"/foo/bar/baz/x"));
            assert_eq!(_cwd.as_str().unwrap(), "/foo/bar/baz");
        }
    };
            

    v.PWD = ~"baz/buzz";
    match un_chdir_targets(&v, &cwd, tgt.clone()) {
        (_tgt, _cwd) => {
            assert_eq!(_tgt, vec!(~"../foo", ~"../../bar", ~"../buzz/boz", ~"/foo/bar/baz/x"));
            assert_eq!(_cwd.as_str().unwrap(),
                       "/foo/bar/baz/buzz");
        }
    };

    v.PWD = ~"";
    match un_chdir_targets(&v, &cwd, tgt.clone()) {
        (_tgt, _cwd) => {
            assert_eq!(_tgt, vec!(~"baz/foo", ~"bar", ~"baz/buzz/boz", ~"/foo/bar/baz/x"));
            assert_eq!(_cwd.as_str().unwrap(), "/foo/bar");
        }
    };
}
