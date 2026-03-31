# Prompt 01 — Project Skeleton + Shared Templates

## Role

Bạn là một **Go Project Architect** chuyên thiết kế monorepo. Bạn tạo cấu trúc thư mục chuẩn, rõ ràng, và viết HTML templates dùng chung cho nhiều Go web frameworks.

---

## Context

Khởi tạo monorepo cho dự án so sánh 4 approaches Go web authentication. Tất cả 4 apps dùng chung templates và database. Chi tiết kiến trúc xem `docs/erd.md` và `docs/api-spec.md`.

---

## Yêu cầu

### 1. Cấu trúc thư mục

```
go-auth-frameworks/
├── shared/
│   └── templates/
│       ├── base.html          # Layout chung (head, nav, footer)
│       ├── login.html         # Trang đăng nhập
│       ├── signup.html        # Trang đăng ký
│       ├── home.html          # Trang chủ (thông tin tài khoản)
│       ├── dashboard.html     # Trang quản trị (admin only)
│       └── error.html         # Trang lỗi chung
├── gin/
│   └── cmd/
│       └── main.go            # package main, func main() — skeleton
├── fiber/
│   └── cmd/
│       └── main.go
├── stdlib/
│   └── cmd/
│       └── main.go
├── echo/
│   └── cmd/
│       └── main.go
├── scripts/
│   └── seed.go                # Script seed admin (Prompt 02)
├── docs/
│   ├── roadmap.md
│   ├── erd.md
│   ├── api-spec.md
│   └── prompt01.md ... prompt15.md
├── Makefile
├── .gitignore
└── README.md
```

### 2. Shared Templates

Mỗi page template phải sử dụng **template inheritance** qua Go `html/template`:

#### base.html

```html
{{define "base"}}
<!DOCTYPE html>
<html lang="vi">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go Auth — {{template "title" .}}</title>
    <style>
        /* CSS reset + basic styling */
        /* Dùng embedded CSS — không cần file CSS riêng */
    </style>
</head>
<body>
    <main>
        {{template "content" .}}
    </main>
</body>
</html>
{{end}}
```

#### login.html

```html
{{define "title"}}Đăng nhập{{end}}
{{define "content"}}
<div class="auth-container">
    <h1>Đăng nhập</h1>
    {{if .Error}}<p class="error">{{.Error}}</p>{{end}}
    <form method="POST" action="/login">
        <input type="email" name="email" placeholder="Email" required>
        <input type="password" name="password" placeholder="Mật khẩu" required>
        <button type="submit">Đăng nhập</button>
    </form>
    <p>Chưa có tài khoản? <a href="/signup">Đăng ký</a></p>

    <!-- Divider + Google OAuth button (Phase 7) -->
    <div class="divider"><span>Hoặc</span></div>
    <a href="/auth/google" class="google-btn">Đăng nhập với Google</a>
</div>
{{end}}
```

#### signup.html

```html
{{define "title"}}Đăng ký{{end}}
{{define "content"}}
<div class="auth-container">
    <h1>Đăng ký</h1>
    {{if .Error}}<p class="error">{{.Error}}</p>{{end}}
    <form method="POST" action="/signup">
        <input type="text" name="name" placeholder="Họ tên" required>
        <input type="email" name="email" placeholder="Email" required>
        <input type="tel" name="phone" placeholder="Số điện thoại" required>
        <input type="password" name="password" placeholder="Mật khẩu" required>
        <input type="password" name="confirm_password" placeholder="Xác nhận mật khẩu" required>
        <button type="submit">Đăng ký</button>
    </form>
    <p>Đã có tài khoản? <a href="/login">Đăng nhập</a></p>
</div>
{{end}}
```

#### home.html

```html
{{define "title"}}Trang chủ{{end}}
{{define "content"}}
<div class="home-container">
    <h1>Thông tin tài khoản</h1>
    <div class="user-info">
        <p><strong>Tên:</strong> {{.User.Name}}</p>
        <p><strong>Email:</strong> {{.User.Email}}</p>
        <p><strong>Điện thoại:</strong> {{.User.Phone}}</p>
        <p><strong>Vai trò:</strong> {{.User.Role}}</p>
    </div>
    <p class="framework-badge">Framework: {{.Framework}}</p>
    <form method="POST" action="/logout">
        <button type="submit" class="logout-btn">Đăng xuất</button>
    </form>
</div>
{{end}}
```

#### dashboard.html

```html
{{define "title"}}Dashboard{{end}}
{{define "content"}}
<div class="dashboard-container">
    <h1>Chào mừng quản trị viên, {{.User.Name}}!</h1>
    <p class="framework-badge">Framework: {{.Framework}}</p>
    <form method="POST" action="/logout">
        <button type="submit" class="logout-btn">Đăng xuất</button>
    </form>
</div>
{{end}}
```

#### error.html

```html
{{define "title"}}Lỗi{{end}}
{{define "content"}}
<div class="error-container">
    <h1>{{.Code}}</h1>
    <p>{{.Message}}</p>
    <a href="/login">Quay lại đăng nhập</a>
</div>
{{end}}
```

### 3. Cơ chế Template Inheritance

Khi render, mỗi page phải được parse **riêng** với `base.html`:

```go
// ✅ ĐÚNG — parse từng cặp
loginTmpl := template.Must(template.ParseFiles("base.html", "login.html"))
signupTmpl := template.Must(template.ParseFiles("base.html", "signup.html"))

// ❌ SAI — parse tất cả cùng lúc (blocks bị override)
templates := template.Must(template.ParseGlob("*.html"))
```

### 4. Skeleton main.go (cho mỗi approach)

```go
package main

import "fmt"

func main() {
    fmt.Println("TODO: implement")
}
```

### 5. .gitignore

```
.env
*.exe
tmp/
vendor/
```

---

## Anti-Patterns (KHÔNG được làm)

❌ Không dùng `template.ParseGlob("*.html")` — sẽ bị block override
❌ Không đặt tên template file trùng với tên define block
❌ Không tạo CSS file riêng — dùng embedded `<style>` trong base.html
❌ Không dùng `register.html` — phải là `signup.html`

---

## Acceptance Criteria

1. 6 template files tồn tại trong `shared/templates/`
2. 4 skeleton main.go tồn tại (gin, fiber, stdlib, echo)
3. Template files parse được không lỗi: `template.ParseFiles("base.html", "login.html")`
4. `signup.html` form có 5 fields: name, email, phone, password, confirm_password
5. `home.html` hiển thị: Name, Email, Phone, Role + nút logout (POST)
6. `dashboard.html` hiển thị: "Chào mừng quản trị viên" + nút logout (POST)
7. Logout dùng `<form method="POST">` (không dùng `<a>` link)

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | Cấu trúc thư mục đúng | 🔲 | |
| 2 | base.html có define "base" + template "title"/"content" | 🔲 | |
| 3 | login.html có form email + password | 🔲 | |
| 4 | login.html có link tới /signup | 🔲 | |
| 5 | login.html có nút Google OAuth (Phase 7) | 🔲 | |
| 6 | signup.html có 5 fields (name, email, phone, password, confirm) | 🔲 | |
| 7 | signup.html có link tới /login | 🔲 | |
| 8 | home.html hiển thị Name, Email, Phone, Role | 🔲 | |
| 9 | home.html có nút logout (POST form) | 🔲 | |
| 10 | home.html hiển thị Framework badge | 🔲 | |
| 11 | dashboard.html có "Chào mừng quản trị viên" | 🔲 | |
| 12 | dashboard.html có nút logout (POST form) | 🔲 | |
| 13 | error.html hiển thị Code + Message | 🔲 | |
| 14 | 4 skeleton main.go tồn tại | 🔲 | |
| 15 | .gitignore có .env | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
