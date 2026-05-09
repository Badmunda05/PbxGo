package database

import (
	"context"
	"log/slog"
	"pbxgo/config"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	client         *mongo.Client
	SudoCollection *mongo.Collection
)

func Init() {
	if config.AppConfig.MongoURL == "" {
		slog.Warn("MONGO_URL not set — sudo users stored in-memory only (lost on restart)")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(options.Client().ApplyURI(config.AppConfig.MongoURL))
	if err != nil {
		slog.Error("MongoDB connect failed", "error", err)
		return
	}

	if err = client.Ping(ctx, nil); err != nil {
		slog.Error("MongoDB ping failed", "error", err)
		client = nil
		return
	}

	SudoCollection = client.Database("pbxgo").Collection("sudo_users")
	slog.Info("MongoDB connected", "db", "pbxgo")
}

func Disconnect() {
	if client != nil {
		if err := client.Disconnect(context.Background()); err != nil {
			slog.Error("MongoDB disconnect error", "error", err)
		}
	}
}

func IsConnected() bool {
	return client != nil
}
