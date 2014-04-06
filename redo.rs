#![crate_id="redo#0.1"]
#![crate_type = "bin"]
#![feature(globs, phase)]

extern crate getopts;
use getopts::{optopt,optflag,getopts,usage};
use std::os;
use std::path::posix::Path;

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
        // FIXME: Run init here
    }
    // FIXME: actually run proper body.
    println!("{}: {} {}", flavour, envs, targets);
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
        None => match Path::new(progname.clone()).with_extension("").filename_str() {
            Some(m) => m.to_owned(),
            None => fail!("Can't match: " + progname)
        }
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
fn test_arg_parse() -> () {
    let x = read_opts([~"redo", ~"--shuffle", ~"foo", ~"-x", ~"bar"]);
    assert!(x.eq(&Ok((~"redo", vec!(~"REDO_SHUFFLE", ~"REDO_XTRACE"), vec!(~"foo", ~"bar")))));
}
