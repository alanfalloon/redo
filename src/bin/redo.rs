extern crate redo;
extern crate rmp;
extern crate rmp_serialize as msgpack;
extern crate rustc_serialize;
extern crate unix_socket;
use msgpack::decode::Error::InvalidMarkerRead;
use msgpack::{Decoder, Encoder};
use redo::protocol::{Operation, Reply, Request, get_sock_path};
use rmp::decode::ReadError::UnexpectedEOF;
use rustc_serialize::{Decodable, Encodable};
use std::net::Shutdown;
use std::path::PathBuf;
use std::process::Command;
use unix_socket::UnixStream;

fn main() {
    let mut targets: Vec<String> = std::env::args().collect();
    let progname = targets.remove(0);
    let sock_path = get_sock_path();
    let mut stream = match UnixStream::connect(&sock_path) {
        Ok(s) => s,
        Err(_) => {
            start_daemon(&progname);
            UnixStream::connect(&sock_path)
                .unwrap_or_else(|e| panic!("{}: {}", sock_path.display(), e))
        },
    };
    {
        let mut encoder = Encoder::new(&mut stream);
        for id in 0..targets.len() {
            let req = Request{
                id: id as u32,
                op: Operation::RedoIfChange,
                target: From::from(&targets[id]),
            };
            req.encode(&mut encoder).unwrap();
        }
    }
    stream.shutdown(Shutdown::Write).unwrap();
    let mut decoder = Decoder::new(&mut stream);
    loop {
        let res: Reply = match Decodable::decode(&mut decoder) {
            Ok(res) => res,
            Err(InvalidMarkerRead(UnexpectedEOF)) => break,
            Err(e) => {
                println!("Unexpected Error: {:?}", e);
                break;
            },
        };
        println!("Response {:?}", res);
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
