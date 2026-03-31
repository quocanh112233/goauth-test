package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
