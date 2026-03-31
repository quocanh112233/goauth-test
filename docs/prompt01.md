# Prompt 01 — Project Skeleton + Shared Templates

## Role

Bạn là một **Go Project Architect** chuyên thiết kế monorepo. Bạn tạo cấu trúc thư mục chuẩn và viết HTML templates dùng chung cho nhiều Go web frameworks.

---

## Context

Khởi tạo monorepo cho dự án so sánh 4 approaches Go web authentication. Chi tiết kiến trúc xem `docs/erd.md`, `docs/api-spec.md`, `docs/conventions.md`.

---

## Yêu cầu

### 1. Cấu trúc thư mục

```
go-auth-frameworks/
├── shared/
│   └── templates/
│       ├── base.html          # Layout chung
│       ├── login.html         # Trang đăng nhập (mặc định khi mở web)
│       ├── signup.html        # Trang đăng ký
│       ├── home.html          # Trang chủ — hiển thị thông tin tài khoản (mọi role)
│       ├── dashboard.html     # Trang quản trị — chỉ dành cho admin
│       └── error.html         # Trang lỗi chung
├── gin/
│   └── cmd/main.go
├── fiber/
│   └── cmd/main.go
├── stdlib/
│   └── cmd/main.go
├── echo/
│   └── cmd/main.go
├── scripts/
│   └── seed.go
├── docs/
│   ├── roadmap.md, erd.md, api-spec.md, conventions.md
│   └── prompt01.md ... prompt16.md
├── Makefile
├── .gitignore
└── README.md
```

### 2. Shared Templates

#### base.html

```html
{{define "base"}}
<!DOCTYPE html>
<html lang="vi">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go Auth — {{template "title" .}}</title>
    <style>/* CSS embedded */</style>
</head>
<body>
    <main>{{template "content" .}}</main>
</body>
</html>
{{end}}
```

#### login.html

```html
<!-- Trang đăng nhập — hiển thị mặc định khi mở web -->
<!-- Nếu đã authenticated → redirect /home (xử lý trong handler, không phải template) -->
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
    <div class="divider"><span>Hoặc</span></div>
    <a href="/auth/google" class="google-btn">Đăng nhập với Google</a>
</div>
{{end}}
```

#### signup.html

```html
<!-- Trang đăng ký — tạo tài khoản mới -->
<!-- Khi đăng ký thành công → redirect /login (handler xử lý) -->
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
<!-- Trang chủ — hiển thị thông tin tài khoản -->
<!-- Truy cập được bởi TẤT CẢ user đã đăng nhập (role=user hoặc admin) -->
<!-- User role=user được redirect tới đây sau login -->
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
<!-- Trang quản trị — CHỈ dành cho admin (role=admin) -->
<!-- User role=user truy cập → redirect /home (handler xử lý) -->
<!-- User role=admin được redirect tới đây sau login -->
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
<!-- Trang lỗi chung: 403, 404, 500... -->
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

```go
// ✅ ĐÚNG — parse từng cặp riêng
loginTmpl := template.Must(template.ParseFiles("base.html", "login.html"))

// ❌ SAI — blocks bị override
templates := template.Must(template.ParseGlob("*.html"))
```

### 4. Skeleton main.go

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

## Anti-Patterns

❌ Không dùng `template.ParseGlob("*.html")`
❌ Không dùng `register.html` — phải là `signup.html`
❌ Không dùng `<a href="/logout">` — logout phải dùng `<form method="POST">`

---

## Acceptance Criteria

1. 6 template files tồn tại trong `shared/templates/`
2. Mỗi template có comment HTML giải thích vai trò + role access
3. `signup.html` form có 5 fields: name, email, phone, password, confirm_password
4. `home.html` hiển thị: Name, Email, Phone, Role + nút logout POST
5. `dashboard.html` + `home.html` đều dùng `<form method="POST" action="/logout">` cho logout

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | 6 template files tồn tại | 🔲 | |
| 2 | Mỗi template có comment vai trò | 🔲 | |
| 3 | login.html form email + password | 🔲 | |
| 4 | login.html link /signup + nút Google | 🔲 | |
| 5 | signup.html 5 fields (name, email, phone, password, confirm) | 🔲 | |
| 6 | home.html: Name, Email, Phone, Role + logout POST | 🔲 | |
| 7 | dashboard.html: "Chào mừng quản trị viên" + logout POST | 🔲 | |
| 8 | error.html: Code + Message | 🔲 | |
| 9 | 4 skeleton main.go | 🔲 | |
| 10 | .gitignore có .env | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
