use std::sync::Arc;
pub struct World {
    _x: u8,
}
impl World {
    pub fn new() -> Arc<World> {
        Arc::new(World{_x:0})
    }
}
impl super::Fs for World {
}
