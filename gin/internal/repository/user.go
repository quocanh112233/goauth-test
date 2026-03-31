package repository

import (
	"context"
	"time"

	"github.com/quocanh112233/goauth-test/gin/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error)
	FindByPhone(ctx context.Context, phone string) (*model.User, error)
	FindByGoogleID(ctx context.Context, googleID string) (*model.User, error)
	UpdateByID(ctx context.Context, id primitive.ObjectID, update bson.M) error
}

type mongoUserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &mongoUserRepository{
		collection: db.Collection("users"),
	}
}

func (r *mongoUserRepository) Create(ctx context.Context, user *model.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *mongoUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func (r *mongoUserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func (r *mongoUserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"phone": phone}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func (r *mongoUserRepository) FindByGoogleID(ctx context.Context, googleID string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"google_id": googleID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func (r *mongoUserRepository) UpdateByID(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}
