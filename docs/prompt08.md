# Prompt 08 — net/http (stdlib): Config + DB + Model + Repository

## Role

Bạn là một **Go Backend Engineer** chuyên Go standard library. Prompt này chỉ xây dựng **tầng data** — giống Gin Prompt 03.

---

## Context

Phase 4 bắt đầu. Approach **zero-framework**. Go 1.22+ required cho enhanced ServeMux (Prompt 09–10). Tầng data (config, db, model, repository) giống hệt Gin — không phụ thuộc framework.

> stdlib được **tách 3 prompt** (08, 09, 10) vì handler/middleware/template phải tự viết toàn bộ, phức tạp hơn các framework có sẵn.

---

## Dependencies

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `stdlib/` folder skeleton |
| 02 | MongoDB Atlas đã seed |

---

## Yêu cầu

### 1. stdlib/go.mod

- Module: `github.com/yourusername/go-auth-frameworks/stdlib`
- Go: 1.22
- Dependencies (**KHÔNG có web framework**):
  - `go.mongodb.org/mongo-driver v1.15.0`
  - `golang.org/x/crypto v0.22.0`
  - `github.com/golang-jwt/jwt/v5 v5.2.1`
  - `github.com/joho/godotenv v1.5.1`

### 2. stdlib/config/config.go

- `Framework = "net/http"`, `Port = "8083"`
- **Có `TemplateDir` + `IsProduction`** (giống Gin, xem `conventions.md`)

### 3. stdlib/db/mongo.go

Copy từ Gin.

### 4. stdlib/internal/model/

- `user.go`: copy từ Gin (Phone, json:"-")
- `session.go`: copy từ Gin

### 5. stdlib/internal/repository/

- `user.go`: 6 methods (copy logic từ Gin)
- `session.go`: 4 methods (copy logic từ Gin)

### 6. stdlib/.env.example

```
MONGO_URI=...
MONGO_DB=goauth
JWT_SECRET=...
PORT=8083
# TEMPLATE_DIR=../shared/templates
# APP_ENV=development
```

---

## Anti-Patterns

❌ Không import web framework — chỉ dùng standard library + data dependencies
❌ Prompt này KHÔNG viết handler/middleware/router — Prompt 09 + 10 sẽ làm

---

## Acceptance Criteria

1. `cd stdlib && go build ./...` pass
2. `go.mod` KHÔNG chứa web framework
3. UserRepository 6 methods, SessionRepository 4 methods
4. Config có TemplateDir + IsProduction

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | go.mod NO framework | 🔲 | |
| 2 | Config: Framework="net/http", Port="8083" | 🔲 | |
| 3 | Config: TemplateDir + IsProduction | 🔲 | |
| 4 | db/mongo.go giống Gin | 🔲 | |
| 5 | User model (Phone, json:"-") | 🔲 | |
| 6 | Session model | 🔲 | |
| 7 | UserRepository 6 methods | 🔲 | |
| 8 | SessionRepository 4 methods | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
