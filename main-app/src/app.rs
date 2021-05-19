use chrono::prelude::*;
use serde::Serialize;
use tokio::fs;
use uuid::Uuid;
use warp::Filter;
use futures::TryFutureExt;
use tracing::error;
use tonic::transport::Channel;
use tonic::Request;
use tokio::sync::Mutex;
use std::sync::Arc;

#[derive(Debug, Serialize)]
struct ResponseEnvelope {
    ts: DateTime<Utc>,
    request_id: Uuid,
    pings: u64
}

pub mod routes {
    use super::*;
    use std::error::Error;

    #[derive(Debug)]
    struct InternalServerError;
    impl warp::reject::Reject for InternalServerError {}

    pub fn status(pingpong_client: Arc<Mutex<crate::pingpong::Client<Channel>>>) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
        warp::get().and_then(move || handlers::status(pingpong_client.clone()))
    }
}

pub mod handlers {
    use super::*;

    pub async fn status(pingpong_client: Arc<Mutex<crate::pingpong::Client<Channel>>>) -> Result<impl warp::Reply, std::convert::Infallible> {
        let request_id = Uuid::new_v4();

        let response = async {
            let mut client = pingpong_client.lock().await;
            client.get_stats(Request::new(())).await.unwrap().into_inner()
        }.await;

        let response = format!("{}: {}\nPing / Pongs: {}", Utc::now(), request_id, response.pings);
        Ok(response)
    }
}
