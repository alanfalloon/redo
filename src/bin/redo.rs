extern crate unix_socket;
extern crate rmp_serialize as msgpack;
extern crate redo;
extern crate rustc_serialize;

use rustc_serialize::Encodable;
use msgpack::Encoder;

use unix_socket::UnixStream;
use std::path::PathBuf;
use std::process::Command;

fn main() {
    let mut targets: Vec<String> = std::env::args().collect();
    let progname = targets.remove(0);
    let sock_path = redo::protocol::get_sock_path();
    let mut stream = match UnixStream::connect(&sock_path) {
        Ok(s) => s,
        Err(_) => {
            start_daemon(&progname);
            UnixStream::connect(&sock_path)
                .unwrap_or_else(|e| panic!("{}: {}", sock_path.display(), e))
        },
    };
    let mut encoder = Encoder::new(&mut stream);
    for id in 0..targets.len() {
        let req = redo::Request{
            id: id as u32,
            op: redo::Operation::RedoIfChange,
            target: From::from(&targets[id]),
        };
        req.encode(&mut encoder).unwrap();
    }
}

fn start_daemon(progname: &str) {
    let mut exe = PathBuf::from(progname);
    exe.pop();
    exe.push("redod");
    let _ = Command::new(exe).status().unwrap_or_else(|e| {
        panic!("failed to execute redo daemon: {}", e)
    });
}
