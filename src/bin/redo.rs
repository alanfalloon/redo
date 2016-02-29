extern crate redo;
extern crate rmp_serialize as msgpack;
extern crate rustc_serialize;
extern crate unix_socket;
use msgpack::Encoder;
use redo::protocol::{Operation, Reply, Request, get_sock_path, StreamDecoder};
use rustc_serialize::Encodable;
use std::net::Shutdown;
use std::path::{Path, PathBuf};
use std::process::Command;
use unix_socket::UnixStream;
use std::iter::Iterator;

fn main() {
    let mut targets: Vec<String> = std::env::args().collect();
    let progname = PathBuf::from(targets.remove(0));
    let op = progname.file_name().unwrap().to_str().unwrap();
    let op = Operation::from_str(&op).unwrap();
    let cwd = PathBuf::from(".").canonicalize().unwrap();
    let targets = targets.into_iter().enumerate().map(|(id, target)|{
        let target = cwd.join(target);
        Request {
            id: id as u32,
            op: op,
            target: target,
        }
    }).collect::<Vec<_>>();
    let mut stream = connect_to_redod(&progname);

    // Write all the command targets
    {
        let mut encoder = Encoder::new(&mut stream);
        for target in &targets {
            target.encode(&mut encoder).unwrap();
        }
    }
    stream.shutdown(Shutdown::Write).unwrap();

    // Collect the results
    for res in StreamDecoder::new(&mut stream) {
        let res: Reply = res.unwrap();
        let req = &targets[res.id as usize];
        println!("{:?} => {:?}", req, res);
    }
}

fn connect_to_redod(progname: &Path) -> UnixStream {
    let sock_path = get_sock_path();
    match UnixStream::connect(&sock_path) {
        Ok(s) => s,
        Err(_) => {
            start_daemon(progname);
            UnixStream::connect(&sock_path)
                .unwrap_or_else(|e| panic!("{}: {}", sock_path.display(), e))
        },
    }
}

fn start_daemon(progname: &Path) {
    let mut exe = PathBuf::from(progname);
    exe.pop();
    exe.push("redod");
    let _ = Command::new(exe).status().unwrap_or_else(|e| {
        panic!("failed to execute redo daemon: {}", e)
    });
}
