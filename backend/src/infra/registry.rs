// Registry module for managing application infrastructure
use std::sync::Arc;

use super::{db, scylla, vector};
use crate::config::config::AppConfig;
use crate::datasource::{dbdao::dao::DBDao, scylladao::dao::ScyllaDao, vectordao::dao::VectorDao};

use anyhow::Result;

pub struct Registry {
    pub db_dao: Arc<DBDao>,
    pub scylla_dao: Arc<ScyllaDao>,
    pub vector_dao: Arc<VectorDao>,
}

impl Registry {
    pub async fn new(c: &AppConfig) -> Result<Self> {
        let db = db::new_db(c.database.clone()).await?;
        let scylla = scylla::new_scylla(c.scylla.clone()).await?;
        let vector = vector::new_vector(c.vector.clone()).await?;
        Ok(Registry {
            db_dao: Arc::new(db),
            scylla_dao: Arc::new(scylla),
            vector_dao: Arc::new(vector),
        })
    }
}
