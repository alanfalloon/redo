extern crate rustc_serialize;
extern crate rmp;
extern crate rmp_serialize as msgpack;
#[macro_use]
extern crate log;

pub mod cmd;
pub mod protocol;
pub mod state;
pub mod world;

#[test]
fn it_works() {
}

pub trait Fs {
}

pub use state::State;
pub use world::World;
