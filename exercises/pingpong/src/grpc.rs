use tonic::{Request, Response, Status};
use std::sync::{Arc};
use crate::app::App;

use proto::pingpong_service_server as server;

pub use server::{PingpongService, PingpongServiceServer};
use proto::{Stats};

mod proto {
    tonic::include_proto!("pingpong");
}

pub struct PingpongServiceHandler {
    state: Arc<App>
}

impl PingpongServiceHandler {
    pub fn new(state: Arc<App>) -> Self {
        PingpongServiceHandler { state }
    }
}

#[tonic::async_trait]
impl PingpongService for PingpongServiceHandler {
    async fn get_stats(&self, _: Request<()>) -> Result<Response<Stats>, Status> {
        match self.state.count_pings().await {
            Err(_) => Err(Status::unavailable("Failed to fetch pings")),
            Ok(pings) => Ok(Response::new(Stats { pings }))
        }
    }
}
