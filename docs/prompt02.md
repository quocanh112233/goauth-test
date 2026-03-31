# Prompt 02 — MongoDB Atlas Setup + Seed Script

## Role

Bạn là một **Database Administrator** chuyên MongoDB Atlas. Bạn thiết kế schema, tạo indexes (bao gồm TTL index), và viết seed script an toàn.

---

## Context

Tạo database `goauth` trên MongoDB Atlas với 2 collections: `users` và `sessions`. Seed 1 admin account để test. Chi tiết schema xem `docs/erd.md`.

---

## Dependencies (Prompt phụ thuộc)

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `scripts/` folder phải tồn tại |

---

## Yêu cầu

### 1. MongoDB Atlas Setup (hướng dẫn)

- Tạo cluster free tier (M0)
- Tạo database user
- Whitelist IP (0.0.0.0/0 cho dev)
- Lấy connection string

### 2. scripts/go.mod

```
module github.com/yourusername/go-auth-frameworks/scripts
go 1.22
```

Dependencies:
- `go.mongodb.org/mongo-driver v1.15.0`
- `golang.org/x/crypto v0.22.0`
- `github.com/joho/godotenv v1.5.1`

> Chạy `go mod tidy` sau khi tạo go.mod.

### 3. scripts/seed.go

Script idempotent tạo:

**a. Collections + Indexes:**

```go
// Collection users
db.CreateCollection(ctx, "users")

// Unique indexes
usersCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
    {Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
    {Keys: bson.D{{Key: "phone", Value: 1}}, Options: options.Index().SetUnique(true)},
})

// Collection sessions
db.CreateCollection(ctx, "sessions")

// TTL index — tự động xóa sessions hết hạn
sessionsCol.Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys:    bson.D{{Key: "expired_at", Value: 1}},
    Options: options.Index().SetExpireAfterSeconds(0),
})

// Lookup indexes cho sessions
sessionsCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
    {Keys: bson.D{{Key: "user_id", Value: 1}}},
    {Keys: bson.D{{Key: "refresh_token", Value: 1}}},
})
```

**b. Admin user:**

```go
adminUser := bson.M{
    "name":       "Admin",
    "email":      "admin@goauth.dev",
    "phone":      "0900000000",
    "password":   bcryptHash("Admin@123"), // cost=12
    "role":       "admin",
    "provider":   "local",
    "google_id":  "",
    "created_at": time.Now(),
    "updated_at": time.Now(),
}
```

**c. Idempotent logic:**

- Tìm user theo email `admin@goauth.dev`
- Nếu đã tồn tại → log "Admin already exists, skipping" → skip
- Nếu chưa → insert → log "Admin created successfully"
- Nếu collection/index đã tồn tại → không lỗi (MongoDB tự skip)

### 4. scripts/.env.example

```
MONGO_URI=mongodb+srv://<user>:<pass>@cluster.mongodb.net/?retryWrites=true&w=majority
MONGO_DB=goauth
```

---

## Hướng dẫn thực hiện

1. Chạy `go mod tidy` sau khi tạo go.mod
2. TTL index `expireAfterSeconds: 0` nghĩa là xóa ngay khi `expired_at` <= hiện tại
3. Script phải load `.env` từ thư mục `scripts/`
4. Bcrypt cost = 12 (production standard)
5. Khi chạy lại nhiều lần → không tạo duplicate, không lỗi

---

## Anti-Patterns (KHÔNG được làm)

❌ Không để password plaintext trong code seed — phải bcrypt hash
❌ Không drop collection/database — chỉ upsert
❌ Không quên TTL index cho sessions

---

## Acceptance Criteria

1. `cd scripts && go run seed.go` pass (lần đầu tạo admin)
2. Chạy lại lần 2 → log "already exists", không lỗi
3. MongoDB Atlas có 2 collections: `users` (2 indexes), `sessions` (3 indexes)
4. TTL index hoạt động: insert session với `expired_at` = quá khứ → document tự xóa sau ≤60s

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | go.mod đúng module + đủ dependencies (có version) | 🔲 | |
| 2 | Tạo collection users + 2 unique indexes (email, phone) | 🔲 | |
| 3 | Tạo collection sessions + TTL index (expired_at) | 🔲 | |
| 4 | Tạo sessions indexes (user_id, refresh_token) | 🔲 | |
| 5 | Admin user có đủ fields: name, email, phone, password, role | 🔲 | |
| 6 | Admin phone = "0900000000" | 🔲 | |
| 7 | Admin password bcrypt cost=12 | 🔲 | |
| 8 | Script idempotent (chạy lại không lỗi) | 🔲 | |
| 9 | .env.example có MONGO_URI + MONGO_DB | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
