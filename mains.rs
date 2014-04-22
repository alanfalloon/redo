
#[cfg(not(test))]
pub fn redo(v: &::vars::Vars, flavour: &str, targets: Vec<~str>) -> () {
    assert_eq!("redo", flavour);
    let cwd = ::std::os::getcwd();
    let (rerooted, new_cwd) = ::state::un_chdir_targets(v, &cwd, targets);
    assert!(::std::os::change_dir(&new_cwd));
    ::builder::build(v, rerooted);
}
