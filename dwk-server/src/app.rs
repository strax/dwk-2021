use chrono::prelude::*;
use serde::Serialize;
use uuid::Uuid;
use warp::Filter;

#[derive(Debug, Serialize)]
struct ResponseEnvelope {
    ts: DateTime<Utc>,
    request_id: Uuid
}

pub mod routes {
    use super::*;

    #[inline]
    pub fn status() -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
        warp::get()
            .map(|| {
                let request_id = Uuid::new_v4();
                warp::reply::json(&ResponseEnvelope { ts: Utc::now(), request_id })
            })
    }
}