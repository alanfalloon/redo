extern crate unix_socket;
extern crate rmp_serialize as msgpack;
extern crate rmp;
extern crate redo;
extern crate rustc_serialize;
extern crate libc;

use rustc_serialize::Decodable;
use msgpack::Decoder;

use unix_socket::{UnixListener, UnixStream};

fn main() {
    let sock_path = redo::protocol::get_sock_path();
    std::fs::create_dir_all(sock_path.parent().unwrap()).unwrap();
    std::fs::remove_file(&sock_path).unwrap();
    let listener = UnixListener::bind(&sock_path)
        .unwrap_or_else(|e| panic!("{}: {}", sock_path.display(), e));
    daemonize().unwrap();
    for stream in listener.incoming() {
        let stream = stream.unwrap();
        std::thread::spawn(|| handle(stream));
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
    let mut decoder = Decoder::new(&mut stream);
    loop {
        let req: redo::Request = match Decodable::decode(&mut decoder) {
            Ok(req) => req,
            Err(msgpack::decode::Error::InvalidMarkerRead(rmp::decode::ReadError::UnexpectedEOF)) => break,
            Err(e) => {
                println!("Unexpected Error: {:?}", e);
                break;
            },
        };
        println!("Request {:?}", req);
    }
}
