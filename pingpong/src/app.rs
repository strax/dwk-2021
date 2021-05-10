use warp::Filter;
use std::sync::{Arc, atomic::{Ordering, AtomicU64}};
use tokio::fs;
use tracing::error;

#[derive(Debug, Default)]
pub struct State {
    counter: AtomicU64
}

pub fn with_state(state: Arc<State>) -> impl Filter<Extract = (Arc<State>,), Error = std::convert::Infallible> + Clone {
    warp::any().map(move || state.clone())
}

pub mod routes {
    use super::*;

    #[inline]
    pub fn ping(state: Arc<State>) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
        warp::path::end()
            .and(with_state(state))
            .and(warp::get())
            .and_then(handlers::ping)
    }
}

mod handlers {
    use super::*;
    use std::convert::Infallible;

    pub async fn ping(state: Arc<State>) -> Result<impl warp::Reply, Infallible> {
        let counter = state.counter.fetch_add(1, Ordering::SeqCst);
        tokio::spawn(async move {
            if let Err(err) = fs::write("/mnt/storage/counter", counter.to_string()).await {
                error!("{}", err.to_string());
            }
        });
        Ok(format!("pong {}", counter))
    }
}