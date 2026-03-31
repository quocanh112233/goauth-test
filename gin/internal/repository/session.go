package repository

import (
	"context"
	"time"

	"github.com/quocanh112233/goauth-test/gin/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionRepository interface {
	Create(ctx context.Context, session *model.Session) error
	FindByRefreshToken(ctx context.Context, token string) (*model.Session, error)
	DeleteByRefreshToken(ctx context.Context, token string) error
	DeleteAllByUserID(ctx context.Context, userID primitive.ObjectID) error
}

type mongoSessionRepository struct {
	collection *mongo.Collection
}

func NewSessionRepository(db *mongo.Database) SessionRepository {
	return &mongoSessionRepository{
		collection: db.Collection("sessions"),
	}
}

func (r *mongoSessionRepository) Create(ctx context.Context, session *model.Session) error {
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, session)
	return err
}

func (r *mongoSessionRepository) FindByRefreshToken(ctx context.Context, token string) (*model.Session, error) {
	var session model.Session
	err := r.collection.FindOne(ctx, bson.M{"refresh_token": token}).Decode(&session)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &session, err
}

func (r *mongoSessionRepository) DeleteByRefreshToken(ctx context.Context, token string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"refresh_token": token})
	return err
}

func (r *mongoSessionRepository) DeleteAllByUserID(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"user_id": userID})
	return err
}
