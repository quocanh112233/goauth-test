# Prompt 13 — GitHub Actions CI/CD + Makefile hoàn chỉnh

## Role

Bạn là một **DevOps Engineer** thành thạo GitHub Actions và Fly.io CLI.

---

## Context

Phase 6. Tất cả 4 apps hoàn chỉnh. Setup CI/CD + Makefile.

---

## Dependencies

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 05 | Gin Dockerfile + fly.toml |
| 07 | Fiber Dockerfile + fly.toml |
| 10 | stdlib Dockerfile + fly.toml |
| 12 | Echo Dockerfile + fly.toml |

---

## Yêu cầu

### 1–4. GitHub Actions Workflows

4 files: `deploy-gin.yml`, `deploy-fiber.yml`, `deploy-stdlib.yml`, `deploy-echo.yml`

- Trigger: push `main`, paths: `<approach>/**` + `shared/**`
- Token riêng: `FLY_API_TOKEN_<APPROACH>`
- `flyctl deploy --remote-only`

### 5. Makefile

```makefile
# Development
run-gin run-fiber run-stdlib run-echo run-all

# Seed
seed

# Deploy
deploy-gin deploy-fiber deploy-stdlib deploy-echo deploy-all

# Secrets
secrets-gin secrets-fiber secrets-stdlib secrets-echo secrets-all

# Test
test-gin test-fiber test-stdlib test-echo test-all test-coverage

# Utility
tidy help
```

### 6. README.md CI/CD section

---

## Acceptance Criteria

1. Push `gin/` → chỉ trigger deploy-gin
2. Push `shared/` → trigger cả 4
3. `make help` hiển thị targets
4. `make run-all` chạy 4 server song song

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | 4 workflow files | 🔲 | |
| 2 | shared/** trigger cả 4 | 🔲 | |
| 3 | Makefile đủ targets | 🔲 | |
| 4 | README CI/CD section | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
