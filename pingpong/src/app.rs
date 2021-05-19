use warp::Filter;
use std::sync::{Arc, atomic::{Ordering, AtomicU64}};
use serde::Serialize;

#[derive(Debug, Default)]
pub struct State {
    pub counter: AtomicU64
}

#[derive(Serialize)]
struct StatsResponse {
    pings: u64
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

    #[inline]
    pub fn stats(state: Arc<State>) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
        warp::path("stats")
            .and(with_state(state))
            .and(warp::get())
            .and_then(handlers::stats)
    }
}

mod handlers {
    use super::*;
    use std::convert::Infallible;

    pub async fn ping(state: Arc<State>) -> Result<impl warp::Reply, Infallible> {
        let counter = state.counter.fetch_add(1, Ordering::SeqCst);
        Ok(format!("pong {}", counter))
    }

    pub async fn stats(state: Arc<State>) -> Result<impl warp::Reply, Infallible> {
        Ok(warp::reply::json(&StatsResponse { pings: state.counter.load(Ordering::SeqCst) }))
    }
}