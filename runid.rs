
use std::path::Path;
use std::io::{
    File,
    FileNotFound,
    FileStat,
    IoError,
    PathAlreadyExists,
};
use std::io::fs::mkdir_recursive;
#[cfg(not(test))]
use time::get_time;
#[cfg(test)]
use time::Timespec;



/**
 * Read the runid from the given file
 *
 * The file is created it if it doesn't exist.
 */
pub fn read(file: &Path) -> u64 {
    let (_, r) = read_created(file);
    r
}

fn read_created(file: &Path) -> (bool, u64) {
    match file.stat() {
        Ok(FileStat{modified: m, ..}) => (false, m / 1000),
        Err(IoError {kind: FileNotFound, ..}) => (true, new_file_mtime(file)),
        Err(e) => fail!("stat({}) failed: {}", file.display(), e)
    }
}

/**
 * Increase the mtime of the given file.
 *
 * If possible, the mtime is simply set to the current time, but if
 * that doesn't cause the mtime to increase, the mtime is explicitly
 * set forward 1 second. The file (and its directory) are created if
 * they don't exist.
 */
pub fn increment(file: &Path) -> u64 {
    let (created, start_mtime) = read_created(file);
    if created {
        return start_mtime
    }
    let now = get_time().sec as u64;
    let new_mtime = if start_mtime >= now {
        start_mtime + 1
    } else {
        now
    };
    let new_mtime_ms = new_mtime * 1000;
    match ::native::io::file::utime(&file.to_c_str(), new_mtime_ms, new_mtime_ms) {
        Ok(()) => (),
        Err(e) => fail!("utime({}) failed: {}", file.display(), e)
    };
    new_mtime
}

fn new_file_mtime(file: &Path) -> u64 {
    match mkdir_recursive(&file.dir_path(), 0o755) {
        Ok(()) => (),
        Err(IoError {kind: PathAlreadyExists, ..}) => (),
        Err(e) => fail!("mkdir({}) failed: {}", file.dir_path().display(), e)
    }
    match File::create(file) {
        Err(e) => fail!("create({}) failed: {}", file.display(), e),
        _ => ()
    }
    // Really? No fstat? I guess I'll need to create a pull-request.
    match file.stat() {
        Ok(FileStat{modified: m, ..}) => m / 1000,
        Err(e) => fail!("stat({}) after create failed: {}", file.display(), e)
    }
 }

#[test]
fn test_runid() -> () {
    let tmpdir_holder = ::std::io::TempDir::new("runidtest").unwrap();
    let tmpdir = tmpdir_holder.path();
    let p = tmpdir.join("a/b/c");
    let runid1 = increment(&p);
    println!("runid1 = {}", runid1);
    assert!(runid1 < 2_000_000_000);
    assert_eq!(read(&p), runid1);
    assert_eq!(read(&p), runid1);
    assert_eq!(increment(&p), 2_000_000_000);
    assert_eq!(read(&p), 2_000_000_000);
    assert_eq!(increment(&p), 2_000_000_001);
    assert_eq!(read(&p), 2_000_000_001);
    assert_eq!(increment(&p), 2_000_000_002);
    assert_eq!(read(&p), 2_000_000_002);
    assert_eq!(read(&p), 2_000_000_002);
}

#[cfg(test)]
fn get_time() -> Timespec {Timespec {sec: 2_000_000_000, nsec: 123456}}
