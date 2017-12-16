use std::str::FromStr;
use hyper::header::{Authorization, Bearer, Header, Raw};
use rocket::{Outcome, Route, State};
use rocket::http::Status;
use rocket::request::{self, Request, FromRequest};
use rocket_contrib::Json;
use failure::Error;

use redirs::{Redir, RedirStore};

pub fn routes() -> Vec<Route> {
    routes![get_all, get_single, create]
}

#[get("/")]
fn get_all(_token: ApiToken, store: State<Box<RedirStore>>) -> Json<Vec<Redir>> {
    Json(store.all_in_cache())
}

#[put("/", format = "application/json", data = "<body>")]
fn create(_token: ApiToken, body: Json<Redir>, store: State<Box<RedirStore>>) -> Result<(), Error> {
    store.put(body.into_inner())
}

#[get("/<ident>")]
fn get_single(
    _token: ApiToken,
    ident: String,
    store: State<Box<RedirStore>>,
) -> Option<Json<Redir>> {
    store.get(&ident).map(|v| Json(v))
}

pub struct ApiToken(String);

fn is_valid(given: &str, expected: &ApiToken) -> bool {
    given == expected.0
}

impl<'a, 'r> FromRequest<'a, 'r> for ApiToken {
    type Error = ();

    fn from_request(request: &'a Request<'r>) -> request::Outcome<ApiToken, ()> {
        let auth: Authorization<Bearer> = match request.headers().get_one("Authorization") {
            Some(h) => Authorization::parse_header(&Raw::from(h)).unwrap(),
            None => return Outcome::Failure((Status::Unauthorized, ())),
        };

        let api_token = request.guard::<State<ApiToken>>().unwrap().inner();
        if !is_valid(&auth.token, api_token) {
            Outcome::Failure((Status::Unauthorized, ()))
        } else {
            Outcome::Success(ApiToken(auth.token.clone()))
        }
    }
}

impl FromStr for ApiToken {
    type Err = ();
    fn from_str(s: &str) -> Result<Self, Self::Err> {
        Ok(ApiToken(s.to_owned()))
    }
}
