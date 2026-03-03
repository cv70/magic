use crate::config::config::ScyllaConfig;
use crate::datasource::scylladao::dao::ScyllaDao;

use anyhow::Result;
use scylla::SessionBuilder;

pub async fn new_scylla(c: ScyllaConfig) -> Result<ScyllaDao> {
    let scylla_user = c.user;
    let scylla_pass = c.pass;
    let scylla_host = c.host;

    let sess = SessionBuilder::new()
        .known_node(scylla_host)
        .user(scylla_user, scylla_pass)
        .build()
        .await?;
    Ok(ScyllaDao::new(sess))
}
