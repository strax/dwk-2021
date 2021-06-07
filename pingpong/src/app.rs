use std::sync::{Arc};
use std::error::Error;
use tracing::error;

use sqlx::PgPool;
use warp::Filter;
use warp::filters::BoxedFilter;
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

pub fn with_state(state: Arc<App>) -> BoxedFilter<(Arc<App>, )> {
    warp::any().map(move || state.clone()).boxed()
}

pub mod routes {
    use super::*;

    #[inline]
    pub fn ping(state: Arc<App>) -> BoxedFilter<(impl warp::Reply,)> {
        warp::path::end()
            .and(with_state(state))
            .and(warp::get())
            .and_then(handlers::ping)
            .boxed()
    }

    #[inline]
    pub fn stats(state: Arc<App>) -> BoxedFilter<(impl warp::Reply,)> {
        warp::path("stats")
            .and(with_state(state))
            .and(warp::get())
            .and_then(handlers::stats)
            .boxed()
    }

    #[inline]
    pub fn health(state: Arc<App>) -> BoxedFilter<(impl warp::Reply,)> {
        warp::path("healthz")
            .and(with_state(state))
            .and(warp::get())
            .and_then(handlers::health)
            .boxed()
    }
}

mod handlers {
    use super::*;
    use std::convert::Infallible;
    use warp::http::StatusCode;
    use warp::Reply;

    pub async fn ping(app: Arc<App>) -> Result<warp::reply::Response, Infallible> {
        if let Err(err) = app.create_ping().await {
            error!("{}", err);
            return Ok(StatusCode::INTERNAL_SERVER_ERROR.into_response())
        }
        match app.count_pings().await {
            Ok(counter) => Ok(format!("pong {}", counter).into_response()),
            Err(err) => {
                error!("{}", err);
                return Ok(StatusCode::INTERNAL_SERVER_ERROR.into_response())
            }
        }
    }

    pub async fn stats(app: Arc<App>) -> Result<warp::reply::Response, Infallible> {
        match app.count_pings().await {
            Ok(pings) => Ok(warp::reply::json(&StatsResponse { pings }).into_response()),
            Err(err) => {
                error!("{}", err);
                Ok(StatusCode::INTERNAL_SERVER_ERROR.into_response())
            }
        }
    }

    pub async fn health(app: Arc<App>) -> Result<impl warp::Reply, Infallible> {
        match app.db.acquire().await {
            Ok(_) => Ok(StatusCode::OK),
            Err(_) => Ok(StatusCode::SERVICE_UNAVAILABLE)
        }
    }
}
