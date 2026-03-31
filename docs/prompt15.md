# Prompt 15 — Google OAuth: Fiber + net/http + Echo

## Role

Bạn là một **Go Backend Engineer** port OAuth từ Gin sang 3 approaches.

---

## Context

Prompt 14 đã implement OAuth cho Gin. Port sang Fiber, stdlib, Echo. Logic giống — chỉ đổi transport + cookie API. Conventions xem `docs/conventions.md`.

---

## Dependencies

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 14 | Gin OAuth (reference) |
| 07 | Fiber app |
| 10 | stdlib app |
| 12 | Echo app |

---

## Yêu cầu chung

1. `go.mod` thêm `golang.org/x/oauth2 v0.19.0`
2. `config.go` thêm Google fields
3. `service/oauth.go` copy từ Gin
4. `handler/oauth.go` — đổi framework API
5. Cookie **Secure = cfg.IsProduction** (tất cả: oauth_state, access_token, refresh_token)
6. Callback: tạo session + set 2 cookies + redirect theo role

### Fiber
- `GOOGLE_REDIRECT_URL=http://localhost:8082/auth/google/callback`

### net/http (stdlib)
- `GOOGLE_REDIRECT_URL=http://localhost:8083/auth/google/callback`

### Echo
- `GOOGLE_REDIRECT_URL=http://localhost:8084/auth/google/callback`

---

## Acceptance Criteria

1. `go build ./...` pass cho cả 3
2. Cookie Secure = cfg.IsProduction
3. OAuth E2E giống Gin

---

## Checklist hoàn thành

| # | Approach | Yêu cầu | Trạng thái |
|---|----------|---------|------------|
| 1 | Fiber | OAuth full + Secure=cfg.IsProduction | 🔲 |
| 2 | stdlib | OAuth full + Secure=cfg.IsProduction | 🔲 |
| 3 | Echo | OAuth full + Secure=cfg.IsProduction | 🔲 |

---

## Report

> Điền sau khi hoàn thành
