# Prompt 14 — Google OAuth: Fiber + net/http + Echo

## Role

Bạn là một **Go Backend Engineer** port OAuth từ Gin sang 3 approaches còn lại.

---

## Context

Prompt 13 đã implement Google OAuth cho Gin. Prompt này port sang Fiber, net/http stdlib, Echo.

Logic OAuth (state, exchange, upsert, session creation) **giống hệt Gin** — chỉ đổi transport layer + cookie API.

---

## Dependencies (Prompt phụ thuộc)

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 13 | Gin OAuth implementation (reference) |
| 07 | Fiber app hoàn chỉnh |
| 09 | stdlib app hoàn chỉnh |
| 11 | Echo app hoàn chỉnh |

---

## Yêu cầu chung cho mỗi approach

1. `go.mod` thêm `golang.org/x/oauth2 v0.19.0`
2. `config/config.go` thêm 3 Google fields
3. `service/oauth.go` copy từ Gin (logic giống hệt)
4. `handler/oauth.go` — đổi framework API
5. `router.go` thêm 2 routes public
6. `.env.example` — `GOOGLE_REDIRECT_URL` đúng port

Callback phải: **tạo session + set 2 cookies + redirect theo role** (giống login flow).

---

### Fiber

- Cookie: `c.Cookie(&fiber.Cookie{...})`
- Query: `c.Query("state")`, `c.Query("code")`
- Read cookie: `c.Cookies("oauth_state")`
- `.env.example`: `GOOGLE_REDIRECT_URL=http://localhost:8082/auth/google/callback`

### net/http (stdlib)

- Cookie: `http.SetCookie(w, &http.Cookie{...})`
- Query: `r.URL.Query().Get("state")`
- Read cookie: `r.Cookie("oauth_state")`
- `.env.example`: `GOOGLE_REDIRECT_URL=http://localhost:8083/auth/google/callback`

### Echo

- Cookie: `c.SetCookie(&http.Cookie{...})`
- Query: `c.QueryParam("state")`
- Read cookie: `c.Cookie("oauth_state")`
- `.env.example`: `GOOGLE_REDIRECT_URL=http://localhost:8084/auth/google/callback`

---

## Anti-Patterns (KHÔNG được làm)

❌ Không thay đổi OAuth logic — chỉ đổi transport layer
❌ Không quên tạo session + 2 cookies trong callback
❌ Không quên redirect theo role

---

## Acceptance Criteria

1. `go build ./...` pass cho cả 3
2. Google OAuth flow hoạt động giống Gin
3. Callback tạo session + set 2 cookies

---

## Checklist hoàn thành

| # | Approach | Yêu cầu | Trạng thái |
|---|----------|---------|------------|
| 1 | Fiber | go.mod + config + service + handler + router | 🔲 |
| 2 | Fiber | Callback: session + 2 cookies + redirect by role | 🔲 |
| 3 | stdlib | go.mod + config + service + handler + router | 🔲 |
| 4 | stdlib | Callback: session + 2 cookies + redirect by role | 🔲 |
| 5 | Echo | go.mod + config + service + handler + router | 🔲 |
| 6 | Echo | Callback: session + 2 cookies + redirect by role | 🔲 |

---

## Report

> Điền sau khi hoàn thành
