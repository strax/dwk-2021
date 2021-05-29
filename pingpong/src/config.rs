use anyhow::{Result, Context};

use sqlx::postgres::PgConnectOptions;

#[derive(Debug)]
pub struct Config {
    pub port: u16,
    pub database: PgConnectOptions
}

fn env_var(key: &str) -> Result<String> {
    std::env::var(key).context(format!("environment variable '{}' is required but not present", key))
}

impl Config {
    pub fn from_env() -> Result<Config> {
        let port = u16::from_str_radix(&env_var("PORT")?, 10)?;
        let database = PgConnectOptions::new()
            .host(&env_var("PGHOST")?)
            .port(u16::from_str_radix(&env_var("PGPORT")?, 10)?)
            .username(&env_var("PGUSER")?)
            .password(&env_var("PGPASS")?)
            .database(&env_var("PGDATABASE")?)
            .application_name("pingpong");
        Ok(Config { port, database })
    }
}