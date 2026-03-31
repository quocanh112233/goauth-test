# API Specification

## Base URLs

| Approach | Local | Production |
|----------|-------|-----------|
| Gin | `http://localhost:8081` | `https://goauth-gin.fly.dev` |
| Fiber | `http://localhost:8082` | `https://goauth-fiber.fly.dev` |
| net/http | `http://localhost:8083` | `https://goauth-stdlib.fly.dev` |
| Echo | `http://localhost:8084` | `https://goauth-echo.fly.dev` |

---

## Pages (Server-Side Rendering)

### GET `/`

Redirect tới `/login`.

| | |
|---|---|
| **Response** | 302 → `/login` |

---

### GET `/login`

Trang đăng nhập. Nếu đã có access token hợp lệ → redirect `/home`.

| | |
|---|---|
| **Response** | 200 — HTML `login.html` |
| **Template data** | `{ Error: string }` |

---

### POST `/login`

Xử lý đăng nhập.

| | |
|---|---|
| **Content-Type** | `application/x-www-form-urlencoded` |
| **Body** | |

| Field | Type | Required | Validation |
|-------|------|----------|-----------|
| `email` | string | ✅ | Email format |
| `password` | string | ✅ | Không rỗng |

**Thành công:**

| | |
|---|---|
| **Set-Cookie** | `access_token` — JWT, HttpOnly, Secure, SameSite=Lax, MaxAge=1800 (30 phút) |
| **Set-Cookie** | `refresh_token` — random hex, HttpOnly, Secure, SameSite=Lax, MaxAge=604800 (7 ngày) |
| **Side effect** | Tạo session trong MongoDB `sessions` collection |
| **Response** | 303 → `/home` (nếu role=user) hoặc `/dashboard` (nếu role=admin) |

**Thất bại:**

| | |
|---|---|
| **Response** | 200 — render lại `login.html` với `{ Error: "email hoặc mật khẩu không đúng" }` |

> **Bảo mật**: Error message giống nhau cho cả "email không tồn tại" và "sai password" (chống user enumeration).

---

### GET `/signup`

Trang đăng ký. Nếu đã có access token hợp lệ → redirect `/home`.

| | |
|---|---|
| **Response** | 200 — HTML `signup.html` |
| **Template data** | `{ Error: string }` |

---

### POST `/signup`

Xử lý đăng ký tài khoản mới.

| | |
|---|---|
| **Content-Type** | `application/x-www-form-urlencoded` |
| **Body** | |

| Field | Type | Required | Validation |
|-------|------|----------|-----------|
| `name` | string | ✅ | Không rỗng |
| `email` | string | ✅ | Email format, unique |
| `phone` | string | ✅ | Không rỗng, unique |
| `password` | string | ✅ | Tối thiểu 6 ký tự |
| `confirm_password` | string | ✅ | Phải khớp `password` |

**Thành công:**

| | |
|---|---|
| **Side effect** | Tạo user trong MongoDB: role="user", provider="local", password=bcrypt(cost=12) |
| **Response** | 303 → `/login` |

**Thất bại:**

| | |
|---|---|
| **Response** | 200 — render lại `signup.html` với `{ Error: "..." }` |

Possible errors:
- `"mật khẩu xác nhận không khớp"`
- `"email đã được sử dụng"`
- `"số điện thoại đã được sử dụng"`

---

### GET `/home`

Trang chủ — hiển thị thông tin tài khoản + nút đăng xuất.

| | |
|---|---|
| **Auth** | ✅ Required (access token cookie) |
| **Role** | Bất kỳ (user hoặc admin) |
| **Response** | 200 — HTML `home.html` |
| **Template data** | `{ User: { Name, Email, Phone, Role }, Framework: string }` |

> **Nội dung hiển thị**: Tên, email, số điện thoại, role. **KHÔNG hiển thị password**.

---

### GET `/dashboard`

Trang quản trị — chỉ dành cho admin.

| | |
|---|---|
| **Auth** | ✅ Required (access token cookie) |
| **Role** | `admin` only |
| **Response 200** | HTML `dashboard.html` — `{ User: { Name }, Framework: string }` |
| **Response 403** | Nếu role ≠ admin → redirect `/home` |

> **Nội dung hiển thị**: "Chào mừng quản trị viên, [Name]!" + nút đăng xuất.

---

### POST `/logout`

Đăng xuất — xóa session + cookies.

| | |
|---|---|
| **Auth** | ✅ Required |
| **Side effect** | Xóa session khỏi MongoDB (bằng refresh_token) |
| **Clear-Cookie** | `access_token` (MaxAge=-1) |
| **Clear-Cookie** | `refresh_token` (MaxAge=-1) |
| **Response** | 303 → `/login` |

> **Tại sao xóa session DB?** Để revoke refresh token — ngay cả khi attacker có cookie, server sẽ không cấp access token mới.

---

## API Endpoints (JSON)

### GET `/api/me`

Lấy thông tin user hiện tại.

| | |
|---|---|
| **Auth** | ✅ Required (access token cookie) |
| **Content-Type** | `application/json` |

**Response 200:**

```json
{
  "user": {
    "id": "665a1b2c3d4e5f6a7b8c9d0e",
    "name": "Nguyễn Văn A",
    "email": "user@example.com",
    "phone": "0901234567",
    "role": "user",
    "created_at": "2026-01-15T10:30:00Z",
    "updated_at": "2026-03-31T08:00:00Z"
  }
}
```

> **KHÔNG bao gồm**: `password`, `provider`, `google_id`.

**Response 401:**

```json
{
  "error": "unauthorized"
}
```

---

### GET `/health`

Health check cho Fly.io monitoring.

| | |
|---|---|
| **Auth** | ❌ Không cần |
| **Content-Type** | `application/json` |

**Response 200:**

```json
{
  "status": "ok",
  "framework": "Gin",
  "db": "connected"
}
```

**Response 503:**

```json
{
  "status": "error",
  "framework": "Gin",
  "db": "disconnected"
}
```

---

## Google OAuth Endpoints (Phase 7)

> Chỉ implement sau khi toàn bộ đăng ký/đăng nhập thủ công hoạt động ổn định.

### GET `/auth/google`

Khởi tạo Google OAuth flow.

| | |
|---|---|
| **Auth** | ❌ Không cần |
| **Set-Cookie** | `oauth_state` — random hex, HttpOnly, Secure, MaxAge=300 (5 phút) |
| **Response** | 307 → Google OAuth URL |

### GET `/auth/google/callback`

Xử lý callback từ Google.

| | |
|---|---|
| **Auth** | ❌ Không cần |
| **Query params** | `state`, `code` |
| **Validation** | `state` query phải khớp `oauth_state` cookie (chống CSRF) |
| **Side effect** | Upsert user (tìm google_id → email → tạo mới) + tạo session |
| **Set-Cookie** | `access_token` + `refresh_token` (giống POST /login) |
| **Clear-Cookie** | `oauth_state` |
| **Response** | 303 → `/home` (hoặc `/dashboard` nếu admin) |

---

## Auth Middleware Flow

```
Request vào
    │
    ▼
Đọc cookie "access_token"
    │
    ├── Có token + JWT valid + chưa expired
    │       → Decode user_id, email, role
    │       → Set user vào context
    │       → ✅ Tiếp tục
    │
    ├── Có token + JWT expired
    │       → Đọc cookie "refresh_token"
    │       → Tìm trong MongoDB sessions
    │       │
    │       ├── Tìm thấy + chưa expired
    │       │       → Cấp access_token mới (set cookie)
    │       │       → Lấy user từ DB
    │       │       → Set user vào context
    │       │       → ✅ Tiếp tục
    │       │
    │       └── Không tìm thấy / expired
    │               → Xóa cả 2 cookies
    │               → 🔒 Redirect /login (HTML) hoặc 401 (API)
    │
    └── Không có token
            → Đọc cookie "refresh_token"
            │
            ├── Có refresh_token + valid trong DB
            │       → Cấp access_token mới
            │       → ✅ Tiếp tục
            │
            └── Không có
                    → 🔒 Redirect /login (HTML) hoặc 401 (API)
```

> **HTML routes** (/, /home, /dashboard) → redirect `/login` khi chưa auth
> **API routes** (/api/me) → return `401 JSON` khi chưa auth

---

## Cookie Specification

| Cookie | Value | HttpOnly | Secure | SameSite | MaxAge | Path |
|--------|-------|----------|--------|----------|--------|------|
| `access_token` | JWT string | ✅ | ✅ | Lax | 1800 (30 min) | `/` |
| `refresh_token` | Random hex 32 bytes | ✅ | ✅ | Lax | 604800 (7 days) | `/` |
| `oauth_state` | Random hex 16 bytes | ✅ | ✅ | Lax | 300 (5 min) | `/` |

---

## JWT Claims (Access Token)

```json
{
  "user_id": "665a1b2c3d4e5f6a7b8c9d0e",
  "email": "user@example.com",
  "role": "user",
  "exp": 1711868400,
  "iat": 1711866600
}
```

| Claim | Type | Mô tả |
|-------|------|-------|
| `user_id` | string | MongoDB ObjectID dạng hex |
| `email` | string | Email của user |
| `role` | string | `"admin"` hoặc `"user"` |
| `exp` | int64 | Unix timestamp hết hạn (iat + 30 phút) |
| `iat` | int64 | Unix timestamp thời điểm tạo |

---

## Route Summary

| Method | Path | Auth | Role | Type | Mô tả |
|--------|------|------|------|------|-------|
| GET | `/` | ❌ | — | Redirect | → `/login` |
| GET | `/login` | ❌ | — | HTML | Trang đăng nhập |
| POST | `/login` | ❌ | — | Form | Xử lý đăng nhập |
| GET | `/signup` | ❌ | — | HTML | Trang đăng ký |
| POST | `/signup` | ❌ | — | Form | Xử lý đăng ký |
| GET | `/home` | ✅ | any | HTML | Trang chủ (thông tin user) |
| GET | `/dashboard` | ✅ | admin | HTML | Trang quản trị |
| POST | `/logout` | ✅ | any | Form | Đăng xuất |
| GET | `/api/me` | ✅ | any | JSON | Lấy thông tin user |
| GET | `/health` | ❌ | — | JSON | Health check |
| GET | `/auth/google` | ❌ | — | Redirect | OAuth initiate (Phase 7) |
| GET | `/auth/google/callback` | ❌ | — | Redirect | OAuth callback (Phase 7) |
