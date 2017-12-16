#![feature(use_extern_macros)]
#![feature(plugin)]
#![plugin(rocket_codegen)]

extern crate clap;
extern crate rocket;
extern crate rocket_contrib;
extern crate serde;
extern crate hyper;
#[macro_use]
extern crate serde_derive;
extern crate serde_json;
#[macro_use]
extern crate failure;
extern crate rand;

mod adm;
mod redirs;

use rocket::request::State;
use rocket::response::Redirect;
use rand::Rng;

use redirs::RedirStore;
use std::str::FromStr;
use adm::ApiToken;

#[get("/")]
fn index() -> &'static str {
    "
    WHAT IS THIS?
    a really simple redirection service

    WHAT CAN IT DO?
    talk to /adm to find out all the routes we currently know about
    "
}

#[get("/<key>")]
fn redir(key: String, store: State<Box<RedirStore>>) -> Option<Redirect> {
    store.get(&key).map(|v| Redirect::to(&v.to))
}

#[error(404)]
fn not_found() -> &'static str {
    "
    404 
    not found
    "
}
#[error(401)]
fn unauthorized() -> &'static str {
    "
    401 
    unauthorized
    "
}

fn rocket(linkfile: &str) -> rocket::Rocket {
    let store = Box::new(RedirStore::from_path(linkfile));
    let api_token = rand::thread_rng()
        .gen_ascii_chars()
        .take(24)
        .collect::<String>();
    println!("   Using protecting access to API at /adm");
    println!("    => api_token: {}", api_token);

    rocket::ignite()
        .manage(store)
        .manage(ApiToken::from_str(&api_token).unwrap())
        .mount("/", routes![index, redir])
        .mount("/adm", adm::routes())
        .catch(errors![not_found, unauthorized])
}

fn main() {
    let matches = clap::App::new("redirsrv")
        .version("0.1.0")
        .author("robertgzr <r@gnzler.io>")
        .about("simple redirection server")
        .arg(
            clap::Arg::with_name("LINKFILE")
                .short("f")
                .long("linkfile")
                .value_name("FILE")
                .help("file used to read/persist redirections")
                .takes_value(true),
        )
        .get_matches();

    let linkfile = matches.value_of("LINKFILE").unwrap_or("./linkfile.json");
    rocket(linkfile).launch();
}
