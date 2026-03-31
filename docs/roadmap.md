# Roadmap — Go Auth Frameworks Comparison

## Mục tiêu

So sánh 4 cách tiếp cận Go web: **Gin**, **Fiber**, **net/http (stdlib)**, **Echo** — cùng implement hệ thống authentication với Access Token + Refresh Token + Sessions.

## Tài liệu tham khảo

| File | Nội dung |
|------|---------|
| `docs/erd.md` | Database schema (users + sessions), indexes, token strategy |
| `docs/api-spec.md` | 12 routes, input/output/method, middleware flow |
| `docs/conventions.md` | Quyết định chung: cookie, error messages, redirect codes, template path |

## PORT Mapping

| Approach | Local Dev | Production (Fly.io) |
|----------|----------|-------------------|
| Gin | 8081 | 8080 |
| Fiber | 8082 | 8080 |
| net/http | 8083 | 8080 |
| Echo | 8084 | 8080 |

## Lộ trình (8 Phases — 16 Prompts)

| Phase | Prompt | Nội dung |
|-------|--------|---------|
| 1 — Foundation | 01, 02 | Skeleton + Templates + Seed DB |
| 2 — Gin | 03, 04, 05 | Config/Data → Logic → Deploy |
| 3 — Fiber | 06, 07 | All layers → Deploy |
| 4 — net/http | 08, 09, 10 | Config/Data → Logic → Deploy |
| 5 — Echo | 11, 12 | All layers → Deploy |
| 6 — CI/CD | 13 | GitHub Actions + Makefile |
| 7 — OAuth | 14, 15 | Google OAuth Gin → port 3 còn lại |
| 8 — Tests | 16 | Unit Tests tất cả |

> **net/http (stdlib)** được chia 3 prompt (giống Gin) vì phải tự viết template renderer, middleware chaining, form parsing — phức tạp hơn các framework có sẵn.

## Quy tắc chuyển Phase

> Build pass + tất cả unit test pass + E2E flow chạy được → mới chuyển sang phase tiếp theo.
