use chrono::prelude::*;
use serde::Serialize;
use tokio::fs;
use uuid::Uuid;
use warp::Filter;
use futures::TryFutureExt;
use tracing::error;

#[derive(Debug, Serialize)]
struct ResponseEnvelope {
    ts: DateTime<Utc>,
    request_id: Uuid,
    pingpongs: u64
}

pub mod routes {
    use super::*;
    use std::error::Error;

    #[derive(Debug)]
    struct InternalServerError;
    impl warp::reject::Reject for InternalServerError {}

    #[inline]
    pub fn status() -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
        warp::get()
            .and_then(|| async {
                let request_id = Uuid::new_v4();
                let pingpongs: u64 = fs::read_to_string("/mnt/storage/counter").await?.parse()?;
                Ok::<_, Box<dyn Error>>(warp::reply::json(&ResponseEnvelope { ts: Utc::now(), request_id, pingpongs }))
            }.map_err(|err| {
                error!("{}", err);
                warp::reject::custom(InternalServerError)
            }))
    }
}