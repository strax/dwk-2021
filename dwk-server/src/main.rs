use futures::future::select;
use tokio::pin;
use tokio::signal::unix::{signal, SignalKind};
use tokio::sync::oneshot;
use uuid::Uuid;
use warp::Filter;

use std::str::FromStr;
use std::net::{SocketAddr, Ipv6Addr};
use std::env;

const DEFAULT_PORT: u16 = 80;

async fn signal_handler(notification: oneshot::Sender<()>) {
    let mut sigint = signal(SignalKind::interrupt()).unwrap();
    let sigint = sigint.recv();
    let mut sigterm = signal(SignalKind::terminate()).unwrap();
    let sigterm = sigterm.recv();
    pin!(sigint);
    pin!(sigterm);
    select(sigint, sigterm).await;
    notification.send(()).unwrap();
}

#[tokio::main]
pub async fn main() {
    let port = match env::var("PORT") {
        Ok(port) => u16::from_str_radix(&port, 10).unwrap(),
        Err(_) => DEFAULT_PORT
    };
    let pid = Uuid::new_v4();
    let (tx, rx) = oneshot::channel();
    tokio::task::spawn(signal_handler(tx));

    let route = warp::get().map(move || {
        let request_id = Uuid::new_v4();
        format!("ProcessId={} RequestId={}", pid, request_id)
    });

    let (addr, server) = warp::serve(route).bind_with_graceful_shutdown(
        SocketAddr::new(Ipv6Addr::from_str("::").unwrap().into(), port),
        async {
            rx.await.ok();
        }
    );
    let t = tokio::task::spawn(server);
    println!("Server started in port {}", addr.port());
    t.await.unwrap();
}
