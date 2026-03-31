# ERD — Database Schema

## Database: MongoDB Atlas (`goauth`)

---

## 1. Collection: `users`

| Field | Type | Constraint | Mô tả |
|-------|------|-----------|-------|
| `_id` | ObjectID | PK, auto | MongoDB tự sinh |
| `name` | string | required | Tên hiển thị (username) |
| `email` | string | required, unique | Email đăng nhập |
| `phone` | string | required, unique | Số điện thoại |
| `password` | string | required* | Bcrypt hash, cost=12. (*Rỗng khi provider="google") |
| `role` | string | required, default="user" | `"admin"` hoặc `"user"` |
| `provider` | string | required, default="local" | `"local"` hoặc `"google"` (Phase 7) |
| `google_id` | string | optional | Google OAuth ID (Phase 7, chỉ khi provider="google") |
| `created_at` | DateTime | auto | Thời điểm tạo |
| `updated_at` | DateTime | auto | Thời điểm cập nhật cuối |

**Indexes:**
```javascript
db.users.createIndex({ "email": 1 }, { unique: true })
db.users.createIndex({ "phone": 1 }, { unique: true })
db.users.createIndex({ "google_id": 1 }, { sparse: true })  // Phase 7
```

> **Lưu ý**: `provider` và `google_id` chỉ được sử dụng ở Phase 7 (Google OAuth). Ở các Phase trước, `provider` luôn = `"local"` và `google_id` = `""`.

---

## 2. Collection: `sessions`

| Field | Type | Constraint | Mô tả |
|-------|------|-----------|-------|
| `_id` | ObjectID | PK, auto | MongoDB tự sinh |
| `user_id` | ObjectID | required | Reference tới `users._id` |
| `refresh_token` | string | required | Token random (32 bytes hex) |
| `expired_at` | DateTime | required | Thời điểm hết hạn (created + 7 ngày) |
| `created_at` | DateTime | auto | Thời điểm tạo session |
| `updated_at` | DateTime | auto | Thời điểm cập nhật |

**Indexes:**
```javascript
// TTL index — MongoDB tự động xóa document khi expired_at đã qua
db.sessions.createIndex({ "expired_at": 1 }, { expireAfterSeconds: 0 })

// Lookup by user_id (cho logout xóa tất cả sessions của user)
db.sessions.createIndex({ "user_id": 1 })

// Lookup by refresh_token (cho middleware auto-refresh)
db.sessions.createIndex({ "refresh_token": 1 })
```

> **TTL Index**: MongoDB background thread kiểm tra mỗi 60 giây và xóa documents có `expired_at` <= thời điểm hiện tại. Không cần cron job hay cleanup thủ công.

---

## Sơ đồ quan hệ

```
┌─────────────────────────────────┐
│            users                │
├─────────────────────────────────┤
│ _id          ObjectID    PK     │
│ name         string      req    │
│ email        string      unique │
│ phone        string      unique │
│ password     string      req    │
│ role         string      req    │
│ provider     string      req    │
│ google_id    string      opt    │
│ created_at   DateTime    auto   │
│ updated_at   DateTime    auto   │
└──────────────┬──────────────────┘
               │
               │  1 : N
               │
┌──────────────┴──────────────────┐
│           sessions              │
├─────────────────────────────────┤
│ _id            ObjectID  PK     │
│ user_id        ObjectID  FK→usr │
│ refresh_token  string    req    │
│ expired_at     DateTime  TTL    │
│ created_at     DateTime  auto   │
│ updated_at     DateTime  auto   │
└─────────────────────────────────┘
```

**Quan hệ**: 1 User → N Sessions (nhiều thiết bị đăng nhập đồng thời)

---

## Auth Token Strategy

```
┌──────────────────────────────────────────────────────────────┐
│                    CLIENT (Browser)                          │
│                                                              │
│  Cookie: access_token   (HttpOnly, Secure, 30 min)  ───┐    │
│  Cookie: refresh_token  (HttpOnly, Secure, 7 days)  ───┤    │
│                                                         │    │
└─────────────────────────────────────────────────────────┼────┘
                                                          │
                                                          ▼
┌──────────────────────────────────────────────────────────────┐
│                    SERVER                                    │
│                                                              │
│  access_token  → JWT (user_id, email, role, exp)            │
│                  Stateless, verify bằng JWT_SECRET           │
│                                                              │
│  refresh_token → Lookup trong MongoDB sessions collection   │
│                  Nếu tìm thấy + chưa hết hạn → cấp         │
│                  access_token mới                            │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

| Token | Lưu trữ | Lifetime | Mục đích |
|-------|---------|----------|---------|
| Access Token (JWT) | Cookie `access_token` | 30 phút | Xác thực mỗi request, chứa user info |
| Refresh Token | Cookie `refresh_token` + MongoDB `sessions` | 7 ngày | Tự động cấp access token mới khi hết hạn |

> **Tại sao refresh token cần lưu cả cookie và DB?**
> - Cookie: để client gửi lên server (middleware tự đọc)
> - DB: để server validate + cho phép revoke (logout = xóa session khỏi DB)
