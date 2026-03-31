# Prompt 04 — Gin: Service + Handler + Middleware

## Role

Bạn là một **Go Backend Engineer** chuyên về authentication với JWT, session management, và Gin framework. Bạn implement hệ thống Access Token + Refresh Token + Auto-refresh middleware.

---

## Context

Tiếp theo Prompt 03 (Config, DB, Model, Repository). Bây giờ implement tầng business logic + transport layer:

- **Service**: signup, login (tạo 2 tokens + session), logout (xóa session), refresh, getme
- **Handler**: nhận request Gin, gọi service, render template hoặc trả JSON
- **Middleware**: validate access token, auto-refresh khi expired

Auth flow chi tiết xem `docs/api-spec.md` (section "Auth Middleware Flow").

---

## Dependencies (Prompt phụ thuộc)

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `shared/templates/*.html` (login, signup, home, dashboard, error) |
| 03 | Config, DB, User Model, Session Model, UserRepository, SessionRepository |

---

## Yêu cầu

### 1. gin/internal/service/auth.go

```go
type AuthService interface {
    Signup(ctx context.Context, name, email, phone, password string) error
    Login(ctx context.Context, email, password string) (accessToken, refreshToken, role string, err error)
    Logout(ctx context.Context, refreshToken string) error
    RefreshAccessToken(ctx context.Context, refreshToken string) (newAccessToken string, err error)
    GetUserByID(ctx context.Context, id string) (*model.User, error)
}
```

Implement `authService`:

- **Signup**:
  - Kiểm tra email trùng → error `"email đã được sử dụng"`
  - Kiểm tra phone trùng → error `"số điện thoại đã được sử dụng"`
  - Hash password bcrypt cost=12
  - Tạo user: Role=`"user"`, Provider=`"local"`, GoogleID=`""`
  - Set `CreatedAt`, `UpdatedAt` = `time.Now()`

- **Login**:
  - Tìm user theo email → error `"email hoặc mật khẩu không đúng"` (chống user enumeration)
  - So sánh bcrypt → **cùng error message** nếu sai
  - Tạo access token JWT (30 phút): claims `user_id`, `email`, `role`, `exp`, `iat`
  - Tạo refresh token: `crypto/rand` 32 bytes → hex encode
  - Tạo session trong DB: `UserID`, `RefreshToken`, `ExpiredAt` = now + 7 ngày
  - Return: accessToken, refreshToken, user.Role

- **Logout**:
  - Xóa session theo refresh token: `sessionRepo.DeleteByRefreshToken()`

- **RefreshAccessToken**:
  - Tìm session theo refresh token
  - Nếu không thấy hoặc expired → error
  - Lấy user theo session.UserID
  - Nếu user không tồn tại → xóa session → error
  - Tạo access token mới → return

- **GetUserByID**: gọi thẳng `userRepo.FindByID()`

Constructor: `NewAuthService(userRepo, sessionRepo, jwtSecret) AuthService`

### 2. gin/internal/handler/auth.go

**ShowLogin** — `GET /login`:
- Nếu đã có valid access token → redirect `/home`
- Render `login.html` với `gin.H{"Error": ""}`

**ShowSignup** — `GET /signup`:
- Nếu đã có valid access token → redirect `/home`
- Render `signup.html` với `gin.H{"Error": ""}`

**Login** — `POST /login`:
- `c.PostForm("email")`, `c.PostForm("password")`
- Gọi `service.Login()`
- Nếu lỗi → render lại `login.html` với error
- Thành công → set 2 cookies:
  ```go
  // Access token cookie — 30 phút
  c.SetCookie("access_token", accessToken, 1800, "/", "", true, true)
  // Params: name, value, maxAge, path, domain, secure, httpOnly

  // Refresh token cookie — 7 ngày
  c.SetCookie("refresh_token", refreshToken, 604800, "/", "", true, true)
  ```
- Redirect theo role:
  ```go
  if role == "admin" {
      c.Redirect(http.StatusSeeOther, "/dashboard")
  } else {
      c.Redirect(http.StatusSeeOther, "/home")
  }
  ```

**Signup** — `POST /signup`:
- `c.PostForm("name")`, `c.PostForm("email")`, `c.PostForm("phone")`, `c.PostForm("password")`, `c.PostForm("confirm_password")`
- Validate: confirm_password khớp → render error nếu khác
- Gọi `service.Signup()`
- Thành công → redirect `/login` (303)

**Logout** — `POST /logout`:
- Đọc refresh token cookie: `refreshToken, _ := c.Cookie("refresh_token")`
- Gọi `service.Logout(refreshToken)`
- Xóa cả 2 cookies:
  ```go
  c.SetCookie("access_token", "", -1, "/", "", true, true)
  c.SetCookie("refresh_token", "", -1, "/", "", true, true)
  ```
- Redirect `/login` (303)

### 3. gin/internal/handler/home.go

**ShowHome** — `GET /home` (auth required, any role):
- `user := c.MustGet("user").(*model.User)`
- Render `home.html`:
  ```go
  c.HTML(200, "home", gin.H{
      "User":      user,
      "Framework": cfg.Framework,
  })
  ```

### 4. gin/internal/handler/dashboard.go

**ShowDashboard** — `GET /dashboard` (auth required, admin only):
- `user := c.MustGet("user").(*model.User)`
- Nếu `user.Role != "admin"` → redirect `/home`
- Render `dashboard.html`:
  ```go
  c.HTML(200, "dashboard", gin.H{
      "User":      user,
      "Framework": cfg.Framework,
  })
  ```

### 5. gin/internal/handler/api.go

**GetMe** — `GET /api/me` (auth required, JSON):
- `user := c.MustGet("user").(*model.User)`
- Trả JSON (password tự ẩn nhờ `json:"-"`):
  ```go
  c.JSON(200, gin.H{"user": user})
  ```

### 6. gin/internal/middleware/auth.go

**RequireAuth** — Middleware auto-refresh:

```go
func RequireAuth(service service.AuthService, jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Đọc access_token cookie
        accessToken, err := c.Cookie("access_token")

        // 2. Nếu có → parse JWT
        if err == nil && accessToken != "" {
            claims, err := parseJWT(accessToken, jwtSecret)
            if err == nil {
                // Token valid → lấy user → set context → continue
                user, _ := service.GetUserByID(c.Request.Context(), claims.UserID)
                if user != nil {
                    c.Set("user", user)
                    c.Next()
                    return
                }
            }
        }

        // 3. Access token missing/invalid/expired → thử refresh
        refreshToken, err := c.Cookie("refresh_token")
        if err != nil || refreshToken == "" {
            redirectOrUnauthorized(c) // redirect /login (HTML) hoặc 401 (API)
            return
        }

        // 4. Auto-refresh: tạo access token mới
        newAccessToken, err := service.RefreshAccessToken(c.Request.Context(), refreshToken)
        if err != nil {
            // Refresh token cũng invalid → xóa cookies → redirect
            c.SetCookie("access_token", "", -1, "/", "", true, true)
            c.SetCookie("refresh_token", "", -1, "/", "", true, true)
            redirectOrUnauthorized(c)
            return
        }

        // 5. Set cookie mới + lấy user
        c.SetCookie("access_token", newAccessToken, 1800, "/", "", true, true)
        claims, _ := parseJWT(newAccessToken, jwtSecret)
        user, _ := service.GetUserByID(c.Request.Context(), claims.UserID)
        if user == nil {
            redirectOrUnauthorized(c)
            return
        }
        c.Set("user", user)
        c.Next()
    }
}

// redirectOrUnauthorized: HTML routes → redirect /login, API routes → 401 JSON
func redirectOrUnauthorized(c *gin.Context) {
    if strings.HasPrefix(c.Request.URL.Path, "/api/") {
        c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
    } else {
        c.Redirect(http.StatusFound, "/login")
        c.Abort()
    }
}
```

---

## Hướng dẫn thực hiện

1. Login set **2 cookies** — access_token (30 phút) + refresh_token (7 ngày)
2. Cookies: `Secure=true`, `HttpOnly=true` (Gin SetCookie tham số = `true, true`)
3. POST redirect dùng **303 (See Other)**
4. Login redirect dựa trên **role**: admin→/dashboard, user→/home
5. Logout dùng **POST** method (không phải GET)
6. Error message login: luôn `"email hoặc mật khẩu không đúng"` (chống user enumeration)
7. Middleware phân biệt HTML vs API routes khi unauthorized

---

## Anti-Patterns (KHÔNG được làm)

❌ Không chỉ set 1 cookie — phải có cả access_token VÀ refresh_token
❌ Không return JWT/refresh token trong response body — chỉ set cookie
❌ Không tiết lộ "email không tồn tại" vs "sai password"
❌ Không dùng GET cho logout — dùng POST
❌ Không dùng 302 cho POST redirect — dùng 303

---

## Acceptance Criteria

1. `cd gin && go build ./...` pass
2. `cd gin && go vet ./...` không warning
3. AuthService interface có đúng 5 methods
4. Login response set 2 cookies
5. Admin login → redirect /dashboard; User login → redirect /home
6. Middleware auto-refresh khi access token expired

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | AuthService có 5 methods (Signup, Login, Logout, RefreshAccessToken, GetUserByID) | 🔲 | |
| 2 | Signup kiểm tra duplicate email VÀ phone | 🔲 | |
| 3 | Signup dùng bcrypt cost=12 | 🔲 | |
| 4 | Login return accessToken + refreshToken + role | 🔲 | |
| 5 | Login tạo session trong DB | 🔲 | |
| 6 | Login error message chống user enumeration | 🔲 | |
| 7 | Refresh token = crypto/rand 32 bytes hex | 🔲 | |
| 8 | Access token JWT: user_id, email, role, exp (30 phút) | 🔲 | |
| 9 | ShowLogin redirect /home nếu đã authenticated | 🔲 | |
| 10 | POST Login set 2 cookies (access + refresh) | 🔲 | |
| 11 | POST Login redirect theo role (303) | 🔲 | |
| 12 | Signup form có 5 fields (name, email, phone, password, confirm) | 🔲 | |
| 13 | Logout dùng POST method | 🔲 | |
| 14 | Logout xóa session DB + xóa 2 cookies | 🔲 | |
| 15 | ShowHome lấy user từ context, render home.html | 🔲 | |
| 16 | ShowDashboard kiểm tra role=admin | 🔲 | |
| 17 | GetMe trả JSON (password ẩn nhờ json:"-") | 🔲 | |
| 18 | Middleware auto-refresh access token | 🔲 | |
| 19 | Middleware phân biệt HTML redirect vs API 401 | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
