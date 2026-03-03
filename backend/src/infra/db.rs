use crate::config::config::DatabaseConfig;
use crate::datasource::dbdao::dao::DBDao;

use anyhow::Result;
use sqlx::{Pool, Postgres};

pub async fn new_db(c: DatabaseConfig) -> Result<DBDao> {
    let db_user = c.user;
    let db_pass = c.pass;
    let db_host = c.host;
    let db_name = c.db_name;

    // 拼接连接字符串
    let connection_string = format!("postgres://{}:{}@{}/{}", db_user, db_pass, db_host, db_name);

    let pool = Pool::<Postgres>::connect(&connection_string).await?;

    Ok(DBDao::new(pool))
}
