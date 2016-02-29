use msgpack::Decoder;
#[cfg(test)]
use msgpack::Encoder;
use msgpack::decode::Error as DecodeError;
use msgpack::decode::Error::InvalidMarkerRead;
use rmp::decode::ReadError::UnexpectedEOF;
use rustc_serialize::Decodable;
use rustc_serialize::Encodable;
use std::env::home_dir;
use std::io::Read;
use std::iter::Iterator;
use std::path::PathBuf;
use std::result::Result;
use std::marker::PhantomData;

#[derive(RustcEncodable, RustcDecodable, PartialEq, Debug, Copy, Clone)]
pub enum Operation {
    RedoIfChange,
    RedoIfCreate,
    Redo,
}

impl Operation {
    pub fn from_str(name: &str) -> Option<Operation> {
        match name {
            "redo-ifchange" => Some(Operation::RedoIfChange),
            "redo-ifcreate" => Some(Operation::RedoIfCreate),
            "redo" => Some(Operation::Redo),
            _ => None
        }
    }
}

#[test]
fn operation_from_string(){
    assert_eq!(Operation::RedoIfChange, Operation::from_str("redo-ifchange").unwrap());
    assert_eq!(Operation::RedoIfCreate, Operation::from_str("redo-ifcreate").unwrap());
    assert_eq!(Operation::Redo, Operation::from_str("redo").unwrap());
    assert_eq!(None, Operation::from_str("redod"));
}

#[derive(RustcEncodable, RustcDecodable, PartialEq, Debug, Clone)]
pub struct Request {
    pub id: u32,
    pub op: Operation,
    pub target: PathBuf,
}

impl Request {
    pub fn new(id: u32, op: Operation, target: PathBuf) -> Request {
        Request{ id: id, op: op, target: target}
    }
}

#[derive(RustcEncodable, RustcDecodable, PartialEq, Debug, Clone)]
pub struct Reply {
    pub id: u32,
    pub target: PathBuf,
}

impl Reply {
    pub fn new(id: u32, target: PathBuf) -> Reply { Reply{id: id, target: target} }
}

pub fn get_sock_path() -> PathBuf {
    let mut sock_path = home_dir().expect("No HOME directory");
    sock_path.push(".redo");
    sock_path.push("redod.sock");
    sock_path
}

pub struct StreamDecoder<'a, T: Decodable, R: 'a + Read> {
    decoder: Decoder<&'a mut R>,
    data_type: PhantomData<T>,
}

impl<'a, T: Decodable, R: Read> StreamDecoder<'a, T, R> {
    pub fn new(reader: &'a mut R) -> Self {
        StreamDecoder {
            decoder: Decoder::new(reader),
            data_type: PhantomData,
        }
    }
}

impl<'a, T: Decodable, R: Read> Iterator for StreamDecoder<'a, T, R> {
    type Item = Result<T, DecodeError>;
    fn next(&mut self) -> Option<Self::Item> {
        match Decodable::decode(&mut self.decoder) {
            Ok(v) => Some(Ok(v)),
            Err(InvalidMarkerRead(UnexpectedEOF)) => None,
            Err(e) => Some(Err(e)),
        }
    }
}

#[test]
fn stream_decoder() {
    fn prop<T: Encodable + Decodable + Eq>(xs: Vec<T>) -> bool {
        let mut bytes = Vec::new();
        {
            let mut enc = Encoder::new(&mut bytes);
            for x in &xs {
                x.encode(&mut enc).unwrap();
            }
        }
        let mut bytes = ::std::io::Cursor::new(bytes);
        let sd = StreamDecoder::new(&mut bytes);
        let mut res = Vec::new();
        for r in sd {
            let r = r.unwrap();
            res.push(r);
        }
        xs == res
    }
    assert!(prop(Vec::<bool>::new()));
    assert!(prop(vec!(true)));
    assert!(prop(vec!(false)));
    assert!(prop(vec!(false,false)));
    assert!(prop(vec!(true,false)));
    assert!(prop(vec!(false,true)));
    assert!(prop(vec!(true,true)));
}
