# Roadmap — Go Auth Frameworks Comparison

## Mục tiêu

So sánh 4 cách tiếp cận Go web: **Gin**, **Fiber**, **net/http (stdlib)**, **Echo** — cùng implement một hệ thống authentication hoàn chỉnh với Access Token + Refresh Token + Sessions.

## Kiến trúc

- **Monorepo**: 1 repo, 4 apps, dùng chung templates + DB
- **Database**: MongoDB Atlas (2 collections: `users`, `sessions`)
- **Auth**: Access Token (JWT, cookie, 30 phút) + Refresh Token (cookie + DB, 7 ngày)
- **Auto-refresh**: Middleware tự cấp lại access token khi hết hạn
- **Deploy**: Fly.io — mỗi app 1 service riêng
- **Go version**: 1.22+ (bắt buộc cho net/http enhanced ServeMux)

## PORT Mapping

| Approach | Local Dev | Production (Fly.io) |
|----------|----------|-------------------|
| Gin | 8081 | 8080 |
| Fiber | 8082 | 8080 |
| net/http | 8083 | 8080 |
| Echo | 8084 | 8080 |

## Pages

| Trang | Path | Mô tả |
|-------|------|-------|
| Login | `/login` | Đăng nhập (trang mặc định) |
| Signup | `/signup` | Đăng ký tài khoản mới |
| Home | `/home` | Thông tin tài khoản + đăng xuất (auth required) |
| Dashboard | `/dashboard` | Quản trị viên (auth + admin only) |

## Tài liệu tham khảo

| File | Nội dung |
|------|---------|
| `docs/erd.md` | Database schema, indexes, token strategy |
| `docs/api-spec.md` | Endpoints, input/output, middleware flow |

## Lộ trình (8 Phases)

| Phase | Prompt | Nội dung |
|-------|--------|---------|
| 1 — Foundation | 01, 02 | Skeleton, templates, seed DB |
| 2 — Gin | 03, 04, 05 | Config → Logic → Deploy |
| 3 — Fiber | 06, 07 | All layers → Deploy |
| 4 — net/http | 08, 09 | All layers → Deploy |
| 5 — Echo | 10, 11 | All layers → Deploy |
| 6 — CI/CD | 12 | GitHub Actions, Makefile |
| 7 — OAuth | 13, 14 | Google OAuth (Gin → port 3 còn lại) |
| 8 — Tests | 15 | Unit Tests cho tất cả |

## Quy tắc chuyển Phase

> Build pass + tất cả unit test pass + E2E flow chạy được → mới chuyển sang phase tiếp theo.
