#![crate_id="redo#0.1"]
#![crate_type = "bin"]
#![feature(globs, phase)]

extern crate getopts;
use getopts::{optopt,optflag,getopts,usage};
use std::os;
use std::path::posix::Path;

fn main() -> () {
    let (flavour, targets) = match read_opts() {
        Ok(r) => r,
        Err(m) => {
            println!("{}", m);
            return;
        }
    };
    if !std::str::eq_slice(flavour, "redo-exec") {
        // FIXME: Run init here
    }
    // FIXME: actually run proper body.
    println!("{}: {}", flavour, targets);
}

fn read_opts() -> Result<(~str, Vec<~str>), ~str> {
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
    let args = os::args();
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
        None => match Path::new(progname.clone()).with_extension("").filename_str() {
            Some(m) => m.to_owned(),
            None => fail!("Can't match: " + progname)
        }
    };
    
    setenv_if_opt(&opts, "color", "REDO_COLOR");
    setenv_if_opt(&opts, "debug", "REDO_DEBUG");
    setenv_if_opt(&opts, "debug-locks", "REDO_DEBUG_LOCKS");
    setenv_if_opt(&opts, "debug-pids", "REDO_DEBUG_PIDS");
    setenv_if_opt(&opts, "keep-going", "REDO_KEEP_GOING");
    setenv_if_opt(&opts, "log", "REDO_LOG");
    setenv_if_opt(&opts, "only-log", "REDO_ONLY_LOG");
    setenv_if_opt(&opts, "overwrite", "REDO_OVERWRITE");
    setenv_if_opt(&opts, "shuffle", "REDO_SHUFFLE");
    setenv_if_opt(&opts, "verbose", "REDO_VERBOSE");
    setenv_if_opt(&opts, "warn-stdout", "REDO_WARN_STDOUT");
    setenv_if_opt(&opts, "xtrace", "REDO_XTRACE");

    return Ok((flavour, opts.free));
}

fn setenv_if_opt(opts: &getopts::Matches, opt: &str, env: &str) {
    if opts.opt_present(opt) {
        os::setenv(env, "1");
    }
}