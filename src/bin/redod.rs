extern crate libc;
extern crate redo;
extern crate rmp;
extern crate rmp_serialize as msgpack;
extern crate rustc_serialize;
extern crate unix_socket;
use msgpack::Encoder;
use redo::protocol::{Request, Reply, get_sock_path, StreamDecoder};
use rustc_serialize::Encodable;
use std::fs::{create_dir_all, remove_file};
use std::io::Write;
use std::sync::atomic::{ATOMIC_USIZE_INIT, AtomicUsize, Ordering};
use std::sync::mpsc::channel;
use std::thread::spawn;
use unix_socket::{UnixListener, UnixStream};

static CONNS: AtomicUsize = ATOMIC_USIZE_INIT;

fn main() {
    let sock_path = get_sock_path();
    create_dir_all(sock_path.parent().unwrap()).unwrap();
    remove_file(&sock_path).unwrap();
    let listener = UnixListener::bind(&sock_path)
        .unwrap_or_else(|e| panic!("{}: {}", sock_path.display(), e));
    daemonize().unwrap();
    for stream in listener.incoming() {
        let stream = stream.unwrap();
        CONNS.fetch_add(1, Ordering::AcqRel);
        spawn(|| {
            handle(stream);
            if CONNS.fetch_sub(1, Ordering::AcqRel) == 1 {
                let stderr = std::io::stderr();
                write!(stderr.lock(), "Goodbye.").unwrap();
                std::process::exit(0);
            }
        });
    }
    unreachable!()
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

fn handle(mut stream_tx: UnixStream) {
    let mut stream_rx = stream_tx.try_clone().unwrap();
    let decoder = StreamDecoder::new(&mut stream_rx);
    let (res_tx, res_rx) = channel::<Reply>();
    let responder = spawn(move || {
        let mut encoder = Encoder::new(&mut stream_tx);
        for res in res_rx {
            res.encode(&mut encoder).unwrap();
        }
    });
    for req in decoder {
        let req: Request = req.unwrap();
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
