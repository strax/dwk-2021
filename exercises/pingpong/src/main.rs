mod app;
mod grpc;
mod config;

use std::str::FromStr;
use std::net::{SocketAddr, IpAddr};
use std::sync::{Arc};

use futures::future::select;
use tokio::pin;
use tokio::signal::unix::{signal, SignalKind};
use tokio::sync::watch;
use tracing::{info, Level};
use tracing_subscriber::{EnvFilter};
use warp::Filter;
use tonic::transport::Server as GrpcServer;
use sqlx::postgres::{PgPoolOptions};

use config::Config;

async fn signal_handler(notification: watch::Sender<()>) {
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
    tracing_subscriber::fmt()
        .with_env_filter(EnvFilter::from_env("LOG_LEVEL").add_directive(Level::INFO.into()))
        .without_time()
        .init();

    let (tx, mut rx) = watch::channel(());
    tokio::task::spawn(signal_handler(tx));

    let config = Config::from_env().unwrap();
    let db = PgPoolOptions::new().max_connections(5).connect_lazy_with(config.database);

    let state = Arc::new(app::App::new(db));

    let routes = app::routes::ping(state.clone())
        .or(app::routes::stats(state.clone()))
        .with(warp::trace::request());
    let routes = match config.base_path {
        None => {
            info!("Base path = /");
            routes.boxed()
        },
        Some(base_path) => {
            info!("Base path = /{}", &base_path);
            warp::path(base_path).and(routes).boxed()
        }
    };

    let mut rx2 = rx.clone();
    let (http_addr, server) = warp::serve(app::routes::health(state.clone()).or(routes))
        .bind_with_graceful_shutdown(SocketAddr::new(IpAddr::from_str("::").unwrap(), config.port), async move {
            rx2.changed().await.ok();
        });
    let http = tokio::task::spawn(server);
    info!(addr = %http_addr, "http.up");

    let grpc = GrpcServer::builder()
        .add_service(grpc::PingpongServiceServer::new(grpc::PingpongServiceHandler::new(state.clone())))
        .serve_with_shutdown(SocketAddr::new(IpAddr::from_str("::").unwrap(), 50051), async move {
            rx.changed().await.ok();
        });
    info!("grpc.up");
    let _ = futures::future::join(http, grpc).await;
    info!("shutdown");
}
