// Application configuration
use anyhow::Result;
use serde::{Deserialize, Serialize};
use std::fs;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct AppConfig {
    pub server: ServerConfig,
    pub database: DatabaseConfig,
    pub redis: RedisConfig,
    pub vector: VectorConfig,
    pub scylla: ScyllaConfig,
    pub llm: LLMConfig,
    pub text_embedding: EmbeddingConfig,
    pub browser: BrowserConfig,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct ServerConfig {
    pub env: String,
    pub host: String,
    pub port: u16,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct DatabaseConfig {
    pub host: String,
    pub user: String,
    pub pass: String,
    pub db_name: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct RedisConfig {
    pub host: String,
    pub port: u16,
    pub pass: Option<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct VectorConfig {
    pub host: String,
    pub api_key: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct ScyllaConfig {
    pub host: String,
    pub user: String,
    pub pass: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct LLMConfig {
    pub base_url: String,
    pub model: String,
    pub timeout: u64,
    pub max_tokens: Option<u32>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct EmbeddingConfig {
    pub base_url: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BrowserConfig {
    pub url: String,
}

impl AppConfig {
    pub fn load_from_yaml_file(path: &str) -> Result<Self> {
        let content = fs::read_to_string(path)
            .map_err(|e| anyhow::anyhow!("failed to read config file '{}': {}", path, e))?;
        let cfg = serde_yaml::from_str::<AppConfig>(&content)
            .map_err(|e| anyhow::anyhow!("failed to parse yaml config '{}': {}", path, e))?;
        Ok(cfg)
    }
}

#[cfg(test)]
mod tests {
    use super::AppConfig;
    use std::fs;
    use std::time::{SystemTime, UNIX_EPOCH};

    #[test]
    fn load_from_yaml_file_success() {
        let unique = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .expect("clock went backwards")
            .as_nanos();
        let path = std::env::temp_dir().join(format!("magic-backend-config-{unique}.yaml"));
        let yaml = r#"
server:
  env: dev
  host: 127.0.0.1
  port: 9999
database:
  host: db
  user: u
  pass: p
  db_name: n
redis:
  host: redis
  port: 6379
  pass: null
vector:
  host: vector
  api_key: k
scylla:
  host: scylla
  user: su
  pass: sp
llm:
  base_url: http://llm
  model: qwen2
  timeout: 30
  max_tokens: 2048
text_embedding:
  base_url: http://embed
browser:
  url: http://browser
"#;
        fs::write(&path, yaml).expect("write temp yaml");

        let cfg =
            AppConfig::load_from_yaml_file(path.to_str().expect("utf8 temp path")).expect("load");
        assert_eq!(cfg.server.port, 9999);
        assert_eq!(cfg.database.host, "db");

        fs::remove_file(path).expect("remove temp yaml");
    }

    #[test]
    fn load_from_yaml_file_missing_returns_error() {
        let err = AppConfig::load_from_yaml_file("/tmp/definitely-not-exist-magic-config.yaml")
            .expect_err("should fail for missing file");
        assert!(
            err.to_string().contains("failed to read config file"),
            "unexpected error: {err}"
        );
    }
}
