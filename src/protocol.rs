use rustc_serialize::Encodable;
use std::env::home_dir;
use std::path::PathBuf;

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
    pub target: PathBuf,
}

#[derive(RustcEncodable, RustcDecodable, PartialEq, Debug)]
pub struct Reply {
    pub id: u32,
    pub target: PathBuf,
}

pub fn get_sock_path() -> PathBuf {
    let mut sock_path = home_dir().expect("No HOME directory");
    sock_path.push(".redo");
    sock_path.push("redod.sock");
    sock_path
}
