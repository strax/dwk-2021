use warp::Filter;
use warp::http::StatusCode;
use std::str::FromStr;
use std::net::{SocketAddr, Ipv6Addr};
use std::env;

const DEFAULT_PORT: u16 = 80;

#[tokio::main]
pub async fn main() {
    let port = match env::var("PORT") {
        Ok(port) => u16::from_str_radix(&port, 10).unwrap(),
        Err(_) => DEFAULT_PORT
    };
    let route = warp::any().and_then(|| async {
        Err::<StatusCode, _>(warp::reject::not_found())
    });
    let server = warp::serve(route).bind(
        SocketAddr::new(Ipv6Addr::from_str("::").unwrap().into(), port)
    );
    let t = tokio::task::spawn(server);
    println!("Server started in port {}", port);
    t.await.unwrap();
}
