package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(uri string) (*mongo.Client, error) {
	var client *mongo.Client
	var err error

	for i := 1; i <= 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err == nil {
			err = client.Ping(ctx, nil)
			if err == nil {
				log.Printf("Successfully connected to MongoDB (Attempt %d)", i)
				return client, nil
			}
		}

		log.Printf("Failed to connect to MongoDB (Attempt %d): %v", i, err)
		if i < 3 {
			time.Sleep(2 * time.Second)
		}
	}

	return nil, err
}

func Disconnect(client *mongo.Client) {
	if client == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Printf("Error disconnecting MongoDB: %v", err)
	} else {
		log.Println("Disconnected from MongoDB")
	}
}
