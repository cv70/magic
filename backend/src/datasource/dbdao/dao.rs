// Database Data Access Object
use sqlx::PgPool;

pub struct DBDao {
    pub db: PgPool,
}

impl DBDao {
    pub fn new(db: PgPool) -> Self {
        DBDao { db }
    }
}
