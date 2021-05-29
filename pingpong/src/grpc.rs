use tonic::{Request, Response, Status};
use std::sync::{Arc};
use crate::app::App;

use proto::pingpong_service_server as server;

pub use server::{PingpongService, PingpongServiceServer};
use proto::{Stats};

mod proto {
    tonic::include_proto!("pingpong");
}

pub struct Endpoint {
    state: Arc<App>
}

impl Endpoint {
    pub fn new(state: Arc<App>) -> Self {
        Endpoint { state }
    }
}

#[tonic::async_trait]
impl PingpongService for Endpoint {
    async fn get_stats(&self, _: Request<()>) -> Result<Response<Stats>, Status> {
        let pings = self.state.count_pings().await.unwrap() as u64;
        Ok(Response::new(Stats { pings }))
    }
}