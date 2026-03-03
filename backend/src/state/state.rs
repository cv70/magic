#![allow(dead_code)]
use std::sync::Arc;

use crate::datasource::{dbdao::dao::DBDao, scylladao::dao::ScyllaDao, vectordao::dao::VectorDao};
use crate::domain::{
    ai_generation::domain::AIGenerationDomain, configuration::domain::ConfigurationDomain,
    content::domain::ContentDomain, financing::domain::FinancingDomain,
    identity::domain::IdentityDomain, news::domain::NewsDomain,
    publishing::domain::PublishingDomain, scheduling::domain::SchedulingDomain,
};
use crate::infra::registry::Registry;

#[derive(Clone)]
pub struct AppState {
    pub financing_domain: FinancingDomain,
    pub news_domain: NewsDomain,
    pub content_domain: ContentDomain,
    pub configuration_domain: ConfigurationDomain,
    pub ai_generation_domain: AIGenerationDomain,
    pub publishing_domain: PublishingDomain,
    pub scheduling_domain: SchedulingDomain,
    pub identity_domain: IdentityDomain,
    pub db_dao: Arc<DBDao>,
    pub scylla_dao: Arc<ScyllaDao>,
    pub vector_dao: Arc<VectorDao>,
}

impl AppState {
    pub fn new(r: &Registry) -> Self {
        Self {
            financing_domain: FinancingDomain::new(r),
            news_domain: NewsDomain::new(r),
            content_domain: ContentDomain::new(r),
            configuration_domain: ConfigurationDomain::new(r),
            ai_generation_domain: AIGenerationDomain::new(r),
            publishing_domain: PublishingDomain::new(r),
            scheduling_domain: SchedulingDomain::new(r),
            identity_domain: IdentityDomain::new(r),
            db_dao: r.db_dao.clone(),
            scylla_dao: r.scylla_dao.clone(),
            vector_dao: r.vector_dao.clone(),
        }
    }
}
