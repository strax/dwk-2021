use tokio::time::{sleep, Duration};
use uuid::Uuid;

#[tokio::main]
pub async fn main() {
    loop {
        let uuid = Uuid::new_v4();
        println!("{}", uuid);
        sleep(Duration::from_secs(5)).await;
    }
}
