use std::path::PathBuf;
use std::fs::File;
use std::collections::BTreeMap;
use std::sync::Mutex;

use failure::Error;
use serde_json;

type RedirCache = Mutex<BTreeMap<String, Redir>>;

pub struct RedirStore {
    fpath: PathBuf,
    cache: RedirCache,
}

impl RedirStore {
    pub fn from_path(path_str: &str) -> Self {
        let p = PathBuf::from(path_str.to_owned());
        let mut c = BTreeMap::new();

        match File::open(&p) {
            Ok(fd) => {
                let list: Vec<Redir> = serde_json::from_reader(fd).unwrap();
                for r in list.iter() {
                    c.entry(r.short.to_owned()).or_insert(r.clone());
                }
            }
            Err(_) => (),
        };

        RedirStore {
            fpath: p,
            cache: RedirCache::new(c),
        }
    }

    /// Get the Redir element mapped to 'short'
    pub fn get(&self, short: &str) -> Option<Redir> {
        match self.cache.try_lock() {
            Ok(ref mutex) => {
                if let Some(v) = mutex.get(short) {
                    Some(v.clone())
                } else {
                    None
                }
            }
            Err(_) => None,
        }
    }

    /// Get all the Redirs from the cache
    pub fn all_in_cache(&self) -> Vec<Redir> {
        match self.cache.try_lock() {
            Ok(ref mutex) => mutex.values().cloned().collect::<Vec<Redir>>(),
            Err(_) => Vec::new(),
        }
    }

    /// Put a new Redir element in the file and cache
    pub fn put(&self, new: Redir) -> Result<(), Error> {
        match File::create(&self.fpath) {
            Ok(fd) => {
                let mut lock = self.cache.try_lock();
                match lock {
                    Ok(ref mut mutex) => {
                        mutex.insert(new.short.clone(), new);
                        serde_json::to_writer_pretty(fd, &mutex.values().collect::<Vec<&Redir>>())?;
                        Ok(())
                    }
                    Err(_) => Err(format_err!("unable to get lock")),
                }
            }
            Err(e) => Err(e)?,
        }
    }
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct Redir {
    pub short: String,
    pub to: String,
}

// impl Redir {
//     pub fn new(short: &str, url: &str) -> Result<Self, Error> {
//         let url = Url::from_str(url)?;

//         Ok(Redir {
//             short: short.to_owned(),
//             to: url.to_string(),
//         })
//     }
// }
