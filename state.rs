
/**
 * Make targets relative to project root.
 */
pub fn canonical_targets(v: &::vars::Vars, targets: Vec<~str>) -> Vec<~str> {
    targets
}

#[test]
fn test_canonical_targets() {
    let mut v : ::vars::Vars = ::vars::v();
    let tgt = vec!(~"foo", ~"../bar", ~"baz/boz");
    let ct = canonical_targets(&v, tgt.clone());
    fail!("v={} tgt={} ct={}", v, tgt, ct);
}
