extern crate rustc_serialize;

use rustc_serialize::Encodable;

pub mod protocol;

#[derive(RustcEncodable, RustcDecodable, PartialEq, Debug)]
pub enum Operation {
    RedoIfChange,
    RedoIfCreate,
    Redo,
}

#[derive(RustcEncodable, RustcDecodable, PartialEq, Debug)]
pub struct Request {
    pub id: u32,
    pub op: Operation,
    pub target: std::path::PathBuf,
}

#[derive(RustcEncodable, RustcDecodable, PartialEq, Debug)]
pub struct Reply {
    id: u32,
    target: Vec<u8>,
}

#[test]
fn it_works() {
}
