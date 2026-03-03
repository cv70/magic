#![allow(dead_code)]
// ScyllaDB Data Access Object
use scylla::Session;

pub struct ScyllaDao {
    pub db: Session,
}

impl ScyllaDao {
    pub fn new(db: Session) -> Self {
        ScyllaDao { db }
    }
}
