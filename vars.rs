use std::from_str::from_str;
use std::os::{getenv,setenv};
use std::default::Default;
use std::char::to_lowercase;

/*
 * A trait like FromStr but a little more forgiving for bool.
 */
trait ParseEnv {
    fn parse(s: &str) -> Option<Self>;
}

impl ParseEnv for ~str {
    fn parse(s: &str) -> Option<~str> {
        Some(s.to_owned())
    }
}

impl ParseEnv for int {
    fn parse(s: &str) -> Option<int> {
        from_str(s)
    }
}

impl ParseEnv for bool {
    fn parse(s: &str) -> Option<bool> {
        let lowered: ~str = s.chars().map(to_lowercase).collect();
        match lowered.as_slice() {
            "true"|"t"|"yes"|"y"|"1" => Some(true),
            "false"|"f"|"no"|"n"|"0" => Some(false),
            _ => None
        }
    }
}

macro_rules! try_bool(
    ( $s:expr ) => ( ParseEnv::parse($s).expect(concat!("not a valid bool: ", stringify!($s))) )
)

#[test]
fn parse_bool() {
    assert_eq!(false, try_bool!("0"));
    assert_eq!(false, try_bool!("N"));
    assert_eq!(false, try_bool!("f"));
    assert_eq!(false, try_bool!("faLsE"));
    assert_eq!(false, try_bool!("false"));
    assert_eq!(false, try_bool!("nO"));
    assert_eq!(false, try_bool!("no"));
    assert_eq!(true, try_bool!("1"));
    assert_eq!(true, try_bool!("T"));
    assert_eq!(true, try_bool!("YEs"));
    assert_eq!(true, try_bool!("t"));
    assert_eq!(true, try_bool!("tRUe"));
    assert_eq!(true, try_bool!("true"));
    assert_eq!(true, try_bool!("y"));
    assert_eq!(true, try_bool!("yes"));
}

/**
 * Read the environment variable and parse it in to its appropriate
 * type.
 */
fn env<T: ParseEnv + Default>(envname: &str) -> T {
    match getenv(envname) {
        None => {
            info!("{} missing.", envname);
            Default::default()
        },
        Some(ev) => match ParseEnv::parse(ev) {
            None => {
                error!("{} invalid value '{}'!", envname, ev);
                Default::default()
            }
            Some(v) => v
        }
    }
}

/*
 * Macro so we only have to specify the list of environment
 * variables/names once. DRY-principle and all that.
 *
 * The ones before the | symbol are all taken straight from
 * environment variables, the ones after use the supplied expression.
 */
macro_rules! configvars(
    ( $($nv:ident : $tv:ty),+ | $($ne:ident : $te:ty = $e:expr),* ) => (
        #[deriving(Eq, Ord, Clone, Show)]
        #[allow(uppercase_variables)]
        pub struct Vars {
            $(pub $nv : $tv),+ ,
            $(pub $ne : $te),+
        }
        pub fn v() -> Vars {
            Vars {
                $($nv: env(concat!("REDO_", stringify!($nv)))),+ ,
                $($ne: $e),+
            }
        }
        static var_names: &'static [&'static str] = &[
            $(concat!("REDO_", stringify!($nv))),+
        ];
    );
)

configvars!(
    // These are the operational varaibles used to track targets
    STARTDIR: ~str,
    PWD: ~str,
    TARGET: ~str,
    DEPTH: int,
    // These are the option variables used to transmit command-line
    // options to sub-processes.
    COLOR: bool,
    DEBUG: int,
    DEBUG_LOCKS: bool,
    DEBUG_PIDS: bool,
    KEEP_GOING: bool,
    LOG: bool,
    ONLY_LOG: bool,
    OVERWRITE: bool,
    SHUFFLE: bool,
    VERBOSE: bool,
    WARN_STDOUT: bool,
    XTRACE: bool |
    RUNID: u64 = {
        let runid_filename = getenv("REDO_RUNID_FILE").unwrap_or(~".redo/runid");
        ::runid::read(&Path::new(runid_filename))
    }
)

#[cfg(not(test))]
pub fn set(opts: Vec<~str>) {
    for opt in opts.iter() {
        setenv(*opt, "1");
    }
}

#[test]
fn load_vars() {
    unsetall();
    setenv("REDO_DEPTH", "1");
    setenv("REDO_KEEP_GOING", "1");
    setenv("REDO_TARGET", "bin/foo/bar");
    let v = ::vars::v();
    assert_eq!(v.STARTDIR, ~"");
    assert_eq!(v.PWD, ~"");
    assert_eq!(v.TARGET, ~"bin/foo/bar");
    assert_eq!(v.DEPTH, 1);
    assert_eq!(v.KEEP_GOING, true);
    assert_eq!(v.LOG, false);
}

pub fn unsetall() {
    for varname in var_names.iter() {
        ::std::os::unsetenv(*varname);
    }
}
