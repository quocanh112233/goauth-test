package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI is required in environment")
	}

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "goauth"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(dbName)
	fmt.Printf("Connected to Database: %s\n", dbName)

	setupUsers(ctx, db)
	setupSessions(ctx, db)

	seedAdmin(ctx, db)

	fmt.Println("Seed process completed successfully!")
}

func setupUsers(ctx context.Context, db *mongo.Database) {
	fmt.Println("Setting up 'users' collection...")
	coll := db.Collection("users")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "phone", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "google_id", Value: 1}},
			Options: options.Index().SetSparse(true),
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		log.Printf("Warning: error creating user indexes: %v\n", err)
	}
}

func setupSessions(ctx context.Context, db *mongo.Database) {
	fmt.Println("Setting up 'sessions' collection...")
	coll := db.Collection("sessions")

	ttlIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "expired_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	lookupIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{{Key: "refresh_token", Value: 1}}},
	}

	_, err := coll.Indexes().CreateOne(ctx, ttlIndex)
	if err != nil {
		log.Printf("Warning: error creating session TTL index: %v\n", err)
	}

	_, err = coll.Indexes().CreateMany(ctx, lookupIndexes)
	if err != nil {
		log.Printf("Warning: error creating session lookup indexes: %v\n", err)
	}
}

func seedAdmin(ctx context.Context, db *mongo.Database) {
	fmt.Println("Seeding admin user...")
	coll := db.Collection("users")

	email := "admin@goauth.dev"
	var existing bson.M
	err := coll.FindOne(ctx, bson.M{"email": email}).Decode(&existing)

	if err == nil {
		fmt.Println("Admin user already exists, skipping.")
		return
	}

	if err != mongo.ErrNoDocuments {
		log.Fatalf("Error checking for existing admin: %v", err)
	}

	password := "Admin@123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	admin := bson.M{
		"name":       "Admin",
		"phone":      "0900000000",
		"password":   string(hashedPassword),
		"role":       "admin",
		"provider":   "local",
		"google_id":  "",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}

	_, err = coll.InsertOne(ctx, admin)
	if err != nil {
		log.Fatalf("Error inserting admin user: %v", err)
	}

	fmt.Println("Admin user created successfully!")
}
