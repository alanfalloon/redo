use std::path::PathBuf;
use std::env::home_dir;

pub fn get_sock_path() -> PathBuf {
    let mut sock_path = home_dir().expect("No HOME directory");
    sock_path.push(".redo");
    sock_path.push("redod.sock");
    sock_path
}
