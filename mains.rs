
#[cfg(not(test))]
pub fn redo(v: &::vars::Vars, flavour: &str, targets: Vec<~str>) -> () {
    let rerooted = ::state::canonical_targets(v, targets);
    ::builder::build(v, rerooted);
}
