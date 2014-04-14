
pub fn build(v: &::vars::Vars, targets: Vec<~str>) -> () {
    fail!("build {} {}", targets, v.DEPTH)
}

#[test]
fn test_builder() {
    build(&::vars::v(), Vec::new());
}
