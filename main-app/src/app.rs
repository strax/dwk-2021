use chrono::prelude::*;
use uuid::Uuid;
use warp::Filter;
use tonic::transport::Channel;
use tonic::Request;
use tokio::sync::Mutex;
use std::sync::Arc;

pub struct Ctx {
    pingpong_client: Mutex<crate::pingpong::Client<Channel>>,
    message: String
}

impl Ctx {
    pub fn new(pingpong_client: crate::pingpong::Client<Channel>, message: String) -> Ctx {
        Ctx {
            pingpong_client: Mutex::new(pingpong_client),
            message
        }
    }
}

pub mod routes {
    use super::*;

    #[derive(Debug)]
    struct InternalServerError;
    impl warp::reject::Reject for InternalServerError {}

    pub fn status(ctx: Arc<Ctx>) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
        warp::get().and_then(move || handlers::status(ctx.clone()))
    }
}

pub mod handlers {
    use super::*;

    pub async fn status(ctx: Arc<Ctx>) -> Result<impl warp::Reply, std::convert::Infallible> {
        let request_id = Uuid::new_v4();

        let response = async {
            let mut client = ctx.pingpong_client.lock().await;
            client.get_stats(Request::new(())).await.unwrap().into_inner()
        }.await;

        let response = format!(
            "{}\n{}: {}\nPing / Pongs: {}",
            ctx.message,
            Utc::now(),
            request_id,
            response.pings
        );
        Ok(response)
    }
}
