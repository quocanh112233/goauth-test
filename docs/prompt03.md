# Prompt 03 — Gin: Config + DB + Model + Repository

## Role

Bạn là một **Go Backend Engineer** chuyên về Clean Architecture.

---

## Context

Phase 2 bắt đầu. Xây dựng tầng data cho Gin. Tham khảo `docs/erd.md`, `docs/conventions.md`.

---

## Dependencies

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `gin/` folder skeleton |
| 02 | MongoDB Atlas đã seed (users + sessions collections, indexes) |

---

## Yêu cầu

### 1. gin/go.mod

- Module: `github.com/yourusername/go-auth-frameworks/gin`
- Go: 1.22
- Dependencies:
  - `github.com/gin-gonic/gin v1.9.1`
  - `go.mongodb.org/mongo-driver v1.15.0`
  - `golang.org/x/crypto v0.22.0`
  - `github.com/golang-jwt/jwt/v5 v5.2.1`
  - `github.com/joho/godotenv v1.5.1`

### 2. gin/config/config.go

```go
type Config struct {
    MongoURI           string
    MongoDB            string
    JWTSecret          string
    Port               string // default "8081"
    Framework          string // hardcode "Gin"
    TemplateDir        string // default "../shared/templates"
    IsProduction       bool   // APP_ENV == "production"
    GoogleClientID     string // Phase 7
    GoogleClientSecret string // Phase 7
    GoogleRedirectURL  string // Phase 7
}

func Load() *Config {
    godotenv.Load() // ignore error

    cfg := &Config{
        MongoURI:    os.Getenv("MONGO_URI"),
        MongoDB:     os.Getenv("MONGO_DB"),
        JWTSecret:   os.Getenv("JWT_SECRET"),
        Port:        getEnv("PORT", "8081"),
        Framework:   "Gin",
        TemplateDir: getEnv("TEMPLATE_DIR", "../shared/templates"),
        IsProduction: os.Getenv("APP_ENV") == "production",
        // Google fields — optional
    }

    // Validate required
    if cfg.MongoURI == "" || cfg.MongoDB == "" || cfg.JWTSecret == "" {
        log.Fatal("Missing required env: MONGO_URI, MONGO_DB, JWT_SECRET")
    }
    return cfg
}
```

> **`TemplateDir`**: Xem `docs/conventions.md` mục 7 — local = `../shared/templates`, Docker = `./shared/templates`.
> **`IsProduction`**: Xem `docs/conventions.md` mục 1 — dùng cho cookie Secure flag.

### 3. gin/db/mongo.go

- `Connect()`: retry 3 lần, delay 2s, timeout 10s
- `Disconnect()`: timeout 5s

### 4. gin/internal/model/user.go

```go
type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name      string             `bson:"name" json:"name"`
    Email     string             `bson:"email" json:"email"`
    Phone     string             `bson:"phone" json:"phone"`
    Password  string             `bson:"password" json:"-"`
    Role      string             `bson:"role" json:"role"`
    Provider  string             `bson:"provider" json:"-"`
    GoogleID  string             `bson:"google_id" json:"-"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
```

### 5. gin/internal/model/session.go

```go
type Session struct {
    ID           primitive.ObjectID `bson:"_id,omitempty"`
    UserID       primitive.ObjectID `bson:"user_id"`
    RefreshToken string             `bson:"refresh_token"`
    ExpiredAt    time.Time          `bson:"expired_at"`
    CreatedAt    time.Time          `bson:"created_at"`
    UpdatedAt    time.Time          `bson:"updated_at"`
}
```

### 6. gin/internal/repository/user.go

6 methods: `Create`, `FindByEmail`, `FindByID`, `FindByPhone`, `FindByGoogleID`, `UpdateByID`

### 7. gin/internal/repository/session.go

4 methods: `Create`, `FindByRefreshToken`, `DeleteByRefreshToken`, `DeleteAllByUserID`

### 8. gin/.env.example

```
MONGO_URI=mongodb+srv://<user>:<pass>@cluster.mongodb.net/?retryWrites=true&w=majority
MONGO_DB=goauth
JWT_SECRET=your-gin-jwt-secret-key-change-me
PORT=8081
# TEMPLATE_DIR=../shared/templates  (default, không cần set local)
# APP_ENV=development               (default, không cần set local)
```

---

## Anti-Patterns

❌ Không return error khi không tìm thấy document — return `nil, nil`
❌ Không quên `json:"-"` trên Password, Provider, GoogleID
❌ Không hardcode template path — dùng `cfg.TemplateDir`

---

## Acceptance Criteria

1. `cd gin && go build ./...` pass
2. Config có `TemplateDir` + `IsProduction` fields
3. UserRepository 6 methods, SessionRepository 4 methods
4. User model có `json:"-"` trên 3 fields

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | go.mod đúng + go mod tidy | 🔲 | |
| 2 | Config: TemplateDir + IsProduction | 🔲 | |
| 3 | Config: Port default "8081", Framework hardcode "Gin" | 🔲 | |
| 4 | db/mongo.go Connect retry 3 lần | 🔲 | |
| 5 | User model có Phone + json:"-" (3 fields) | 🔲 | |
| 6 | Session model đủ 6 fields | 🔲 | |
| 7 | UserRepository 6 methods | 🔲 | |
| 8 | SessionRepository 4 methods | 🔲 | |
| 9 | FindBy* return nil, nil khi not found | 🔲 | |
| 10 | .env.example có comment TEMPLATE_DIR + APP_ENV | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
