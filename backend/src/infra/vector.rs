use crate::config::config::VectorConfig;
use crate::datasource::vectordao::dao::VectorDao;

use anyhow::Result;
use qdrant_client::Qdrant;

pub async fn new_vector(c: VectorConfig) -> Result<VectorDao> {
    let url = format!("http://{}", c.host);
    let client = Qdrant::from_url(&url)
        .api_key(c.api_key) // 在这里填入你的 API Key
        .build()?;

    Ok(VectorDao::new(client))
}
