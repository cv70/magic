# Mandatory YAML Config Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** backend 仅支持 `--config <yaml>` 配置来源，移除程序内环境变量解析。

**Architecture:** 启动参数强制要求 `--config`，配置统一通过 `AppConfig::load_from_yaml_file` 读取并反序列化。删除默认 `load()` 环境变量分支，文档与示例同步更新。

**Tech Stack:** Rust, Axum, serde_yaml, cargo test

---

### Task 1: 强制 `--config` 参数
**Files:**
- Modify: `backend/src/main.rs`
- Test: `backend/src/main.rs`

### Task 2: 删除环境变量配置解析
**Files:**
- Modify: `backend/src/config/config.rs`
- Test: `backend/src/config/config.rs`

### Task 3: 文档和示例配置
**Files:**
- Create: `backend/config.yaml`
- Modify: `README.md`

### Task 4: 验证
**Files:**
- Verify: `backend` 全量测试
