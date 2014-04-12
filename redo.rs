#![crate_id="redo#0.1"]
#![crate_type = "bin"]
#![feature(globs, phase)]

extern crate getopts;
extern crate native;
extern crate time;
use getopts::{optopt,optflag,getopts,usage};
use std::os;
use std::path::Path;

mod mains;
mod runid;

#[cfg(not(test))]
#[start]
pub fn start(argc: int, argv: **u8) -> int {
    native::start(argc, argv, main)
}

#[cfg(not(test))]
fn main() -> () {
    let (flavour, envs, targets) = match read_opts(os::args()) {
        Ok(r) => r,
        Err(m) => {
            println!("{}", m);
            return;
        }
    };
    for env in envs.iter() {
        os::setenv(*env, "1");
    }
    if !std::str::eq_slice(flavour, "redo-exec") {
        let bin_dir = os::getcwd().join(os::self_exe_path().unwrap());
        init_path(&bin_dir);
        match init_start_dir() {
            None => 0,
            Some(runidfile) => runid::increment(&runidfile)
        };
    }
    match flavour.as_slice() {
        "redo" => mains::redo(flavour, targets),
        _ => fail!("Unrecognized redo flavour: {}", flavour)
    }
}

fn read_opts(args: &[~str]) -> Result<(~str, Vec<~str>, Vec<~str>), ~str> {
    let opt_defs = ~[
        //    optopt("j", "jobs", "maximum number of jobs to build at once", "JOBS"),
        optflag("d", "debug", "print dependency checks as they happen"),
        optflag("v", "verbose", "print commands as they are read from .do files (variables intact)"),
        optflag("x", "xtrace", "print commands as they are executed (variables expanded)"),
        optflag("k", "keep-going", "keep going as long as possible even if some targets fail"),
        optflag("", "overwrite", "overwrite files even if generated outside of redo"),
        optflag("", "log", "activate log recording (slower)"),
        optflag("", "only-log", "print only failed targets from log"),
        optflag("", "shuffle", "randomize the build order to find dependency bugs"),
        optflag("", "debug-locks", "print messages about file locking (useful for debugging)"),
        optflag("", "debug-pids", "print process ids as part of log messages (useful for debugging)"),
        optflag("", "version", "print the current version and exit"),
        optflag("h", "help", "print the current help and exit"),
        optflag("", "color", "force enable color (--no-color to disable)"),
        optflag("", "warn-stdout", "warn if stdout is used"),
        optopt("", "main", "Choose which redo flavour to execute", "PROG-NAME"),
    ];
    let progname = args[0].clone();
    let opts = match getopts(args.tail(), opt_defs) {
        Ok(m) => m,
        Err(e) => return Err(e.to_err_msg())
    };
    if opts.opt_present("h") {
        return Err(usage(format!("{} [options] targets...

Rebuild targets", progname), opt_defs));
    };
    if opts.opt_present("version") {
        return Err(~"FIXME version here");
    };
    let flavour: ~str = match opts.opt_str("main") {
        Some(f) => f,
        None => Path::new(progname).with_extension("").filename_str().unwrap().to_owned()
    };
    
    // These are the environment variables that communicate the
    // options to all the child instances of redo.
    let mut env_set = Vec::new();
    if opts.opt_present("color") { env_set.push(~"REDO_COLOR"); }
    if opts.opt_present("debug") { env_set.push(~"REDO_DEBUG"); }
    if opts.opt_present("debug-locks") { env_set.push(~"REDO_DEBUG_LOCKS"); }
    if opts.opt_present("debug-pids") { env_set.push(~"REDO_DEBUG_PIDS"); }
    if opts.opt_present("keep-going") { env_set.push(~"REDO_KEEP_GOING"); }
    if opts.opt_present("log") { env_set.push(~"REDO_LOG"); }
    if opts.opt_present("only-log") { env_set.push(~"REDO_ONLY_LOG"); }
    if opts.opt_present("overwrite") { env_set.push(~"REDO_OVERWRITE"); }
    if opts.opt_present("shuffle") { env_set.push(~"REDO_SHUFFLE"); }
    if opts.opt_present("verbose") { env_set.push(~"REDO_VERBOSE"); }
    if opts.opt_present("warn-stdout") { env_set.push(~"REDO_WARN_STDOUT"); }
    if opts.opt_present("xtrace") { env_set.push(~"REDO_XTRACE"); }

    return Ok((flavour, env_set, opts.free));
}

#[test]
fn test_read_opts() -> () {
    let x = read_opts([~"redo", ~"--shuffle", ~"foo", ~"-x", ~"bar"]);
    assert!(x.eq(&Ok((~"redo", vec!(~"REDO_SHUFFLE", ~"REDO_XTRACE"), vec!(~"foo", ~"bar")))));
    let x = read_opts([~"/foo/bar/redo-ifchange.exe", ~"bar"]);
    assert!(x.eq(&Ok((~"redo-ifchange", vec!(), vec!(~"bar")))));
    let x = read_opts([~"/foo/bar/redo-ifchange.exe", ~"bar", ~"--main=redo-stamp"]);
    assert!(x.eq(&Ok((~"redo-stamp", vec!(), vec!(~"bar")))));
}

fn init_path(bin_dir: &Path) -> bool {
    // Ensure that REDO is set to the redo bin, and make sure that
    // the path to this executable is first in the PATH
    // environment variable.
    if os::getenv("REDO").is_none() {
        os::setenv("REDO", bin_dir.join("redo").as_str().unwrap());
        let new_path = match os::getenv("PATH") {
            None => bin_dir.as_str().unwrap().to_owned(),
            Some(p) => bin_dir.as_str().unwrap() + ":" + p
        };
        os::setenv("PATH", new_path);
        true
    } else {
        false
    }
}

#[test]
fn test_init_path() -> () {
    os::unsetenv("REDO");
    os::unsetenv("PATH");
    let bin_dir = Path::new("/foo/bar");
    assert!(init_path(&bin_dir));
    assert_env("REDO", "/foo/bar/redo");
    assert_env("PATH", "/foo/bar");
    assert!(!init_path(&bin_dir));
    os::unsetenv("REDO");
    os::setenv("PATH", "/bin:/sbin");
    assert!(init_path(&bin_dir));
    assert_env("REDO", "/foo/bar/redo");
    assert_env("PATH", "/foo/bar:/bin:/sbin");
}

#[cfg(test)]
fn assert_env(env: &str, val: &str) -> () {
    assert!(std::str::eq_slice(os::getenv(env).unwrap(), val));
}

fn init_start_dir() -> Option<Path> {
    // Ensure REDO_STARTDIR is captured. This is the top of the
    // build. This also indicates a new build starting, increment
    // the runid.
    if os::getenv("REDO_STARTDIR").is_none() {
        os::setenv("REDO_STARTDIR", os::getcwd().as_str().unwrap());
        os::setenv("REDO_RUNID_FILE", ".redo/runid");
        Some(Path::new(".redo/runid"))
    } else {
        None
    }
}

#[test]
fn test_start_dir() -> () {
    os::unsetenv("REDO_STARTDIR");
    os::unsetenv("REDO_RUNID_FILE");
    assert!(init_start_dir() == Some(Path::new(".redo/runid")));
    assert!(init_start_dir().is_none());
}
