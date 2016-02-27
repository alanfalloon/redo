extern crate libc;
extern crate redo;
extern crate rmp;
extern crate rmp_serialize as msgpack;
extern crate rustc_serialize;
extern crate unix_socket;
use msgpack::decode::Error::InvalidMarkerRead;
use msgpack::{Decoder, Encoder};
use redo::protocol::{Request, Reply, get_sock_path};
use rmp::decode::ReadError::UnexpectedEOF;
use rustc_serialize::{Decodable, Encodable};
use std::fs::{create_dir_all, remove_file};
use std::sync::mpsc::channel;
use std::thread::spawn;
use unix_socket::{UnixListener, UnixStream};

fn main() {
    let sock_path = get_sock_path();
    create_dir_all(sock_path.parent().unwrap()).unwrap();
    remove_file(&sock_path).unwrap();
    let listener = UnixListener::bind(&sock_path)
        .unwrap_or_else(|e| panic!("{}: {}", sock_path.display(), e));
    daemonize().unwrap();
    for stream in listener.incoming() {
        let stream = stream.unwrap();
        spawn(|| handle(stream));
    }
}

fn daemonize() -> Result<(), std::io::Error> {
    let pid = unsafe { libc::fork() };
    if pid < 0 {
        Err(std::io::Error::last_os_error())
    } else if pid != 0 {
        std::process::exit(0)
    } else {
        Ok(())
    }
}

fn handle(mut stream: UnixStream) {
    let mut decoder = Decoder::new(stream.try_clone().unwrap());
    let (res_tx, res_rx) = channel::<Reply>();
    let responder = spawn(move || {
        let mut encoder = Encoder::new(&mut stream);
        for res in res_rx {
            res.encode(&mut encoder).unwrap();
        }
    });
    loop {
        let req: Request = match Decodable::decode(&mut decoder) {
            Ok(req) => req,
            Err(InvalidMarkerRead(UnexpectedEOF)) => break,
            Err(e) => {
                println!("Unexpected Error: {:?}", e);
                break;
            },
        };
        let res_tx = res_tx.clone();
        spawn(move ||{
            let res = Reply { id: req.id, target: req.target.clone() };
            println!("Request {:?}", req);
            res_tx.send(res).unwrap();
        });
    }
    drop(res_tx);
    responder.join().unwrap();
}
