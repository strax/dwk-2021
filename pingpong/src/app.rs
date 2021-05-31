use std::sync::{Arc};
use std::error::Error;

use sqlx::PgPool;
use warp::Filter;
use serde::Serialize;

#[derive(Debug)]
pub struct App {
    pub db: PgPool
}

impl App {
    pub fn new(db: PgPool) -> App {
        App { db }
    }

    pub async fn count_pings(&self) -> Result<i64, Box<dyn Error>> {
        let value = sqlx::query_scalar("SELECT COUNT(*) FROM pings").fetch_one(&self.db).await?;
        Ok(value)
    }

    pub async fn create_ping(&self) -> Result<(), Box<dyn Error>> {
        sqlx::query("INSERT INTO pings VALUES (gen_random_uuid(), NOW())").execute(&self.db).await?;
        Ok(())
    }
}

#[derive(Serialize)]
struct StatsResponse {
    pings: i64
}

pub fn with_state(state: Arc<App>) -> warp::filters::BoxedFilter<(Arc<App>, )> {
    warp::any().map(move || state.clone()).boxed()
}

pub mod routes {
    use super::*;

    #[inline]
    pub fn ping(state: Arc<App>) -> warp::filters::BoxedFilter<(impl warp::Reply,)> {
        warp::path::end()
            .and(with_state(state))
            .and(warp::get())
            .and_then(handlers::ping)
            .boxed()
    }

    #[inline]
    pub fn stats(state: Arc<App>) -> warp::filters::BoxedFilter<(impl warp::Reply,)> {
        warp::path("stats")
            .and(with_state(state))
            .and(warp::get())
            .and_then(handlers::stats)
            .boxed()
    }
}

mod handlers {
    use super::*;
    use std::convert::Infallible;

    pub async fn ping(state: Arc<App>) -> Result<impl warp::Reply, Infallible> {
        state.create_ping().await.unwrap();
        let counter = state.count_pings().await.unwrap();
        Ok(format!("pong {}", counter))
    }

    pub async fn stats(state: Arc<App>) -> Result<impl warp::Reply, Infallible> {
        let pings = state.count_pings().await.unwrap();
        Ok(warp::reply::json(&StatsResponse { pings }))
    }
}
