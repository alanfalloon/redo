
use std::path::Path;
use std::io::{IoError, PathAlreadyExists, FileStat};
use std::io::fs::mkdir_recursive;
use time::{get_time, Timespec};

/**
 * Increase the mtime of the given file.
 *
 * If possible, the mtime is simply set to the current time, but if
 * that doesn't cause the mtime to increase, the mtime is explicitly
 * set forward 1 second. The file (and its directory) are created if
 * they don't exist.
 */
pub fn increment(file: Path) -> () {
    match mkdir_recursive(&file.dir_path(), 0o755) {
        Ok(()) => (),
        Err(IoError {kind: PathAlreadyExists, ..}) => (),
        Err(e) => fail!("mkdir({}) failed: {}", file.dir_path().display(), e)
    }
    let start_mtime = match file.stat() {
        Ok(FileStat{modified: m, ..}) => m,
        Err(_) => 0
    };
    let now = match get_time() {
        Timespec {sec: s, nsec: ns} => s * 1000 + ns as i64 / 1000
    } as u64;
    let new_mtime = if start_mtime >= now {
        start_mtime + 1
    } else {
        now
    };
    match ::native::io::file::utime(&file.to_c_str(), new_mtime, new_mtime) {
        Ok(()) => (),
        Err(e) => fail!("utime({}) failed: {}", file.display(), e)
    }
}
