
#[cfg(not(test))]
pub fn redo(flavour: &str, targets: Vec<~str>) -> () {
    fail!("{} {}", flavour, targets);
}
