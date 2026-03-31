# Báo cáo đánh giá — Prompt 01 & Prompt 02

> Ngày: 2026-03-31
> Trạng thái: Cần sửa trước khi tiếp tục Prompt 03

---

## Prompt 01 — Project Skeleton + Shared Templates

### Checklist kết quả

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | 6 template files tồn tại | ✅ | base, login, signup, home, dashboard, error |
| 2 | Mỗi template có comment vai trò | ❌ | Không có comment HTML giải thích vai trò/role access |
| 3 | login.html form email + password | ✅ | |
| 4 | login.html link /signup + nút Google | ✅ | Có Google SVG icon thật |
| 5 | signup.html 5 fields | ✅ | name, email, phone, password, confirm_password |
| 6 | home.html: Name, Email, Phone, Role + logout POST | ✅ | |
| 7 | dashboard.html: "Chào mừng quản trị viên" + logout POST | ⚠️ | Text nằm trong `<p>` thay vì `<h1>` |
| 8 | error.html: Code + Message | ✅ | |
| 9 | 4 skeleton main.go | ✅ | gin, fiber, stdlib, echo |
| 10 | .gitignore có .env | ✅ | Có thêm .DS_Store, *.log |

### Tổng kết

| Trạng thái | Số lượng |
|------------|---------|
| ✅ Tuân thủ | 8 |
| ⚠️ Thiếu sót | 1 |
| ❌ Vi phạm | 1 |
| 🔲 Chưa implement | 0 |
| ➕ Ngoài yêu cầu | 3 |

### Ngoài yêu cầu ➕

- CSS embedded glassmorphism design đẹp, vượt xa yêu cầu tối thiểu
- Google SVG icon thật thay vì text đơn giản
- .gitignore mở rộng thêm `.DS_Store`, `*.log`

---

## Prompt 02 — MongoDB Seed Script

### Checklist kết quả

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | go.mod đúng module + dependencies | ⚠️ | `go 1.25.5` không tồn tại — Go mới nhất là 1.22 |
| 2 | Collection users + 2 unique indexes (email, phone) | ✅ | |
| 3 | Collection sessions + TTL index (expired_at) | ✅ | |
| 4 | Sessions indexes (user_id, refresh_token) | ✅ | |
| 5 | Admin user đủ fields | ✅ | |
| 6 | Admin phone = "0900000000" | ❌ | Dùng `"0123456789"` thay vì `"0900000000"` |
| 7 | Admin password bcrypt cost=12 | ✅ | |
| 8 | Script idempotent | ✅ | |
| 9 | .env.example có MONGO_URI + MONGO_DB | ✅ | |

### Tổng kết

| Trạng thái | Số lượng |
|------------|---------|
| ✅ Tuân thủ | 7 |
| ⚠️ Thiếu sót | 1 |
| ❌ Vi phạm | 1 |
| 🔲 Chưa implement | 0 |
| ➕ Ngoài yêu cầu | 1 |

### Ngoài yêu cầu ➕

- Thêm `google_id` sparse index ngay từ Prompt 02 (Phase 7 ready ahead of time)

---

## Danh sách việc cần sửa trước Prompt 03

### 🔴 Ưu tiên cao — BẮT BUỘC sửa

#### 1. Lỗi template inheritance trong tất cả page templates

**Vấn đề:** Tất cả page templates (login, signup, home, dashboard, error) đang dùng `{{template "base" .}}` ở đầu file. Đây là vi phạm cơ chế `html/template` của Go — sẽ gây lỗi hoặc render sai khi các frameworks gọi `ExecuteTemplate(w, "base", data)`.

**Giải thích cơ chế đúng:**
Khi parse `template.ParseFiles("base.html", "login.html")` rồi gọi `ExecuteTemplate(w, "base", data)`:

- Go execute template có tên `"base"` (defined trong `base.html`)
- Bên trong `base.html` có `{{template "title" .}}` và `{{template "content" .}}`
- Go tự động tìm các `{{define "title"}}` và `{{define "content"}}` trong cùng template set
- Hoạt động đúng mà **không cần** gọi `{{template "base" .}}` trong page templates

**File cần sửa:** `shared/templates/login.html`, `signup.html`, `home.html`, `dashboard.html`, `error.html`

**Cách sửa — xóa dòng đầu tiên của mỗi file:**

```html
<!-- XÓA dòng này khỏi tất cả page templates -->
{{template "base" .}}

<!-- CHỈ GIỮ các define blocks -->
{{define "title"}}Đăng nhập{{end}}
{{define "content"}}
...
{{end}}
```

---

#### 2. Admin phone sai trong seed.go

**File:** `scripts/seed.go`

**Hiện tại:**

```go
"phone": "0123456789",
```

**Phải sửa thành:**

```go
"phone": "0900000000",
```

**Lý do:** Spec trong `docs/prompt02.md` yêu cầu cụ thể `"0900000000"`. Nếu để sai, các unit test ở Prompt 16 có thể fail khi test dựa trên seed data.

---

### 🟡 Ưu tiên thấp — Nên sửa

#### 3. Go version sai trong scripts/go.mod

**File:** `scripts/go.mod`

**Hiện tại:**

```
go 1.25.5
```

**Sửa thành:**

```
go 1.22
```

**Lý do:** Go 1.25.5 không tồn tại. Phiên bản yêu cầu là 1.22 (theo `docs/prompt02.md`). Có thể gây lỗi khi người khác clone repo và chạy `go mod tidy`.

---

#### 4. Thiếu comment HTML trong page templates

**File:** Tất cả page templates trong `shared/templates/`

**Yêu cầu gốc (prompt01.md):** Mỗi template phải có comment HTML giải thích vai trò + role access.

**Ví dụ cần thêm vào `home.html`:**

```html
<!-- Trang chủ — hiển thị thông tin tài khoản -->
<!-- Truy cập được bởi TẤT CẢ user đã đăng nhập (role=user hoặc admin) -->
<!-- User role=user được redirect tới đây sau login -->
```

**Lý do:** Không ảnh hưởng runtime, nhưng quan trọng cho blog — người đọc cần hiểu role access policy từ template.

---

## Checklist sửa trước khi bắt đầu Prompt 03

| # | Việc cần làm | File | Ưu tiên | Xong chưa |
|---|-------------|------|---------|-----------|
| 1 | Xóa `{{template "base" .}}` khỏi login.html | shared/templates/login.html | 🔴 Cao | 🔲 |
| 2 | Xóa `{{template "base" .}}` khỏi signup.html | shared/templates/signup.html | 🔴 Cao | 🔲 |
| 3 | Xóa `{{template "base" .}}` khỏi home.html | shared/templates/home.html | 🔴 Cao | 🔲 |
| 4 | Xóa `{{template "base" .}}` khỏi dashboard.html | shared/templates/dashboard.html | 🔴 Cao | 🔲 |
| 5 | Xóa `{{template "base" .}}` khỏi error.html | shared/templates/error.html | 🔴 Cao | 🔲 |
| 6 | Đổi phone admin thành "0900000000" | scripts/seed.go | 🔴 Cao | 🔲 |
| 7 | Đổi `go 1.25.5` thành `go 1.22` | scripts/go.mod | 🟡 Thấp | 🔲 |
| 8 | Thêm comment HTML vai trò vào 5 page templates | shared/templates/*.html | 🟡 Thấp | 🔲 |
