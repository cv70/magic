# Backend Config CLI Design

**Goal:** backend 在运行时支持 `--config /path/to/file`，并从 YAML 文件加载配置。

## Scope
- 新增命令行参数 `--config`。
- `--config` 存在时从 YAML 读取 `AppConfig`。
- 未提供 `--config` 时保持现有默认配置加载逻辑。
- 提供单元测试覆盖 YAML 加载成功/失败场景。

## Design
1. 参数解析
- 在 `main.rs` 内新增轻量参数解析函数，识别 `--config <path>`。
- 不引入新 CLI 框架，保持最小改动。

2. 配置加载
- 在 `AppConfig` 新增 `load_from_yaml_file(path: &str) -> Result<Self>`。
- 使用 `std::fs::read_to_string` + `serde_yaml::from_str`。
- 错误信息包含配置文件路径与失败原因。

3. 启动逻辑
- `main` 中根据是否传入 `--config` 分支：
  - 有参数：`AppConfig::load_from_yaml_file`。
  - 无参数：`AppConfig::load()`。

4. 验证
- 单元测试：
  - 正常 YAML 可加载并断言关键字段。
  - 不存在文件返回错误。
