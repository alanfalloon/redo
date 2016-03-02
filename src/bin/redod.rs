extern crate env_logger;
extern crate libc;
#[macro_use]
extern crate log;
extern crate redo;
extern crate rmp;
extern crate rmp_serialize as msgpack;
extern crate rustc_serialize;
extern crate unix_socket;
use msgpack::Encoder;
use redo::protocol::{Request, Reply, get_sock_path, StreamDecoder};
use rustc_serialize::Encodable;
use std::fs::{create_dir_all, remove_file};
use std::sync::atomic::{ATOMIC_USIZE_INIT, AtomicUsize, Ordering};
use std::sync::mpsc::channel;
use std::thread::spawn;
use unix_socket::{UnixListener, UnixStream};

static CONNS: AtomicUsize = ATOMIC_USIZE_INIT;

fn main() {
    env_logger::init().unwrap();
    let sock_path = get_sock_path();
    create_dir_all(sock_path.parent().unwrap()).unwrap();
    remove_file(&sock_path).unwrap();
    let listener = UnixListener::bind(&sock_path)
        .unwrap_or_else(|e| panic!("{}: {}", sock_path.display(), e));
    daemonize().unwrap();
    for (conn_id, stream) in listener.incoming().enumerate() {
        let stream = stream.unwrap();
        debug!("New connection {}.", &conn_id);
        CONNS.fetch_add(1, Ordering::AcqRel);
        spawn(move || {
            handle(stream);
            debug!("Done connection {}.", conn_id);
            if CONNS.fetch_sub(1, Ordering::AcqRel) == 1 {
                debug!("Goodbye.");
                std::process::exit(0);
            }
        });
    }
    unreachable!()
}

fn daemonize() -> Result<(), std::io::Error> {
    trace!("Daemonize.");
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
            trace!("Send {:?}", res);
            res.encode(&mut encoder).unwrap();
        }
    });
    for req in decoder {
        let req: Request = req.unwrap();
        let res_tx = res_tx.clone();
        spawn(move ||{
            debug!("Recv {:?}", req);
            let res = Reply::new(req.id, req.target);
            res_tx.send(res).unwrap();
        });
    }
    drop(res_tx);
    responder.join().unwrap();
}
