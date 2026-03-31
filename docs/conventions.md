# Conventions — Quyết định chung cho dự án

File này ghi lại các quyết định thiết kế dùng chung cho tất cả 4 approaches. Khi implement từng framework, tham khảo file này thay vì đọc lại toàn bộ `api-spec.md`.

---

## 1. Cookie

| Cookie | MaxAge | HttpOnly | Secure | SameSite | Path |
|--------|--------|----------|--------|----------|------|
| `access_token` | 1800 (30 min) | ✅ | `cfg.IsProduction` | Lax | `/` |
| `refresh_token` | 604800 (7 days) | ✅ | `cfg.IsProduction` | Lax | `/` |
| `oauth_state` | 300 (5 min) | ✅ | `cfg.IsProduction` | Lax | `/` |

> **Secure flag**: `cfg.IsProduction` thay vì hardcode `true`. Local dev dùng HTTP (không có HTTPS), nếu `Secure=true` thì cookie sẽ không gửi → login fail. Set `APP_ENV=production` trên Fly.io để bật.

---

## 2. Config chung

Mỗi approach đều có các fields sau trong Config:

```go
type Config struct {
    MongoURI           string // required
    MongoDB            string // required
    JWTSecret          string // required
    Port               string // default: 808{1,2,3,4}
    Framework          string // hardcode per approach
    TemplateDir        string // default: "../shared/templates"
    IsProduction       bool   // APP_ENV == "production"
    // Phase 7:
    GoogleClientID     string // optional
    GoogleClientSecret string // optional
    GoogleRedirectURL  string // optional
}
```

`TemplateDir`:
- **Local dev**: `../shared/templates` (chạy từ `<approach>/cmd/`)
- **Docker**: `./shared/templates` (binary ở `/app`, templates copy vào `/app/shared/templates/`)
- Đọc từ env `TEMPLATE_DIR`, default `../shared/templates`

`IsProduction`:
- Đọc `APP_ENV` env → nếu `== "production"` → `true`, còn lại `false`

---

## 3. Redirect Status Code

| Tình huống | Status Code | Lý do |
|-----------|-------------|-------|
| POST thành công → redirect | **303 See Other** | Tránh browser resend POST khi refresh |
| GET redirect (ví dụ `/` → `/login`) | **302 Found** | Redirect GET thông thường |
| OAuth redirect → Google | **307 Temporary Redirect** | Giữ nguyên method (GET) |

---

## 4. Error Messages

| Tình huống | Message | Lý do |
|-----------|---------|-------|
| Login: email không tồn tại | `"email hoặc mật khẩu không đúng"` | Chống user enumeration |
| Login: sai password | `"email hoặc mật khẩu không đúng"` | **Cùng message** |
| Signup: email trùng | `"email đã được sử dụng"` | Rõ ràng |
| Signup: phone trùng | `"số điện thoại đã được sử dụng"` | Rõ ràng |
| Signup: password mismatch | `"mật khẩu xác nhận không khớp"` | Rõ ràng |

---

## 5. JWT Claims

```json
{
  "user_id": "hex ObjectID",
  "email": "user@example.com",
  "role": "user",
  "exp": 1711868400,
  "iat": 1711866600
}
```

- Signing method: `HS256`
- Expire: `iat + 30 phút`
- Secret: `cfg.JWTSecret`

---

## 6. Refresh Token

- Generate: `crypto/rand` 32 bytes → `hex.EncodeToString()` (64 chars)
- Lưu: cookie + MongoDB `sessions` collection
- Expire: 7 ngày (TTL index tự xóa)

---

## 7. Template Path Resolution

Tất cả approach dùng `cfg.TemplateDir` thay vì hardcode path:

```go
// Config load
cfg.TemplateDir = getEnv("TEMPLATE_DIR", "../shared/templates")
```

| Môi trường | TEMPLATE_DIR | Lý do |
|-----------|-------------|-------|
| Local dev | `../shared/templates` (default) | Binary ở `<approach>/cmd/`, templates ở `../shared/` |
| Docker | `./shared/templates` | Binary ở `/app/`, templates copy vào `/app/shared/templates/` |

Trong fly.toml:
```toml
[env]
  TEMPLATE_DIR = "./shared/templates"
```

---

## 8. Middleware: HTML vs API Response

Khi user chưa auth:
- **HTML routes** (`/home`, `/dashboard`, `/logout`) → redirect `/login`
- **API routes** (`/api/me`) → return `401 JSON {"error": "unauthorized"}`

Cách phân biệt: `strings.HasPrefix(path, "/api/")`

---

## 9. Role-based Redirect

Sau login thành công:
- `role == "admin"` → redirect `/dashboard`
- `role == "user"` → redirect `/home`

---

## 10. Bcrypt

- Cost: **12** (production standard)
- Package: `golang.org/x/crypto/bcrypt`

---

## 11. fly.toml chung

```toml
[env]
  PORT = "8080"
  MONGO_DB = "goauth"
  APP_ENV = "production"
  TEMPLATE_DIR = "./shared/templates"
```

**Secrets** (set qua CLI, KHÔNG đặt trong file):
```bash
fly secrets set MONGO_URI="..." JWT_SECRET="..." -a goauth-<approach>
```
