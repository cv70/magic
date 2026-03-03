#![allow(dead_code)]
// Vector Database Data Access Object
use qdrant_client::Qdrant;

pub struct VectorDao {
    pub db: Qdrant,
}

impl VectorDao {
    pub fn new(db: Qdrant) -> Self {
        VectorDao { db }
    }
}
