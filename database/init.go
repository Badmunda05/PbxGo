package database

import (
	"context"
	"log"
	"main/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var SudoCollection *mongo.Collection

func Init() {
	if config.AppConfig.MongoURL == "" {
		log.Println("⚠️  MongoDB URL not provided, using in-memory storage only")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig.MongoURL))
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("❌ Failed to ping MongoDB: %v", err)
	}

	MongoClient = client
	SudoCollection = client.Database("pbxgo").Collection("sudo_users")

	log.Println("✅ MongoDB connected successfully")
}
