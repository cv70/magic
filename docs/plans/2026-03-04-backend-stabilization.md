# Backend Stabilization Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Make backend compile, run, and provide usable core CRUD APIs across existing domains.

**Architecture:** Keep current Axum + domain/dao layering, fix compile/signature breaks first, then implement minimal correct CRUD behaviors with consistent API responses. Prefer small targeted edits and verify with cargo check/test after each logical batch.

**Tech Stack:** Rust 2024, Axum 0.7, Tokio, SQLx, Serde.

---

### Task 1: Establish failing baseline and inventory breakages
**Files:**
- Modify: `backend/src/**` (as needed)
- Test: `backend` cargo commands

1. Run `cargo check` and record initial failures.
2. Inspect failing files and module interfaces.
3. Group fixes into syntax/import/signature/behavior categories.

### Task 2: Restore compile pipeline
**Files:**
- Modify: `backend/src/main.rs`
- Modify: `backend/src/domain/**`
- Modify: `backend/src/datasource/**`
- Modify: `backend/src/state/**`

1. Fix syntax errors and invalid string literals.
2. Align API/domain function signatures.
3. Remove duplicate/conflicting impl blocks and broken imports.
4. Re-run `cargo check` until clean.

### Task 3: Implement minimal correct domain CRUD behavior
**Files:**
- Modify: `backend/src/domain/financing/*`
- Modify: `backend/src/domain/news/*`
- Modify: `backend/src/domain/content/*`
- Modify: `backend/src/domain/configuration/*`
- Modify: `backend/src/domain/ai_generation/*`
- Modify: `backend/src/domain/publishing/*`
- Modify: `backend/src/domain/scheduling/*`
- Modify: `backend/src/domain/identity/*`

1. Ensure each domain method delegates correctly to DAO and returns meaningful Result.
2. Replace placeholder `0`/empty returns on DB failure with explicit errors.
3. Keep API response schema stable with current structures.

### Task 4: Add focused tests for corrected behavior
**Files:**
- Create/Modify: `backend/tests/*` and/or unit tests in modules

1. Add tests for mapping/error behaviors where feasible without full DB integration.
2. Verify fail/pass cycle for new tests.

### Task 5: Verification before completion
**Files:**
- N/A

1. Run `cargo fmt --check`.
2. Run `cargo check`.
3. Run `cargo test`.
4. Report exact results and remaining gaps.
