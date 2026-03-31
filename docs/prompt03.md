# Prompt 03 — Gin: Config + DB + Model + Repository

## Role

Bạn là một **Go Backend Engineer** chuyên về Clean Architecture. Bạn thiết kế data layer rõ ràng, tách biệt config, database, model, và repository.

---

## Context

Phase 2 bắt đầu. Prompt này xây dựng tầng data cho Gin — config, database connection, models (User + Session), và repositories (UserRepository + SessionRepository).

Tham khảo: `docs/erd.md` cho schema chi tiết, `docs/api-spec.md` cho token strategy.

---

## Dependencies (Prompt phụ thuộc)

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `gin/` folder skeleton |
| 02 | MongoDB Atlas đã seed (2 collections, indexes) |

---

## Yêu cầu

### 1. gin/go.mod

- Module: `github.com/yourusername/go-auth-frameworks/gin`
- Go version: 1.22
- Dependencies:
  - `github.com/gin-gonic/gin v1.9.1`
  - `go.mongodb.org/mongo-driver v1.15.0`
  - `golang.org/x/crypto v0.22.0`
  - `github.com/golang-jwt/jwt/v5 v5.2.1`
  - `github.com/joho/godotenv v1.5.1`

> Chạy `go mod tidy` sau khi tạo go.mod.

### 2. gin/config/config.go

```go
type Config struct {
    MongoURI           string
    MongoDB            string
    JWTSecret          string
    Port               string
    Framework          string // Hardcode = "Gin"
    GoogleClientID     string // Phase 7
    GoogleClientSecret string // Phase 7
    GoogleRedirectURL  string // Phase 7
}
```

- Load từ env bằng `os.Getenv()`
- Load `.env` file bằng `godotenv.Load()` (ignore error — production không có .env)
- `Port` default = `"8081"` nếu env rỗng
- `Framework` hardcode = `"Gin"` (không đọc từ env)
- Validate bắt buộc: `MongoURI`, `MongoDB`, `JWTSecret` — nếu thiếu → log.Fatal
- Google fields: **không bắt buộc** (Phase 7), chỉ log warning

### 3. gin/db/mongo.go

```go
package db

func Connect(cfg *config.Config) (*mongo.Client, *mongo.Database) {
    // Context timeout 10 giây
    // Retry tối đa 3 lần, delay 2 giây giữa mỗi lần
    // Ping để verify connection
    // Return client + database
}

func Disconnect(client *mongo.Client) {
    // Context timeout 5 giây
    // client.Disconnect(ctx)
}
```

### 4. gin/internal/model/user.go

```go
type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name      string             `bson:"name" json:"name"`
    Email     string             `bson:"email" json:"email"`
    Phone     string             `bson:"phone" json:"phone"`
    Password  string             `bson:"password" json:"-"`          // json:"-" ẩn khi serialize
    Role      string             `bson:"role" json:"role"`
    Provider  string             `bson:"provider" json:"-"`          // Phase 7
    GoogleID  string             `bson:"google_id" json:"-"`         // Phase 7
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
```

> **Quan trọng**: `Password`, `Provider`, `GoogleID` có tag `json:"-"` — không bao giờ xuất hiện trong JSON response (xem endpoint `/api/me`).

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

```go
type UserRepository interface {
    Create(ctx context.Context, user *model.User) error
    FindByEmail(ctx context.Context, email string) (*model.User, error)
    FindByID(ctx context.Context, id string) (*model.User, error)
    FindByPhone(ctx context.Context, phone string) (*model.User, error)
    FindByGoogleID(ctx context.Context, googleID string) (*model.User, error) // Phase 7
    UpdateByID(ctx context.Context, id string, update bson.M) error           // Phase 7
}
```

Implement `userRepository` struct với `*mongo.Collection`.

- `FindByEmail`: `filter := bson.M{"email": email}`
- `FindByID`: convert string → ObjectID, `filter := bson.M{"_id": objID}`
- `FindByPhone`: `filter := bson.M{"phone": phone}`
- `FindByGoogleID`: `filter := bson.M{"google_id": googleID}`
- `UpdateByID`: `collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})`
- Constructor: `NewUserRepository(db *mongo.Database) UserRepository`
- Collection name: `"users"`

### 7. gin/internal/repository/session.go

```go
type SessionRepository interface {
    Create(ctx context.Context, session *model.Session) error
    FindByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error)
    DeleteByRefreshToken(ctx context.Context, refreshToken string) error
    DeleteAllByUserID(ctx context.Context, userID primitive.ObjectID) error
}
```

Implement `sessionRepository`:

- `Create`: `collection.InsertOne(ctx, session)`
- `FindByRefreshToken`: `filter := bson.M{"refresh_token": refreshToken}`
- `DeleteByRefreshToken`: `collection.DeleteOne(ctx, bson.M{"refresh_token": refreshToken})`
- `DeleteAllByUserID`: `collection.DeleteMany(ctx, bson.M{"user_id": userID})` — dùng khi muốn revoke all sessions
- Constructor: `NewSessionRepository(db *mongo.Database) SessionRepository`
- Collection name: `"sessions"`

### 8. gin/.env.example

```
MONGO_URI=mongodb+srv://<user>:<pass>@cluster.mongodb.net/?retryWrites=true&w=majority
MONGO_DB=goauth
JWT_SECRET=your-gin-jwt-secret-key-change-me
PORT=8081
```

> **Không có** `FRAMEWORK` — đã hardcode trong config.

---

## Hướng dẫn thực hiện

1. Error handling: wrap error với context bằng `fmt.Errorf("finding user by email: %w", err)`
2. Context timeout: repository methods dùng context từ caller (handler truyền xuống)
3. `FindBy*` methods: nếu không tìm thấy → return `nil, nil` (không phải error)
4. `json:"-"` tag đảm bảo password không lộ qua API
5. Phase 7 methods (`FindByGoogleID`, `UpdateByID`) implement sẵn nhưng chưa dùng tới Phase 7

---

## Anti-Patterns (KHÔNG được làm)

❌ Không return error khi không tìm thấy document — return `nil, nil`
❌ Không hardcode connection string — đọc từ env
❌ Không quên `json:"-"` trên Password, Provider, GoogleID
❌ Không quên Session model và SessionRepository

---

## Acceptance Criteria

1. `cd gin && go build ./...` pass
2. `cd gin && go vet ./...` không warning
3. UserRepository có đúng 6 methods
4. SessionRepository có đúng 4 methods
5. User model có `json:"-"` trên Password, Provider, GoogleID

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | go.mod đúng module + dependencies có version | 🔲 | |
| 2 | Config load đúng + validate MongoURI/MongoDB/JWTSecret | 🔲 | |
| 3 | Config.Framework hardcode = "Gin" | 🔲 | |
| 4 | Config.Port default = "8081" | 🔲 | |
| 5 | db/mongo.go Connect có retry 3 lần | 🔲 | |
| 6 | db/mongo.go Disconnect có context timeout | 🔲 | |
| 7 | User model có Phone field | 🔲 | |
| 8 | User model có json:"-" trên Password/Provider/GoogleID | 🔲 | |
| 9 | Session model có đủ 6 fields | 🔲 | |
| 10 | UserRepository có 6 methods (bao gồm FindByPhone) | 🔲 | |
| 11 | SessionRepository có 4 methods | 🔲 | |
| 12 | FindBy* return nil, nil khi không tìm thấy | 🔲 | |
| 13 | .env.example có 4 biến, không có FRAMEWORK | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
