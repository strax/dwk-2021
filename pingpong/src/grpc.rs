use tonic::{transport::Server, Request, Response, Status};
use std::sync::{Arc, atomic::Ordering};
use crate::app::State;

use proto::pingpong_service_server as server;

pub use server::{PingpongService, PingpongServiceServer};
use proto::{Stats};

mod proto {
    tonic::include_proto!("pingpong");
}

pub struct Endpoint {
    state: Arc<State>
}

impl Endpoint {
    pub fn new(state: Arc<State>) -> Self {
        Endpoint { state }
    }
}

#[tonic::async_trait]
impl PingpongService for Endpoint {
    async fn get_stats(&self, _: Request<()>) -> Result<Response<Stats>, Status> {
        Ok(Response::new(Stats { pings: self.state.counter.load(Ordering::SeqCst) }))
    }
}