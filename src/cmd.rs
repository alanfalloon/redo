use std::sync::Arc;
use super::{State, Fs};
use std::path::{Path,PathBuf};

#[derive(PartialEq, Eq, Debug, Copy, Clone)]
pub enum CmdError {
}

pub type CmdResult = Result<(PathBuf, State), CmdError>;

pub fn redo<T: Fs>(_world: Arc<T>, _target: &Path) -> CmdResult {
    unimplemented!();
}

pub fn redo_ifchange<T: Fs>(_world: Arc<T>, _target: &Path) -> CmdResult {
    unimplemented!();
}

pub fn redo_ifcreate<T: Fs>(_world: Arc<T>, _target: &Path) -> CmdResult {
    unimplemented!();
}
